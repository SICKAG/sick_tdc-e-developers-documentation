---
title: GNSS Examples
sidebar_position: 207
---

import Admonition from '@theme/Admonition';


# GNSS Examples


In this section, GNSS on the TDC-E device is discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that checks device availability and fetches GNSS data is provided.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/gnss.proto
```

This installs needed tools, and generates files for the gRPC service. The application also requires an access token to connect to the gRPC service the TDC-E provides. Please refer to [gRPC Usage](/getting-started/grpc-usage) for instructions on how to generate an access token for the TDC-E.

Once you've obtained the token, navigate to `grpc-gnss/pkg/auth/token.json`. Paste your generated token in the `access_token` field. This token is then used to create a `context` which will be used to create gRPC calls.

The application then creates a gRPC client that connects to the service unix socket. It does so by using the generated gRPC Go file specification.

```go
conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
```

The application first checks the GNSS device status by calling the HAL service for getting the device status. If the transceiver isn't set up, it sends a request to enable the receiver and antenna of the device.

```go
func setUpGnss(client protos.GnssClient, ctx context.Context) {
	req := &protos.DeviceManagementRequest{
		EnableReciever: true,
		EnableAntenna:  true,
	}
	_, err := client.SetDeviceStatus(ctx, req)
	if err != nil {
		log.Fatalf("could not set status: %v", err)
	}
}
```

The application then starts streaming GNSS data. This includes a timestamp, the position of the device, device speed, course, the number of satellites, fix type, fix quality and HDOP. An example log of the application is provided below:

```log
2023/09/20 13:01:52 Timestamp: 2024-12-19T11:13:53Z
2023/09/20 13:01:52 Latitude: 46.286890516666666
2023/09/20 13:01:52 Longitude: 16.321311299999998
2023/09/20 13:01:52 Altitude: 183.8
2023/09/20 13:01:52 Speed (Kph): 0
2023/09/20 13:01:52 Speed (Mph): 0
2023/09/20 13:01:52 Course: 0
2023/09/20 13:01:52 Satellites: 10
2023/09/20 13:01:52 Fix Type: 3
2023/09/20 13:01:52 Fix Quality: 1
2023/09/20 13:01:52 HDOP: 0.8
```

To stream the GNSS data, the application sends a HAL Service request to the TDC device. A stream is opened, which lists GNSS messages indefinitely. This is shown in the code below:

```go
func streamGnssJson(client protos.GnssClient, ctx context.Context) {
	stream, err := client.StreamDataJson(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("Error while calling StreamDataJson: %v", err)
	}
	log.Printf("--------------------------------")
	for {
		data, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error while receiving data: %v", err)
		}

		log.Printf("Timestamp: %v", data.Timestamp)
		log.Printf("Latitude: %v", data.Latitude)
		log.Printf("Longitude: %v", data.Longitude)
		log.Printf("Altitude: %v", data.Altitude)
		log.Printf("Speed (Kph): %v", data.SpeedKph)
		log.Printf("Speed (Mph): %v", data.SpeedMph)
		log.Printf("Course: %v", data.Course)
		log.Printf("Satellites: %v", data.Satellites)
		log.Printf("Fix Type: %v", data.FixType)
		log.Printf("Fix Quality: %v", data.FixQuality)
		log.Printf("HDOP: %v", data.Hdop)
		log.Printf("--------------------------------")

		time.Sleep(1 * time.Second)
	}
}
```

### 1.2. Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the GNSS service is the following:

```protobuf
/**
 * GNSS Service.
 *
 * Service for reading GNSS data
 */
 syntax = "proto3";

 package hal.gnss;
 
 import "google/protobuf/empty.proto";
 
 option go_package = "./protos;protos";

enum NMEAFrequency {
    _1Hz = 0;
    _2Hz = 1;
    _5Hz = 2;
    _10Hz = 3;
}

message NMEAFrequenciesList {
    repeated string frequencies = 1;
}

message NMEAFrequencySettings {
    NMEAFrequency frequency = 1;
}

message GnssDataJson {
    double latitude = 1;     // The latitude of the GNSS position in decimal degrees.
    double longitude = 2;    // The longitude of the GNSS position in decimal degrees.
    double altitude = 3;     // The altitude of the GNSS position in meters.
    string timestamp = 4;    // The timestamp of the GNSS record in ISO 8601 format.
    double speed_kph = 5;    // The speed of the object in kilometers per hour.
    double speed_mph = 6;    // The speed of the object in miles per hour.
    double course = 7;       // The direction of travel in degrees, where 0 is north.
    double satellites = 8;   // The number of satellites used to determine the position.
    double fix_type = 9;     // Indicates the type of GNSS fix (e.g., 1 for no fix, 2 for 2D fix, 3 for 3D fix).
    double fix_quality = 10;  // Indicates the GNSS fix quality (e.g., 0 for invalid, 1 for GPS fix, 2 for DGPS fix).
    double hdop = 11;        // The Horizontal Dilution of Precision, a measure of the GNSS signal's accuracy.
}

message GnssData{
    string sentence = 1;
}

message DeviceManagementRequest {
    bool enableReciever = 1;
    bool enableAntenna = 2;
}

message DeviceManagementResponse {
    bool enableReciever = 1;
    bool enableAntenna = 2;
    bool sessionActive = 3;
    bool rebootRequired = 4;
}

message GnssConstellationsSettings {
    bool GLONASS = 1;
    bool BDS = 2;
    bool Galileo = 3;
}

