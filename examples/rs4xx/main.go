package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

/* Creates modbus serial client */
func createHandler() *modbus.RTUClientHandler {
	/* setting handler */
	handler := modbus.NewRTUClientHandler("/dev/ttymxc1")
	handler.SlaveId = 254
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = 2 * time.Second

	return handler
}

func viewResults(client modbus.Client) {
	startAddress := uint16(0)
	quantity := uint16(2)
	results, err := client.ReadHoldingRegisters(startAddress, quantity)
	if err != nil {
		fmt.Println("Error handling data reading: ", err)
	}
	if len(results) >= 2 {
		humidity := float64(results[1]) / 10.0
		temperature := float64(results[3]) / 10.0

		/* Create result and map */
		result := map[string]float64{
			"Humidity":    humidity,
			"Temperature": temperature,
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			log.Fatal("Failed to marshal JSON:", err)
		}
		fmt.Println(string(resultJSON))
	} else {
		log.Fatal("Invalid number of registers returned.")
	}
}

func main() {

	handler := createHandler()

	/* initializing and connecting the client */
	client := modbus.NewClient(handler)
	err := handler.Connect()
	if err != nil {
		log.Fatal("Failed to connect to Modbus device:", err)
	}
	defer handler.Close()

	//processResults(client)
	viewResults(client)
}
