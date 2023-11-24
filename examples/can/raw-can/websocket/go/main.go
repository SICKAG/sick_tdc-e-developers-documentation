package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type CanBus struct {
	CanBusName                  string `json:"CanBusName"`
	Id                          int    `json:"Id"`
	Data                        []byte `json:"Data"`
	IsErrorFrame                bool   `json:"IsErrorFrame"`
	IsExtendedFrameFormat       bool   `json:"IsExtendedFrameFormat"`
	IsRemoteTransmissionRequest bool   `json:"IsRemoteTransmissionRequest"`
}

/* opens websocket by creating a url object with the provided scheme, host and path */
func OpenWebsocket(scheme, host, path string) (*websocket.Conn, error) {

	serverUrl := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket: ", err)
		return nil, err
	}

	/* to close connection write defer conn.Close() in calling package */
	return conn, nil
}

/* listens on the provided websocket connection; use after opening the websocket and open in separate goroutine*/
/* implemented with an infinite loop */
func ListenOnWS(conn *websocket.Conn) ([]byte, error) {

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			return nil, err
		}
		/* returns bytes from the websocket */
		return message, nil
	}

}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := OpenWebsocket("ws", "192.168.0.100:31768", "/ws/tdce/can-a/data")
		if err != nil {
			fmt.Println("Error opening websocket: ", err)
		}
		defer conn.Close()

		for {
			msg, erro := ListenOnWS(conn)
			if erro != nil {
				fmt.Println("Error fetching data: ", erro)
			}
			var canBus CanBus
			json.Unmarshal(msg, &canBus)
			fmt.Printf("Received Object: %v\n", canBus)
		}

	}()
	wg.Wait()
}
