package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"worksp/webapi"
)

// initiating variables; global
// needing mutex for locks on values
var (
	token                string
	startedTime          float64
	started              bool
	totalGlowingTime     float64
	totalGlowingTimeLock sync.Mutex
)

// creating a dio struct to show values
type Dio struct {
	DioName   string `json:"DioName"`
	Value     int    `json:"Value"`
	Direction string `json:"Direction"`
}

// defining struct tokenResponse
// declares field named token with type string -> token = field tag; json key named token
type TokenResponse struct {
	Token string `json:"token"`
}

func getToken() string {
	tokenURL := "http://192.168.0.100:59801/user/Service/token"
	password := "servicelevel"

	form := url.Values{}
	form.Add("password", password)

	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Response status code:", resp.StatusCode)
		return ""
	}

	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	return tokenResp.Token
}

func setDios(token string, js []byte) {
	//needs mutex
	url := "http://192.168.0.100:59801/tdce/dio/SetStates"
	headers := map[string]string{"Content-Type": "application/json"}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(js))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", headers["Content-Type"])

	//get response from server
	client.Do(req)

	//total glowing time
	fmt.Println("Total Glowing Time: ", totalGlowingTime)

	if totalGlowingTime > 0 {
		client.Do(req)
		startGlowTime := time.Now()

		for totalGlowingTime > 0 {
			// calculating remianing time with time.since the glow started
			// times 1000 for milliseconds
			remainingTime := totalGlowingTime - time.Since(startGlowTime).Seconds()*1000
			if remainingTime <= 0 {
				break
			}
			// sleep for 1 ms
			time.Sleep(time.Millisecond)
		}

		// create slice
		js = []byte(`[{ "DioName": "DIO_A", "Value": 0, "Direction": "Output" }]`)
		// new POST request with buffer and slice
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(js))
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", headers["Content-Type"])
		client.Do(req)

		// if the total glowing time isn't 0, insert into the database and print sleep time
		totalGlowingTimeLock.Lock()
		if totalGlowingTime > 0 {
			modifyDatabase(totalGlowingTime)
			fmt.Println("Sleep: ", totalGlowingTime)
			// reset glow time
			totalGlowingTime = 0
		}
		// unlock the resource and exit function
		totalGlowingTimeLock.Unlock()
	}
}

// connect to database
func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:TDC_arch2023@tcp(192.168.0.100:3306)/diobase")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// insert values into database
func modifyDatabase(duration float64) {
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO dios (duration, whattime) VALUES (?, CURRENT_TIMESTAMP())", duration)
	fmt.Println("1 record inserted.")
}

// fetching DIO state
func fetchCurrState(token string) int {
	resp := getDio(token, "http://192.168.0.100:59801/tdce/dio/GetState/DIO_B")
	return resp.Value
}

// gets DIO state
func getDio(token, url string) Dio {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)

	var dio Dio
	// decode the body of the response and store in dio
	json.NewDecoder(resp.Body).Decode(&dio)
	return dio
}

// stops counting time and calculates total glowing time
func stopTime() {
	started = false

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	elapsed := float64(endTime) - startedTime

	//lock so only one routine may update the value
	totalGlowingTimeLock.Lock()
	totalGlowingTime += elapsed
	totalGlowingTimeLock.Unlock()

	js := []byte(`[{ "DioName": "DIO_A", "Value": 1, "Direction": "Output" }]`)
	go setDios(token, js)
}

// starts time
func startTime() {
	startedTime = float64(time.Now().UnixNano() / int64(time.Millisecond))
	started = true
}

// infinite for loop, fetches state then starts or stops counting time accordingly
func welcome() {
	for {
		state := fetchCurrState(token)
		if state == 0 {
			if started {
				stopTime()
			}
		} else {
			if !started {
				startTime()
			}
		}
	}
}

func main() {
	token = getToken()

	// wait group used for syncing threads; more efficient than for loop in main
	// two wait groups
	var wg sync.WaitGroup
	wg.Add(2)

	// internally created function; decrements wait group when finished and initialized web api
	go func() {
		defer wg.Done()
		webapi.InitializeWebApi()
	}()

	//runs main program and decrements work group when process ends
	go func() {
		defer wg.Done()
		welcome()
	}()

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}
