package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var (
	token   string
	anaval  string
	voltage float64
)

// creating a structs to show values
type Analog struct {
	AinName string `json:"AinName"`
	State   string `json:"State"`
}

type AnalogVal struct {
	AinName string  `json:"AinName"`
	Value   float64 `json:"Value"`
}

// for websocket listener object
type AnalogValueChange struct {
	AinName       string  `json:"AinName"`
	PreviousValue float64 `json:"PreviousValue"`
	NewValue      float64 `json:"NewValue"`
}

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

// connect to database
func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:TDC_arch2023@tcp(192.168.0.100:3306)/analog_base")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func sendRequest(token string, url string, method string, js ...[]byte) *http.Response {

	client := &http.Client{}

	var reqBody io.Reader
	if len(js) > 0 {
		reqBody = bytes.NewBuffer(js[0])
	} else {
		reqBody = nil
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	return resp
}

func getAllAnalogStates(token string) []Analog {
	url := "http://192.168.0.100:59801/tdce/analog-inputs/GetStates"
	resp := sendRequest(token, url, "GET")
	defer resp.Body.Close()
	var analogs []Analog
	err := json.NewDecoder(resp.Body).Decode(&analogs)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return analogs
}

func getAllAnalogValues(token string) []AnalogVal {
	url := "http://192.168.0.100:59801/tdce/analog-inputs/GetValues"
	resp := sendRequest(token, url, "GET")
	defer resp.Body.Close()
	var analogValues []AnalogVal
	err := json.NewDecoder(resp.Body).Decode(&analogValues)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return nil
	}
	return analogValues
}

func printAnalogState(token string, ain string) {
	url := "http://192.168.0.100:59801/tdce/analog-inputs/GetState/" + ain
	resp := sendRequest(token, url, "GET")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the content of the response body as a string
	fmt.Printf("\nAnalog Input %s, State: %s", ain, string(body))
}

func getAnalogValue(token string, ain string) (float64, error) {
	url := "http://192.168.0.100:59801/tdce/analog-inputs/GetValue/" + ain
	resp := sendRequest(token, url, "GET")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return 0, err
	}

	// Parse the response body into a float64
	value, err := strconv.ParseFloat(string(body), 32)
	if err != nil {
		fmt.Println("Error parsing float64: ", err)
		return 0, err
	}

	return float64(value), nil
}

func addToDb(value float64) {
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO analogi (value, whattime) VALUES (?, CURRENT_TIMESTAMP())", value)
	fmt.Println("\n1 record inserted.")
}

func listenAinVal() {
	// create webapi

	serverUrl := url.URL{
		Scheme: "ws",
		Host:   "192.168.0.100:31768",
		Path:   "/ws/tdce/analog-inputs/value",
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket: ", err)
	}
	// close when websocket dial is over
	defer conn.Close()

	// handling OS signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		// infinite listening loop
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message: ", err)
				return
			}
			// save the message in json
			var avchange AnalogValueChange
			errM := json.Unmarshal(message, &avchange)
			if errM != nil {
				fmt.Println("Error decoding JSON: ", errM)
				continue
			}
			fmt.Printf("Received AnalogValueChange: %+v\n", avchange)
			addToDb(avchange.NewValue)
		}
	}()

	// interruption signal closing
	select {
	case <-interrupt:
		log.Println("Received interrupt signal. Closing WebSocket connection...")
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Error closing WebSocket:", err)
		}
		time.Sleep(time.Second)
		return
	}
}

func httpGetAin() {

	token = getToken()
	analogs := getAllAnalogStates(token)
	fmt.Println("Printing all analog states:")
	for _, analog := range analogs {
		fmt.Printf("AinName: %s, State: %s\n", analog.AinName, analog.State)
	}
	analogVals := getAllAnalogValues(token)
	fmt.Println("\nAnalog Values:")
	for _, analogVals := range analogVals {
		fmt.Printf("\nAnalog Name: %s, Value: %f\n", analogVals.AinName, analogVals.Value)
	}

	ain := "AIN_A"

	printAnalogState(token, "AIN_A")
	value, err := getAnalogValue(token, ain)
	if err != nil {
		fmt.Print("Problem during function execution.")
		return
	}
	fmt.Printf("\nAnalog Name: %s, Value: %f", ain, value)
}

func main() {

	// creating two goroutines; one works with rest api http gets, other with websocket state change

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {

		// close goroutine after finishing
		defer wg.Done()
		httpGetAin()

	}()

	go func() {
		// close goroutine after finishing
		defer wg.Done()
		listenAinVal()

	}()

	wg.Wait()

}
