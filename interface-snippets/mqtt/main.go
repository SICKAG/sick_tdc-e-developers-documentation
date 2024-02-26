package main

import (
	mq "mqtt-conn/mqttset"
	"time"
)

func connectToBroker() {

	brokerAddress := "tcp://192.168.0.100:1883"
	clientId := "clientest"

	// NOTE: for safety, specify the username and password parameters in another file
	username := "user1"
	password := "password"

	client := mq.CreateMqttClientPass(brokerAddress, clientId, username, password)
	mq.ConnectClientToBroker(client)
	for {
		mq.PublishMessage("test", []byte("test"), client, 0)
		time.Sleep(10 * time.Second)
	}
}

func main() {
	connectToBroker()
}
