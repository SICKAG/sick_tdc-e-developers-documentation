package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mqtt/db"
	mq "mqtt/mqttset"
	o2 "mqtt/request"
	ws "mqtt/websockets"
	"net"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

type gps struct {
	Altitude           float32 `json:"Altitude"`
	Course             *string `json:"Course"`
	Fix                int     `json:"Fix"`
	GpsFixAvailable    bool    `json:"GpsFixAvailable"`
	Hdop               float32 `json:"Hdop"`
	Latitude           float32 `json:"Latitude"`
	Longitude          float32 `json:"Longitude"`
	NumberOfSatellites int     `json:"NumberOfSatellites"`
	SpeedKnots         float32 `json:"SpeedKnots"`
	SpeedMph           float32 `json:"SpeedMph"`
	Time               string  `json:"Time"`
}

type modem struct {
	Rssi         int    `json:"Rssi"`
	DataLinkType string `json:"dataLinkType"`
}

type modemFull struct {
	GsmRegistrationStatus   string `json:"gsmRegistrationStatus"`
	UtranRegistrationStatus string `json:"utranRegistrationStatus"`
	AccessTechnology        string `json:"accessTechnology"`
	DataLinkType            string `json:"dataLinkType"`
	OperatorName            string `json:"operatorName"`
	SimStatus               string `json:"simStatus"`
	LatestError             string `json:"latestError"`
	Imei                    string `json:"imei"`
	Ccid                    string `json:"ccid"`
	Imsi                    string `json:"imsi"`
	Rssi                    int    `json:"rssi"`
	Ip                      string `json:"ip"`
	Gateway                 string `json:"gateway"`
	Dns1                    string `json:"dns1"`
	Dns2                    string `json:"dns2"`
	LocalIp                 string `json:"localIp"`
	RemoteIp                string `json:"remoteIp"`
	Segment                 string `json:"segment"`
	SegmentType             string `json:"segmentType"`
	Persist                 bool   `json:"persist"`
	Name                    string `json:"name"`
	Type                    string `json:"type"`
	Enabled                 bool   `json:"enabled"`
	State                   string `json:"state"`
}

type messageObject struct {
	Rfid      string `json:"Rfid"`
	Gps       gps    `json:"Gps"`
	ModemData modem  `json:"ModemData"`
}

/* Defining parameters */
var (
	wg               sync.WaitGroup
	brokerAddress    string
	topic            string
	quos             int
	clientId         string
	username         string
	dbms             string
	connectionString string
	modemUrl         string
	currentGps       gps
	lastBestGps      gps
	token            string
	rfid             string
	modemData        modem
	messageTimer     *time.Timer
	dbTimer          *time.Timer
	client           mqtt.Client
	conn             *sql.DB
)

/* Sets parameters for working with mqtt */
func setParameters() {
	topic = "gps"
	quos = 2
	modemUrl = "http://192.168.0.100/devicemanager/api/v1/networking/modem/ppp0/details"

	/* Setting connection parameters */
	dbms = "mysql"
	connectionString = "root:TDC_arch2023@tcp(localhost:3306)/gpsmqtt"
}

/* Sets Bearer Authentication token from OAuth2.0 */
func setToken() {
	for {
		token = o2.Authorize()
		/* Sleep for 59 minutes before reset */
		if token != "" {
			time.Sleep(59 * time.Minute)
		}
	}
}

/* Opens connection to websocket and receives last websocket object */
func getWsData(conn *websocket.Conn) gps {
	receivedObject, erro := ws.ListenOnWS(conn)
	var gps_obj gps
	if erro != nil {
		fmt.Println("Error fetching gps: ", erro)
	}
	json.Unmarshal(receivedObject, &gps_obj)
	return gps_obj
}

/* Fetches GPS data from websocket */
func fetchGpsData() {
	conn, _ := ws.OpenWebsocket("ws", "192.168.0.100:31768", "/ws/tdce/gps/data")
	defer conn.Close()

	for {
		currentGps = getWsData(conn)
		/* If the fix of the currently fetched gps object is equal to or better than the value in lastBestGps, set lastBestGps to this current value */
		// 0 - no fix, 1 - best fix, 2 - ok fix
		if currentGps.Fix <= lastBestGps.Fix && currentGps.Fix != 0 {
			lastBestGps = currentGps
		}
	}

}

/* Serves a TCP on localhost:1234 */
func serveTcp() {
	listener, err := net.Listen("tcp", "localhost:5247")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 5247")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleClient(conn)
	}
}

/* Creates data that will be sent to broker */
func createMessage(rfid string) ([]byte, error) {
	var msg messageObject
	msg.Rfid = rfid
	msg.Gps = lastBestGps
	msg.ModemData = modemData

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

/* Makes OAuth2.0 authenticated request to REST API for fetching modem data */
func fetchModemData() {
	var modfull modemFull
	err := json.Unmarshal(o2.MakeROPCRequest(modemUrl, token), &modfull)
	if err != nil {
		fmt.Println("Error decoding: ", err)
	}
	/* Setting values that will be shown in message */
	modemData.Rssi = modfull.Rssi
	modemData.DataLinkType = modfull.DataLinkType
}

/* Creates messages that will be published */
func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	clientDisconnected := false

	for !clientDisconnected {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				clientDisconnected = true
			} else {
				fmt.Println("Error:", err)
			}
		}
		if n > 0 {
			fmt.Printf("Received: %s\n", buffer[:n])
			rfid = string(buffer[:n])
			fetchModemData()
			msge, err := createMessage(rfid)
			check(err)
			publishToMqtt(msge)
		}
	}
}