message GnssConstellations {
    repeated string constellations = 1;
}


 /**
  * Service exposing GNSS functions.
  */
  service Gnss {
    // Provides parsed Gnss data as JSON
    rpc StreamDataJson(google.protobuf.Empty) returns (stream GnssDataJson) {}
  
    // Provides raw stream that is returned by the Gnss receiver
    rpc StreamDataRaw(google.protobuf.Empty) returns (stream GnssData) {}

    // Provides a way do get receiver and antenna status
    rpc GetDeviceStatus(google.protobuf.Empty) returns (DeviceManagementResponse) {}
  
    // Provides a way do disable/enable receiver and antenna
    rpc SetDeviceStatus(DeviceManagementRequest) returns (google.protobuf.Empty) {}
  
    // Provides a way to set preferred constellations
    rpc SetConstellations(GnssConstellationsSettings) returns (google.protobuf.Empty) {}
  
    // Provides a way to get current preferred constellations
    rpc GetConstellations(google.protobuf.Empty) returns (GnssConstellationsSettings) {}
  
    // Provides a way to get a list of available constellations
    rpc ListAvailableConstellations(google.protobuf.Empty) returns (GnssConstellations) {}
  
    // Provides a way to set NMEA output frequency
    rpc SetNmeaOutputFrequency(NMEAFrequencySettings) returns (google.protobuf.Empty) {}
  
    // Provides a way to set NMEA output frequency
    rpc GetNmeaOutputFrequency(google.protobuf.Empty) returns (NMEAFrequencySettings) {}
  
    // Provides a way to get a list of available NMEA output frequencies
    rpc ListNmeaOutputFrequencies(google.protobuf.Empty) returns (NMEAFrequenciesList) {}
  }

```

</details>

To get the GNSS device status, use `grpcurl`, which is an open-source utility for accessing gRPC services via the shell. For help setting up the `grpcurl` command, refer to [gRPC Usage](/getting-started/grpc-usage).

Use the following line: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> hal.gnss.Gnss.GetDeviceStatus
```

The `token` field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage). The `device-ip:grpc_server_port` is the TDC-E IP address and the gRPC serving port. For example, if the `token` value was `token` and the address and port were `192.168.0.100:8081`, you would use the following line to see the device status of the GNSS service.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.gnss.Gnss.GetDeviceStatus
```

The response should be in this format:

```log
"enableReciever": true,
"enableAntenna": true,
"sessionActive": true,
"rebootRequired": false
```

Additionally, you can use the `gRPC Clicker` VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

### 1.3. Application Deployment

#### 1.2.1. Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.


Open a terminal and paste the following commands:

```bash
docker build -t gnss-grpc-app .
docker save -o gnss-grpc-app.tar gnss-grpc-app:latest
```

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 1.2.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```


The application is built as `gnss-grpc`.

```dockerfile
RUN go build -o gnss-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `gnss-grpc` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gnss-grpc .

CMD ["./gnss-grpc"]
```

#### 1.2.3. Deploying to Portainer

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## 2. Node-RED gRPC Example

<a href="../code/node-red-examples/grpc-gnss.json" download>Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>Download Node</a>

This section describes usage of the the GNSS device on the TDC-E. A Node-RED gRPC application was created to demonstrate reading from the GNSS device.

For implementation, the following nodes are used:

- `inject` node
- `gRPC call` node
- `debug` node.

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

See a screenshot of the application below:

![GNSS Node-RED Example](/img/node-red-gnss.png)

## 3. Lua Example

<a href="../code/lua-examples/gnss.lua" download>Download Example Code</a>

A Lua application is provided as a GNSS usage example.

The script first creates a GNSS handle and enables the GNSS service. The handle is used to fetch all needed data from the receiver.

```lua
gnssHandle = Gnss.create()
gnssHandle:enable()
```

The following metrics are fetched:
- position,
- signal metrics,
- NMEA sentences,
- speed, and
- course.

It also provides a timestamp of the precise time when the application fetches the data. To get the position of your device, which includes latitude, longitude and altitude, the following code is used:

```lua
local x,y,z = gnssHandle:getPosition()
print("Latitude: " .. x)
print("Longitude: " .. y)
print("Altitude: " .. z)
```

To get signal metrics, the following code is used:

```lua
local q,w,e,r = gnssHandle:getSignalMetrics()
print("Number of satellites: " .. q)
print("Fix type: " .. w)
print("Fix quality: " .. e)
print("HDOP: " .. r)
```

This includes the number of satellites, the type of the fix used, fix quality and HDOP.

NMEA sentences are fetched using the code below:

```lua
local h=gnssHandle:getNmeaSentence('RMC')
print('RMC sentence: ' .. h)
h=gnssHandle:getNmeaSentence('GSA')
print('GSA sentence: ' .. h)
h=gnssHandle:getNmeaSentence('GGA')
print('GGA sentence: ' .. h)
h=gnssHandle:getNmeaSentence('GSV')
print('GSV sentence: ' .. h)
h=gnssHandle:getNmeaSentence('VTG')
print('VTG sentence: ' .. h)
```

In the end, speed and course is fetched:

```lua
print("Speed in KPH: " .. gnssHandle:getSpeed('KPH'))
print("Speed MPH:" .. gnssHandle:getSpeed('MPH'))
print("Course " .. gnssHandle:getCourse())
```

See the result of the application run below.

```log
[12:08:27.359: INFO: AppEngine] Starting app: 'gnss' (priority: LOW)
GNSS Receiver is enabled
-------------------------------
Timestamp: 2024-12-19T10:21:54Z
Latitude: 46.28689185
Longitude: 16.321358916667
Altitude: 181.2
Speed [KPH]: 0.0
Speed [MPH]: 0.0
Course: 36.3
Satellite number: 9.0
Fix type: 3.0
Fix quality: 1.0
HDOP: 0.9
-------------------------------
```
