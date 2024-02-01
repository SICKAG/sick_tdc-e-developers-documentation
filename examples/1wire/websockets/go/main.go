package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	websocket "websocket-1wire/websockets"
)

type Wire1 struct {
	Family         interface{} `json:"Family"`
	FamilyAsString string      `json:"FamilyAsString"`
	FullPath       string      `json:"FullPath"`
	Id             interface{} `json:"Id"`
	IdAsString     string      `json:"IdAsString"`
	LastSeenTime   string      `json:"LastSeenTime"`
	DeviceDetails  string      `json:"DeviceDetails"`
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine for fetching 1wire data
	go func() {
		defer wg.Done()

		conn, err := websocket.OpenWebsocket("ws", "192.168.0.100:31768", "/ws/tdce/onewire/data")
		if err != nil {
			fmt.Println("Error opening websocket: ", err)
			return
		}
		defer conn.Close()

		for {
			msg, err := websocket.ListenOnWS(conn)
			if err != nil {
				fmt.Println("Error fetching data: ", err)
				return
			}

			/* 1wire object */
			var wire1 []Wire1
			err = json.Unmarshal(msg, &wire1)
			if err != nil {
				fmt.Println("Error decoding JSON: ", err)
				return
			}

			for _, item := range wire1 {
				fmt.Printf("Received object: %+v\n", item)

				/* Parsed in case of working with temperature */
				temp, err := strconv.ParseFloat(item.DeviceDetails, 32)
				if err != nil {
					fmt.Println("Error parsing: ", err)
					continue
				}
				fmt.Printf("Temperature: %f\n", temp)
			}
		}
	}()

	wg.Wait()
}
