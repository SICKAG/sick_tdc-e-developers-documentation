package main

import (
	"encoding/json"
	"fmt"
	o2 "modemmsg/request"
	"sync"
	"time"
)

type modemMsg struct {
	Index   int    `json:"index"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Pdu     string `json:"pdu"`
}

type postMsg struct {
	PhoneNumber string `json:"phoneNumber"`
	Content     string `json:"content"`
}

var (
	token string
	msg   []modemMsg
)

/* Function for handling errors */
func handleError(err error) {
	fmt.Println("Error: ", err)
	panic(err)
}

/* Sets OAuth2.0 token every 59 minutes (expires after 60 minutes) */
func setToken() {
	for {
		token = o2.Authorize()
		/* Sleep for 59 minutes before reset */
		if token != "" {
			time.Sleep(59 * time.Minute)
		}
	}
}

/* Fetches mdoem messages from REST API */
func fetchMessages() {
	for {
		time.Sleep(2 * time.Second)
		err := json.Unmarshal(o2.MakeROPCRequest("http://192.168.0.100/devicemanager/api/v1/networking/modem/ppp0/sms/messages", token), &msg)
		if err != nil {
			handleError(err)
		}
		fmt.Println(msg)
	}
}

/* Posts message with defined parameters */
func postMessage() {
	time.Sleep(5 * time.Second)
	if token != "" {
		// insert phone number here
		phoneNumber := "+XXXXXXXXXXXX"
		content := "hello world"
		url := "http://192.168.0.100/devicemanager/api/v1/networking/modem/ppp0/sms/messages"

		// setting values
		var poster postMsg
		poster.PhoneNumber = phoneNumber
		poster.Content = content

		// setting up a JSON object from created object
		jsonData, err := json.Marshal(poster)
		if err != nil {
			handleError(err)
		}
		o2.PostROPCMessage(url, token, jsonData)
	} else {
		fmt.Println("Token wasn't fetched. Try again.")
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	// setting token
	go func() {
		defer wg.Done()
		setToken()
	}()

	// fecth REST API messages
	go func() {
		defer wg.Done()
		fetchMessages()
	}()

	// send REST API messages
	go func() {
		defer wg.Done()
		postMessage()
	}()

	// wait for all goroutines to finish
	wg.Wait()
}
