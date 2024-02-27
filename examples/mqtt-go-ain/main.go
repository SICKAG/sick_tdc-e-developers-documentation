package main

import (
	mq "mqtt-conn/mqttset"
	ws "mqtt-conn/websockets"
	"time"

	"github.com/gorilla/websocket"
)

/* Function panics upon receiving error */
func check(err error) {
	if err != nil {
		panic(err)
	}
}

/* Fetches AIN value from websocket */
func fetchAinVal(conn *websocket.Conn) []byte {

	obj, err := ws.ListenOnWS(conn)
	check(err)
	return obj
}

/* Sets broker parameters, connects to broker and sends AIN message */
func connectToBroker() {

	/* Setting up MQTT broker */
	brokerAddress := "tcp://192.168.0.100:1883"
	clientId := "clientest"

	// for safety, specify the username and password parameters in another file
	username := "user1"
	password := "password"

	client := mq.CreateMqttClientPass(brokerAddress, clientId, username, password)
	mq.ConnectClientToBroker(client)

	/* Setting up WebSocket */
	conn, err := ws.OpenWebsocket("ws", "192.168.0.100:31768", "/ws/tdce/analog-inputs/value")
	check(err)
	defer conn.Close()

	// fetching AIN value every ten seconds
	// publishes with qos 0
	for {
		mq.PublishMessage("ainval", fetchAinVal(conn), client, 0)
		time.Sleep(10 * time.Second)
	}

}

func main() {
	connectToBroker()
}