/* Inserts or deletes data from database */
func modifyDB(conn *sql.DB, method string, data []byte, id int) {
	var query string

	switch method {
	case "INSERT":
		var msg messageObject
		json.Unmarshal(data, &msg)
		var courseValue string

		/* Set course to null if there is no data for it; do the same for gateway if it will be implemented */
		if msg.Gps.Course != nil && *msg.Gps.Course != "" {
			courseValue = fmt.Sprintf("'%s'", *msg.Gps.Course)
		} else {
			courseValue = "NULL"
		}

		query = fmt.Sprintf(`INSERT INTO mqtt (Altitude, Course, Fix, GpsFixAvailable, Hdop, Latitude, 
        Longitude, NumberOfSatellites, SpeedKnots, SpeedMph, TimeSt, Rssi, DataLinkType, 
        Rfid) VALUES (%f, %s, %d, %t, %f, %f, %f, %d, %f, %f, '%s', %d, '%s', '%s');`,
			msg.Gps.Altitude, courseValue, msg.Gps.Fix, msg.Gps.GpsFixAvailable,
			msg.Gps.Hdop, msg.Gps.Latitude, msg.Gps.Longitude, msg.Gps.NumberOfSatellites,
			msg.Gps.SpeedKnots, msg.Gps.SpeedMph, msg.Gps.Time, msg.ModemData.Rssi,
			msg.ModemData.DataLinkType, msg.Rfid)

		db.ModifyDB(conn, query)
		break

	case "DELETE":
		query = fmt.Sprintf("DELETE FROM mqtt WHERE Id = %d", id)
		db.ModifyDB(conn, query)
		break
	}
}

/* Converts database mqtt object to message object for publishing */
func convertToMessageObjects(queuedMessages []db.MQTTData) []messageObject {
	var messageObjects []messageObject

	for _, d := range queuedMessages {
		gpsData := gps{
			Altitude:           d.Altitude,
			Course:             d.Course,
			Fix:                d.Fix,
			GpsFixAvailable:    d.GpsFixAvailable,
			Hdop:               d.Hdop,
			Latitude:           d.Latitude,
			Longitude:          d.Longitude,
			NumberOfSatellites: d.NumberOfSatellites,
			SpeedKnots:         d.SpeedKnots,
			SpeedMph:           d.SpeedMph,
			Time:               d.Time,
		}

		modemData := modem{
			Rssi:         d.Rssi,
			DataLinkType: d.DataLinkType,
		}

		messageObj := messageObject{
			Rfid:      d.Rfid,
			Gps:       gpsData,
			ModemData: modemData,
		}

		messageObjects = append(messageObjects, messageObj)
	}
	return messageObjects
}

/* Tries to publish the MQTT message */
func tryToPublish(msge []byte) {
	success := mq.PublishMessage(topic, msge, client, byte(quos))
	if success == 1 {
		messageTimer.Reset(5 * time.Minute)
	}
}

/* Handles message publishing */
/* If no message is published and timer has run out, insert the message into database; else publish the message and if it is a success, reset the timer */
func publishToMqtt(msge []byte) {
	select {
	// if the message timer ran out, publish insert to database and still try to publish afterwards; otherwise just publish the message
	case <-messageTimer.C:
		modifyDB(conn, "INSERT", msge, 0)
	default:
		tryToPublish(msge)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func checkDatabase() {
	for {
		select {
		// if the timer runs out
		case <-dbTimer.C:
			queuedMessages := db.SelectValues(conn, "SELECT * FROM mqtt")
			messageObjects := convertToMessageObjects(queuedMessages)

			for i, msg := range messageObjects {
				jsonData, err := json.Marshal(msg)
				if err != nil {
					fmt.Println("Error marshaling JSON: ", err)
					return
				}
				/* If message publishing is success, the message will be deleted from the database */
				/* If it isn't successfully published, it will stay in the database and its deletion will be retried when timer runs out again */
				success := mq.PublishMessage(topic, jsonData, client, byte(quos))

				// if the message is published successfully, delete the message from the database
				if success == 1 {
					modifyDB(conn, "DELETE", nil, queuedMessages[i].Id)
				}
			}
			/* After checking the messages stored in the database, the timer should be reset */
			dbTimer.Reset(30 * time.Minute)
		}
	}
}

func main() {

	setParameters()

	/* Open database connection */
	conn = db.Connect(dbms, connectionString)
	defer conn.Close()

	/* Open broker connection - broker stays online even if message isn't published */
	brokerAddress = "tcp://localhost:1883"
	// brokerAddress = "tcp://192.168.0.100:1883"
	clientId = "clientest"
	username = "testerE"
	client = mq.CreateMqttClient(brokerAddress, clientId, username)
	mq.ConnectClientToBroker(client)

	/* Creating timers */
	messageTimer = time.NewTimer(5 * time.Minute)
	dbTimer = time.NewTimer(30 * time.Minute)

	/* Start 4 separate goroutines */
	wg.Add(4)

	// Authorize
	go func() {
		defer wg.Done()
		setToken()
	}()

	// GPS websocket
	go func() {
		defer wg.Done()
		fetchGpsData()
	}()

	// TCP server
	go func() {
		defer wg.Done()
		serveTcp()
	}()

	// Database Check Timer
	go func() {
		defer wg.Done()
		checkDatabase()
	}()

	wg.Wait()
}

/* DATABASE STRUCTURE */
/*create table mqtt (
	Id int primary key auto_increment,
	Altitude float,
    Course varchar(45),
    Fix int,
    GpsFixAvailable bool,
    Hdop float,
    Latitude float,
    Longitude float,
    NumberOfSatellites int,
    SpeedKnots float,
    SpeedMph float,
    TimeSt varchar(45),
	Rssi integer,
    DataLinkType varchar(45),
    Rfid varchar(45) );*/
