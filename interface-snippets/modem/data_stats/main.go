package main

import (
	"encoding/json"
	"fmt"
	o2 "modem/request"
	"sync"
	"time"
)

var (
	token     string
	modemUrls []string
	modfull   modemFull
	modstat   modemStatistics
)

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

type modemStatistics struct {
	Name            string `json:"name"`
	PacketsSent     int    `json:"packetsSent"`
	PacketsReceived int    `json:"packetsReceived"`
	BytesSent       int    `json:"bytesSent"`
	BytesReceived   int    `json:"bytesReceived"`
	ErrorsRx        int    `json:"errorsRx"`
	ErrorsTx        int    `json:"errorsTx"`
	DroppedRx       int    `json:"droppedRx"`
	DroppedTx       int    `json:"droppedTx"`
	Multicast       int    `json:"multicast"`
	Collisions      int    `json:"collisions"`
	Carrier_errors  int    `json:"carrier_errors"`
	Over_errors     int    `json:"over_errors"`
}

func setToken() {
	for {
		token = o2.Authorize()
		/* Sleep for 59 minutes before reset */
		if token != "" {
			time.Sleep(59 * time.Minute)
		}
	}
}

/* Function for handling errors */
func handleError(err error) {
	fmt.Println("Error: ", err)
	panic(err)
}

/* Sets up environment for fetching and printing data */
func fetchData() {
	for {
		if token != "" {
			fetchModemData()
			fmt.Printf("Received details data: %v\n", modfull)
			fmt.Printf("Received statistics data: %v\n", modstat)
			time.Sleep(time.Second)
		}
	}
}

/* Makes OAuth2.0-authenticated request to REST API for fetching modem data */
func fetchModemData() {
	err := json.Unmarshal(o2.MakeROPCRequest("http://192.168.0.100/devicemanager/api/v1/networking/modem/ppp0/details", token), &modfull)
	if err != nil {
		handleError(err)
	}
	err = json.Unmarshal(o2.MakeROPCRequest("http://192.168.0.100/devicemanager/api/v1/networking/modem/ppp0/statistics", token), &modstat)
	if err != nil {
		handleError(err)
	}
}

func main() {
	var wg sync.WaitGroup
	// Adds 2 processes to WaitGroup
	wg.Add(2)

	// Authorize periodically
	go func() {
		defer wg.Done()
		setToken()
	}()
	// Fetch modem data
	go func() {
		defer wg.Done()
		fetchData()
	}()
	wg.Wait()
}
