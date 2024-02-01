package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// fetches data indefinitely until signal interrupt
func getData() {
	// setting up URL to fetch data from
	serverUrl := url.URL{
		Scheme: "ws",                  //wss is secure
		Host:   "192.168.0.100:31768", // address
		Path:   "/ws/tdce/gps/data",   //endpoint path
	}

	// dials a websocket with the created URL
	conn, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket: ", err)
	}
	// makes sure to close the connection after work with websocket has finished
	defer conn.Close()

	// handling OS signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// goroutine for listening to messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message: ", err)
				return
			}
			// prints the message
			fmt.Printf("received message: %s\n", message)
		}
	}()

	select {
	// in case of an interruption signal, close the websocket normally, then sleep for a second
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

func main() {
	getData()
}
