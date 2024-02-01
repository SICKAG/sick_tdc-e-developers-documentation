package main

import (
	"fmt"
	"log"
	"sync"
	websocket "websocket-rs232/websockets"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	/* Open connection for reading and writing to websocket */
	conn, err := websocket.OpenWebsocket("ws", "192.168.0.100:31768", "/ws/tdce/rs232/data")
	if err != nil {
		fmt.Println("Error opening websocket: ", err)
	}
	defer conn.Close()

	/* Goroutine for data fetching */
	go func() {
		defer wg.Done()
		fmt.Println("Listening on websocket...")
		for {
			msg, erro := websocket.ListenOnWS(conn)
			if erro != nil {
				fmt.Println("Error fetching data: ", erro)
				break
			}
			receivedString := string(msg)
			fmt.Printf("Received string: %s\n", receivedString)
		}
	}()

	/* Sending data to websocket */
	/* Specify data here... */
	message := "dGVzdAo="
	if err := websocket.SendToWS(conn, []byte(message)); err != nil {
		log.Println("Error sending message: ", err)
	}
	wg.Wait()
}
