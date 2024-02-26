/* Package created 19.09.2023. */
/* Handles MQTT connections, message sending and subscription */

package mqttset

import (
	"fmt"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var wg sync.WaitGroup

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %+v\n", err)
}

/* Creates options for new mqtt client with password */
func CreateMqttClientPass(brokerAddress string, clientId string, username string, password string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(brokerAddress).SetClientID(clientId).SetUsername(username).SetPassword(password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	return client
}

/* Function publishes a message or queues it to Broker */
/* Testing phase -> return successes  */
func PublishMessage(topic string, message []byte, client mqtt.Client, quos byte) int {
	/* creating a success channel because of internal goroutine */
	successChan := make(chan int, 1)
	defer close(successChan)
	go func() {
		token := client.Publish(topic, quos, false, message)
		if token.WaitTimeout(5 * time.Second) {
			fmt.Printf("Published message: %s\n", message)
			successChan <- 1
		} else {
			/* Message is queued */
			fmt.Printf("Failed to publish message: %s. Message queued.\n", message)
			successChan <- 0
		}
	}()

	/* waiting for the goroutine to signal success or failure */
	success := <-successChan
	return success
}

/* Connects to a broker */
func ConnectClientToBroker(client mqtt.Client) {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error connecting to broker: ", token.Error())
		os.Exit(1)
	}
}

/* Function subscribes a client to a broker */
func SubscribeToTopic(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s.\n", topic)
}
