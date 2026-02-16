---
title: "TDC-E Documentation"
---


<!-- Source: docs/code-samples/tdce-serial-examples.md -->


# Serial Examples

In this section, the Serial interface of the TDC-E device and examples of its usage are discussed. Programming examples are given and thoroughly explained.

> **ℹ️ Info**
>
> The examples below are written for use inside **your user workspace container**.
>
> When creating **additional containers**, make sure to use **host paths** when mapping volumes or devices.
>
> See [Device-specific Paths](/getting-started/working-with-docker#13-device-specific-paths) for the correct mount points and device nodes available on your device.


> **ℹ️ TDC-E Serial Interfaces**
>
>
> The TDC-E device has two separate serial interfaces:
>
>
>   - **SERIAL_1** (`/dev/ttymxc5`): RS232 only
>   - **SERIAL_2** (`/dev/ttymxc1`): RS422 or RS485
>
>
> Other devices have a single **SERIAL** interface that supports multiple modes.
>


## Go Example

This section handles setting up a serial device using gRPC, and creating serial Go applications. 

The first application is made and tested with the 'Leaf Wetness Sensor', which operates with a 'Baudrate' of '9600', using the communication protocol 'MODBUS', and with 'RS485'. For this application, an 'Isolated Soil Sensor' was used.

The second application was created by simulating a 'RS422' device using two TDC-E devices and connecting their serial interfaces. The reading and writing functionalities are implemented as separate applications. The 'RS232' device was created using loopback mode by connecting TX and RX wires.

### Setting up Serial Device

In this subsection, setting up the serial device is discussed. Using the dedicated serial HAL service, setup of the following parameters is possible:
 
 - Mode
 - Transceiver Power
 - Slew Rate

In the following subsections, setting up the mode of a serial device is discussed. This is done using the Serial HAL Service 'hal.serial.Serial'. Examples using 'grpcurl' are given below. For more information about using gRPC services, refer to [gRPC Usage](/getting-started/grpc-usage).

#### Setting Up Serial Mode

The Serial HAL service allows setting up the serial mode. The TDC-E device supports three modes across its two serial interfaces:

 - RS232 (SERIAL_1 only)
 - RS422 (SERIAL_2 only)
 - RS485 (SERIAL_2 only)

To set up your serial mode, use the 'hal.serial.Serial.SetMode' service. Examples are given below.

**RS232 (SERIAL_1):**

```
grpcurl -emit-defaults -H 'Authorization: Bearer {token}' -plaintext -
d '{"interfaceName":"SERIAL_1","mode":"RS232"}' 192.168.0.100:8081 hal.serial.
Serial.SetMode
{}
```

**RS422 or RS485 (SERIAL_2):**

```
grpcurl -emit-defaults -H 'Authorization: Bearer {token}' -plaintext -
d '{"interfaceName":"SERIAL_2","mode":"RS422"}' 192.168.0.100:8081 hal.serial.
Serial.SetMode
{}
```

To check changes to the serial mode, use the HAL service 'hal.serial.Serial.GetMode'.

```
grpcurl -emit-defaults -H 'Authorization: Bearer {token}' -plaintext -
d '{"interfaceName":"SERIAL_2"}' 192.168.0.100:8081 hal.serial.Serial.GetMode
{
  "mode": "RS422"
}
```

#### Viewing Serial Statistics

The Serial HAL service provides a means to view serial statistics. The HAL service 'hal.serial.Serial.GetStatistics' is used. An example of seeing serial statistics for the RS422/RS485 interface (SERIAL_2) is given below.

```
grpcurl -d '{"interfaceName":"SERIAL_2"}' -H 'Authorization: Bearer token' -
plaintext 192.168.0.100:8081 hal.serial.Serial.GetStatistics
{
  "txCount": "1050",
  "rxCount": "982"
}
```

For RS232 interface (SERIAL_1), use '"interfaceName":"SERIAL_1"'.

#### Proto File

> **📝 TDC-E Limitations**
>
>
> The `SetTermination` and `GetTermination` services are not available on TDC-E devices. 
> Both SERIAL_1 and SERIAL_2 interfaces on TDC-E do not support termination configuration.
>


<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the Serial service is the following:

```bash
syntax = "proto3";

package hal.serial;

import "google/protobuf/empty.proto";

option go_package = "./protos;protos";

// Request to set transceiver power
message SetTransceiverPowerRequest {
  string interfaceName = 1; // e.g., "serial1"
  bool powerOn = 2;         // true to power on, false to power off
}

// Enum for serial interface modes
enum SerialMode {
  RS485 = 0; // RS485 mode
  RS422 = 1; // RS422 mode
  RS232 = 2; // RS232 mode
}

// Enum for serial interface slew rate mode
enum SlewRateMode {
  HIGH_SPEED = 0; // High speed mode 20 Mbps
  SLOW_SPEED = 1;  // Slow speed mode 500 kbps
}

// Request to set the mode of a serial interface
message SetModeRequest {
  string interfaceName = 1;  // e.g., "serial1"
  SerialMode mode = 2;       // Desired mode (RS485, RS422 or RS232)
}

// Response containing available modes of a serial interface
message GetAvailableModesResponse {
  repeated string modes = 1;   // Available modes (RS485 and RS422)
}

// Response containing available slew rates of a serial interface
message GetAvailableSlewRatesResponse {
  repeated string slewRates = 1;   //
 Available slew rates (HIGH_SPEED and SLOW_SPEED)
}

// Request to set the termination of a serial interface
message InterfaceNameRequest {
  string interfaceName = 1;   // e.g., "serial1"
}

// Request to set the termination of a serial interface
message SetTerminationRequest {
  string interfaceName = 1;   // e.g., "serial1"
  bool enableTermination = 2; // true to enable, false to disable
}

// Request to set the slew rate of a serial interface
message SetSlewRateRequest {
  string interfaceName = 1;   // e.g., "serial1"
  SlewRateMode mode = 2;      // Desired mode (HIGH_SPEED or LOW_SPEED)
}

// Request to set the baud rate of a serial interface
message SetBaudRateRequest {
   string interfaceName = 1;   // e.g., "serial1" 
   int32 baudRate = 2; // Desired baud rate
}

// Request to list available serial interfaces
message ListInterfacesRequest {}

// Response containing the list of serial interfaces
message ListInterfacesResponse {
  repeated string interfaces = 1; //
 List of available serial interface names (e.g., ["serial1", "serial2"])
}

// Request to get the status of a serial interface
message GetStatusRequest {
  string interfaceName = 1;  // e.g., "serial1"
}

// Response containing current power status of a serial interface
message GetTransceiverPowerResponse {
  bool powerOn = 1; // true if power is on, false if power is off
}

// Response containing current serial mode of a serial interface
message GetModeResponse {
  SerialMode mode = 1; // Current mode (RS485 or RS422)
}

// Response containing current termination status of a serial interface
message GetTerminationResponse {
  bool terminationEnabled = 1; // true if termination is enabled,
 otherwise false
}

// Response containing current slew rate mode of a serial interface
message GetSlewRateResponse {
  SlewRateMode slewRateMode = 1; // Current slew rate mode
}

// Response containing current baud rate of a serial interface
message GetBaudRateResponse {
  int64 baudRate = 1; // Current baud rate
}

// Response containing current RX count of a serial interface
message GetStatisticsResponse {
  int64 txCount = 1; // Current TX count
  int64 rxCount = 2; // Current RX count
}

// Detailed serial device data
message DetailedSerialDeviceData {
  string name = 1; // Device name
  string devPath = 2; // Device path in /dev directory
  repeated string alternativeDevPaths = 3; // Alternative device paths in /dev/
serial/by-id directory
  string description = 4; //
 Device description (dynamic devices might have generic simple description)
}

// List of serial interfaces with detailed information
message DetailedSerialDevices {
  repeated DetailedSerialDeviceData serialDevices = 1; //
 list of serial devices with detailed information
}

// SerialInterfaceService definition
service Serial {
  // Lists available serial interfaces
  rpc ListInterfaces(ListInterfacesRequest) returns (ListInterfacesResponse);

  // Gets the available Slew Rates of Serial interface
  rpc GetAvailableSlewRates(google.protobuf.
Empty) returns (GetAvailableSlewRatesResponse);

  // Sets the transceiver power (on/off)
  rpc SetTransceiverPower(SetTransceiverPowerRequest) returns (google.protobuf.
Empty);

  // Gets the transceiver power (on/off)
  rpc GetTransceiverPower(InterfaceNameRequest) returns (GetTransceiverPowerRes
ponse);

  // Sets the mode of the serial interface (RS485 or RS422)
  rpc SetMode(SetModeRequest) returns (google.protobuf.Empty);

  // Gets the mode of the serial interface (RS485 or RS422)
  rpc GetMode(InterfaceNameRequest) returns (GetModeResponse);

  // Gets the available modes of the serial interface (RS485 and RS422)
  rpc GetAvailableModes(google.protobuf.
Empty) returns (GetAvailableModesResponse);

  // Sets the termination of the serial interface (on/off)
  rpc SetTermination(SetTerminationRequest) returns (google.protobuf.Empty);

  // Gets the termination of the serial interface (on/off)
  rpc GetTermination(InterfaceNameRequest) returns (GetTerminationResponse);

  // Sets the slew rate mode of the serial interface (high speed or slow speed)
  rpc SetSlewRate(SetSlewRateRequest) returns (google.protobuf.Empty);

  // Gets the slew rate mode of the serial interface (high speed or slow speed)
  rpc GetSlewRate(InterfaceNameRequest) returns (GetSlewRateResponse);

  // Sets the baud rate of the serial interface
  rpc SetBaudRate(SetBaudRateRequest) returns (google.protobuf.Empty);

  // Gets the baud rate of the serial interface
  rpc GetBaudRate(InterfaceNameRequest) returns (GetBaudRateResponse);

  // Gets the baud rate of the serial interface
  rpc GetStatistics(InterfaceNameRequest) returns (GetStatisticsResponse);

  // Lists available serial interfaces with additional information
  rpc ListDetailedInterfaces(google.protobuf.
Empty) returns (DetailedSerialDevices);
}
```

</details>

### RS485 Example

In this section, the RS485 application is implemented.

```c
Download Example Code from embedded file: "modbus-serial-sensor-tdce.zip"
```

#### Application Implementation

The application uses a single '.go' file to run. The 'go.bug.st/serial' package is used to work with the serial port on '/dev/ttymxc1' (SERIAL_2 for RS485). The port mode is set and the communication to the port is opened.

```go
func setupPort() serial.Port {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
	}
	port, err := serial.Open("/dev/ttymxc1", mode)
	if err != nil {
		log.Fatal(err)
	}
	return port
}
```

A message for fetching the temperature and humidity of the sensor is created. For the 'Isolated Soil Sensor', a message in a specific format is sent to the serial port which prompts the sensor to return the required values.

To send a message to the serial port, the 'Write' function is used.

```go
_, err := port.Write(message)
	if err != nil {
		log.Fatalf("Error writing to serial port: %v", err)
	}
```

The program then sleeps for a second so that the sensor has enough time to process the message and send data back to the host. For reading the serial port data, the 'Read' function is used.

```go
response := make([]byte, 256)
	n, err := port.Read(response)
	if err != nil {
		log.Fatalf("Error reading from serial port: %v", err)
	}
```

The result is then printed and parsed to a human-readable format.

```go
func parseValues(rawData []byte) {
	tempHigh := rawData[3]
	tempLow := rawData[4]
	temperatureRaw := (uint16(tempHigh) << 8) | uint16(tempLow)
	realTemperature := float64(temperatureRaw) / 10.0

	humidityHigh := rawData[5]
	humidityLow := rawData[6]
	humidityRaw := (uint16(humidityHigh) << 8) | uint16(humidityLow)
	realHumidity := float64(humidityRaw) / 10.0

	fmt.Printf("Temperature: %.2f °C\n", realTemperature)
	fmt.Printf("Humidity: %.2f %%\n", realHumidity)
}
```

The results are printed to the console.

#### Application Deployment

This section describes the Go application deployment.

**1.2.2.1. Dockerfile**

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```docker
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o modbus-serial .

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/modbus-serial .

CMD ["./modbus-serial"]

```

Open a terminal and paste the following commands:

```bash
docker build -t modbus-serial .
docker save -o modbus-serial.tar modbus-serial:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

**1.2.2.2. Dockerfile Breakdown**

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' file to the directory, then downloads necessary files.

```docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```docker
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0
```

The application image is built as 'modbus-serial'.

```docker
RUN go build -o modbus-serial .
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'serial' application is started.

```docker
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/modbus-serial .

CMD ["./modbus-serial"]
```

**1.2.2.3. Deploying to Portainer**

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). 

As soon as the image and container are set up, the application starts running.


### RS232 / RS422 Example

In this section, the RS232 / RS422 application is discussed. The RS232 and RS422 examples are mostly the same, with the exception of using the '/dev/ttymxc5' port for RS232, while RS422 uses the '/dev/ttymxc1' port.

For clarity, example usage is shown using RS422.

Download the sample applications below:

- RS232: ```c
Download Example RS232 Code from embedded file: "serial-rs232-tdce.zip"
``` Code</a>

- RS422: ```c
Download Example RS422 Code from embedded file: "serial-rs422-tdce.zip"
``` Code</a>

#### Application Implementation

Both applications need to connect to the TDC-E device's port responsible for serial communication in order to read data sent via RS422. To achieve this, a function is created to connect to the '/dev/ttymxc1' (SERIAL_2) port using the 'go.bug.st/serial' Go package.

This function configures the serial connection, setting the parity, data bits (8), stop bits, and baud rate. It then attempts to open the '/dev/ttymxc1' port. If successful, the port is returned; otherwise, an error is logged, and the application is terminated.

```go
func setupPort() serial.Port {
	mode := &serial.Mode{
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		BaudRate: 9600,
	}

	port, err := serial.Open("/dev/ttymxc1", mode)
	if err != nil {
		log.Fatal(err)
	}
	return port
}
```

Then, a simple goroutine is started for both the writing and reading application.

To read the data, the following function is used:

```go
func readData(port serial.Port) {
	buf := make([]byte, 128)
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Printf("Error reading data: %v\n", err)
			return
		}

		if n > 0 {
			fmt.Printf("Received: %s\n", string(buf[:n]))
		}
	}
}
```

The readData function continuously reads from the port, storing the incoming data in a byte buffer. If data is received, it is printed to the console.

To send data through the port, the 'writeData' function is used in the writer application:

```go
func writeData(port serial.Port, message []byte) {
	for {
		_, err := port.Write(message)
		if err != nil {
			log.Printf("Error writing data: %v\n", err)
			return
		} else {
			fmt.Printf("Sent: %s\n", message)
		}
		time.Sleep(2 * time.Second)
	}
}
```

This function accepts the port and a message in byte format. In this example, the application sends a 'Hello world!' message every two seconds. If writing to the port fails, an error message is logged.

#### Application Deployment

This section describes the Go application deployment.

**1.3.2.1. Dockerfile**

To deploy the applications, Go containers should be created and deployed to the TDC-E devices. To that end, a 'Dockerfile' is created for both applications. As 'Dockerfiles' are identical, this documentation will focus on showing a single 'Dockerfile', but to create two applications, make sure to rename the 'serial' tag in the file. The file is shown below.

```docker
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o serial .

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/serial .

CMD ["./serial"]

```

Open a terminal and paste the following commands:

```bash
docker build -t serial .
docker save -o serial.tar serial:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

**1.3.2.2. Dockerfile Breakdown**

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' file to the directory, then downloads necessary files.

```docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```docker
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0
```

The application image is built as 'serial'.

```docker
RUN go build -o serial .
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'serial' application is started.

```docker
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/serial .

CMD ["./serial"]
```

**1.3.2.3. Deploying to Portainer**

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer).

As soon as the image and container are set up, the application starts running.


## Node-RED Examples 

A NodeRed application is provided as a Serial link usage example. Two scripts are given; one writes data to the TDC-E serial port, while the other reads from this port.

For Serial Out the palette 'node-red-node-serialport' is required. Install it over the NodeRed UI.

### Serial Node Example

```c
Download Example Code from embedded file: "serial-inout-tdce.json"
```

#### Write to Serial-Port

Drag a 'Serial Out' node onto your flow.

Double click for settings:
- Add new Serial Port by clicking on the '+' in the Serial Port row
- The Serial Port is on '/dev/ttymxc5' for RS232 on SERIAL_1
- The Serial Port is on '/dev/ttymxc1' for RS422/RS485 on SERIAL_2
- Match the Settings to your Other Serial Device

Deploy the flow. If a " connected" appears under the 'Serial Out' node, NodeRed successfully connected to the serial port.

Add an 'Inject' node.

Edit the 'Inject' node:
- Change payload to String, with a value of "Hello world!"
- Set Repeat to interval every 5 seconds

Connect the two Nodes and deploy your flow.

#### Read from Serial-Port

Drag a 'Serial In' node onto your flow.

Double click for settings:
- Add new Serial Port by clicking on the '+' in the Serial Port row
- The Serial Port is on '/dev/ttymxc5' for RS232 on SERIAL_1
- The Serial Port is on '/dev/ttymxc1' for RS422/RS485 on SERIAL_2
- Match the Settings to your Other Serial Device

Deploy the flow. If a " connected" appears under the 'Serial Out' node, NodeRed successfully connected to the serial port.

Add a 'Debug' node and connect the two Nodes and deploy your flow. Now all incoming serial data gets printed onto you debug window, which you can find in the top right corner.

![Serial Node-RED Example](static/img/serial-node-red.png)

### Serial gRPC Example

```c
Download Example Code from embedded file: "grpc-serial-tdce.json"
```

The 'gRPC' node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For 'gRPC' node installation, import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

The 'gRPC' node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- 'Server':'unix:///var/run/hal/hal.sock'
- a 'Proto File' is provided 

To test out any of the listed functionalities, use the 'inject' node. The result should appear in the 'debug' node. 

![Serial RS422 Node-RED Example](static/img/grpc-serial.png)

## Lua Example

A Lua application is provided as a Serial link usage example. Two scripts are given; one writes data to the TDC-E serial port, while te other reads from this port. For 'RS485' communication, a 'DIGITUS USB to Serial Adapter' is connected to the TDC-E device for reading and writing data. For 'RS422' communication, two TDC-E devices are connected via a serial interface. For 'RS232', loopback is created by connecting the TX and RX wires.

> **⚠️ Warning**
>
> Termination for 'SERIAL_1' and 'SERIAL_2' cannot be set as it is unsupported! Setting termination will result in configuration failure.


To change between the 'RS232', 'RS422' and 'RS485' standards, find the following line of code in the '.lua' example and set it accordingly. For example, to set the 'RS485' standard, use the following:

```lua
S2:setType("RS485")
```

### Writing to the Serial Device

```c
Download Example Code from embedded file: "serial-write-tdce.lua"
```

The first '.lua' script writes data to the serial device. It creates a serial connection to the device. The type of the device is set and baud rate is set to '115200'.

```lua
S2 = SerialCom.create('SER2')
  
S2:setType("RS485")
S2:setBaudRate(115200)
```

A connection is opened and a 'Hello world!' message is created. Then, a timer is created to send this message to the RS485 device periodically every 5 seconds. The message is transmitted the following way:

```lua
local Retb = S2:transmit(message)
print("Transmitted " .. Retb .. " bytes.")
```

Note that a timer is implemented instead of a 'Sleep' service. This is because 'Sleep' will cause all code to wait for the specified time, while timers operate locally, meaning that only the 'rs-write.lua' application script will sleep for the specified time. 

The script prints the engine version at the end of the script. See the result of the application run below.

```
[15:52:06.041: INFO: AppEngine] Starting app: 'serial' (priority: LOW)
Transmitted 13 bytes.
Transmitted 13 bytes.
Transmitted 13 bytes.
```

### Reading from the Serial Device

```c
Download Example Code from embedded file: "serial-read-tdce.lua"
```

The second '.lua' script reads data from the RS485 device. The script creates a serial connection and sets the baudrate to '115200', then opens the connection. A 'Callback' function is created to read from the device.

A 'Callback' function is created to read from the device.

```lua
S2:register('OnReceive', Callback)
```

All received data is then printed to the console.

```lua
function Callback()
  data = S2:receive(1000)
  print(data)
end
```

See the result of the application run below.
```
[16:23:45.843: INFO: AppEngine] Starting app: 'serial' (priority: LOW)
Received data: Hello world!
Received data: Welcome!
```


<!-- Source: docs/intro.md -->


# TDC-E Developers Documentation

This documentation serves as an insight into the TDC-E device's capabilities and programmability, providing extensive descriptions of the services available on the device. It aims to help users explore the TDC-E and what it has to offer. The developer's documetation has the following file structure:


- [Getting Started](category/getting-started)
- [Code Samples](category/code-samples)
- [AppEngine CROWN Documentation](/app-engine-docs/)

## About Sections

### Getting Started

The _Getting Started_ section serves as an introduction to the TDC-E application development and service usage. It covers topics like application environment setup, application deployment or handling services.

### Code Samples

The _Code Samples_ section provides code snippets for working with the TDC-E interfaces. Interfaces have multiple examples and are written in the Go, Node-RED and / or Lua programming languages for example application usage.

### AppEngine CROWN Documentation

The _AppEngine CROWN Documentation_ section offers a comprehensive overview of AppEngine's capabilities. It provides detailed descriptions of the Lua crowns available on the platform, explaining their functions and practical applications.


<!-- Source: docs/code-samples/one-wire-examples.md -->


# One-Wire Examples

In this section, the One-Wire service on the TDC-E device is discussed. Programming examples are given and thoroughly explained.

## Go gRPC Example

```c
Download Example Code from embedded file: "grpc-one-wire.zip"
```

In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that reads status and temperature / data from the device is given.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file is then placed in a separate package, and from it, gRPC Go files are generated using the following commands:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/one-wire-service.proto
```

This installs needed 'Protoc' files, and generates files for the gRPC service.

First, dial options are set to dial a UNIX socket:

```go
dialOptions := []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return net.Dial("unix", socketPath)
	}),
}
```

The application then creates a gRPC client that connects to the service unix socket.

```go
conn, err := grpc.NewClient("passthrough:///unixsocket", dialOptions...)
if err != nil {
	log.Fatalf("failed to connect to Unix socket: %v", err)
}
defer conn.Close()

client := protos.NewOneWireClient(conn)
ctx = metadata.NewOutgoingContext(ctx, metadata.New(nil))
```

A one-wire instance is created. It's passed the created client. The application then enables the one-wire service before getting the status of the one-wire and listing all available devices.

```go
oneWire := wire.NewOneWire(client)

_, err = oneWire.Enable(ctx)
if err != nil {
	log.Printf("Error enabling 1-wire device: %v", err)
	return
}

_, _ = oneWire.GetStatus(ctx)
_, _ = oneWire.ListDevices(ctx)

```

A ticker channel is created, and the application then enters a loop which continuously reads temperatures from all devices. The loop listens for a stop signal via 'ctx.Done()'. On pressing 'Ctrl+C', the applications stops data collection and the one-wire service is disabled.

```go
for {
	select {
	case <-ctx.Done():
		log.Println("Disabling 1-wire service...")
		disableCtx, cancel := context.WithTimeout(context.Background(), 3*time.
Second)
		defer cancel()
		
		if _, err := oneWire.Disable(disableCtx); err != nil {
			log.Printf("Error disabling 1-wire: %v", err)
		}
		return

	case <-ticker.C:
		err := oneWire.ReadTemperatureFromAllDevices(ctx)
		if err != nil {
			log.Printf("Error reading temperatures: %v", err)
		}
	}
}
```

The 'Wire.go' package also contains example functions on how to read device data and temperature data from a single device.

```go
func (t *OneWire) ReadTemperature(ctx context.Context, deviceId string) (*gen.
TemperatureReading, error) {
	req := &gen.ReadTemperatureRequest{
		DeviceId: deviceId,
	}
	res, err := t.Client.ReadTemperature(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("Temperature for device %s: %s\n", deviceId, res.String())
	log.Printf("----------------------------------------\n")

	return res, nil
}

func (t *OneWire) ReadData(ctx context.Context, deviceId string, offset int32,
 length int32) (*gen.ReadDataResponse, error) {
	req := &gen.ReadDataRequest{
		DeviceId: deviceId,
		Offset:   offset,
		Length:   length,
	}
	res, err := t.Client.ReadDeviceData(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("Data from device %s: %s\n", deviceId, res.String())
	log.Printf("----------------------------------------\n")

	return res, nil
}
```

See an example print of the data stream after connecting to a device below:

```bash
2026/01/20 13:09:55 1-Wire service enabled: success:true message:"enabled"
2026/01/20 13:09:55 ----------------------------------------
2026/01/20 13:09:55 1-Wire Status: enabled:true device_count:1 master_path:"/
sys/bus/w1/devices" active_devices:"10-000803675ae2"
2026/01/20 13:09:55 ----------------------------------------
2026/01/20 13:09:55 Available 1-Wire Devices: devices:{device_id:"10-
000803675ae2" family_code:"10" serial_number:"000803675ae2" device_type:
"DS18S20 (Temperature Sensor)" available:true last_seen:"2026-01-20T13:09:
55Z"} 
2026/01/20 13:09:55 ----------------------------------------
2026/01/20 13:09:58 1-Wire Temperatures: readings:{device_id:"10-
000803675ae2" temperature_celsius:85 timestamp:1768914598 valid:
true} conversion_time_ms:879 
2026/01/20 13:09:58 ----------------------------------------
2026/01/20 13:10:00 1-Wire Temperatures: readings:{device_id:"10-
000803675ae2" temperature_celsius:85 timestamp:1768914600 valid:
true} conversion_time_ms:879 
2026/01/20 13:10:00 ----------------------------------------
2026/01/20 13:10:02 1-Wire Temperatures: readings:{device_id:"10-
000803675ae2" temperature_celsius:85 timestamp:1768914602 valid:
true} conversion_time_ms:895 
2026/01/20 13:10:02 ----------------------------------------
^C2026/01/20 13:10:03 Cancel signal received, shutting down...
2026/01/20 13:10:03 Disabling 1-wire service...
2026/01/20 13:10:03 1-Wire service disabled: success:true message:"disabled"
2026/01/20 13:10:03 ----------------------------------------
```

### Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the OneWire service is the following:

```bash
syntax = "proto3";

package hal.onewire;

option go_package = "gitlab.sickcn.net/GBC05/devicebuilder/devices/tdc/
services/1-wire.git/internal/service/gen";

service OneWire {
  rpc Enable(EnableRequest) returns (EnableResponse);
  rpc Disable(DisableRequest) returns (DisableResponse);
  rpc GetStatus(StatusRequest) returns (StatusResponse);
  rpc GetHotPlugStatus(HotPlugStatusRequest) returns (HotPlugStatusResponse);

  rpc ListDevices(ListDevicesRequest) returns (ListDevicesResponse);
  rpc GetDeviceInfo(DeviceInfoRequest) returns (DeviceInfo);

  // Real-time device event streaming via netlink uevents (instant hot-
plug detection)
  rpc WatchDeviceEvents(EventWatchRequest) returns (stream DeviceEvent);

  // Temperature sensor operations
  rpc ReadTemperature(ReadTemperatureRequest) returns (TemperatureReading);
  rpc ReadAllTemperatures(BulkReadRequest) returns (BulkReadResponse);

  // Generic device operations (RFID ROM reading)
  rpc ReadDeviceData(ReadDataRequest) returns (ReadDataResponse);
}

message StatusRequest {}

message StatusResponse {
  bool enabled = 1;
  int32 device_count = 2;
  string master_path = 3;
  repeated string active_devices = 4;
}

message HotPlugStatusRequest {}

message HotPlugStatusResponse {
  // Detection mode: "uevent" (netlink/instant) or "polling" (2s interval)
  string detection_mode = 1;
  // True if netlink uevent watcher is active
  bool uevent_enabled = 2;
  // Fallback detection message (if uevent unavailable)
  string status_message = 3;
}

message EnableRequest {}

message EnableResponse {
  bool success = 1;
  string message = 2;
}

message DisableRequest {}

message DisableResponse {
  bool success = 1;
  string message = 2;
}

message ListDevicesRequest {}

message ListDevicesResponse {
  repeated DeviceInfo devices = 1;
}

message DeviceInfoRequest {
  string device_id = 1;
}

message DeviceInfo {
  string device_id = 1;
  string family_code = 2;
  string serial_number = 3;
  string device_type = 4;
  bool available = 5;
  string last_seen = 6;  // ISO 8601 formatted timestamp (e.g., "2026-01-12T10:
30:45Z")
  map<string, string> properties = 7;
}

message EventWatchRequest {}

message DeviceEvent {
  enum EventType {
    DEVICE_ADDED = 0;
    DEVICE_REMOVED = 1;
    DEVICE_ERROR = 2;
  }
  EventType type = 1;
  DeviceInfo device = 2;
  string message = 3;
  string event_type_name = 4;  // Human-readable event type (e.g.,
 "DEVICE_ADDED")
}

message ReadTemperatureRequest {
  string device_id = 1;
  int32 resolution = 2;
}

message TemperatureReading {
  string device_id = 1;
  double temperature_celsius = 2;
  int64 timestamp = 3;
  bool valid = 4;
  string error = 5;
}

message BulkReadRequest {
  repeated string device_ids = 1;
  int32 resolution = 2;
}

message BulkReadResponse {
  repeated TemperatureReading readings = 1;
  int64 conversion_time_ms = 2;
}

message ReadDataRequest {
  string device_id = 1;
  int32 offset = 2;
  int32 length = 3;
}

message ReadDataResponse {
  bytes data = 1;
  bool success = 2;
  string error = 3;
}
```

</details>

#### gRPC services

To get the one-wire service status, use 'grpcurl', which is an open-source utility for accessing gRPC services via the shell. For help setting up the 'grpcurl' command, refer to [gRPC Usage](/getting-started/grpc-usage).

Use the following line: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -
plaintext <device_ip>:<grpc_server_port> hal.onewire.OneWire.GetStatus
```

The 'token' field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage). The 'device-ip:grpc_server_port' is the TDC-E IP address and the gRPC serving port. For example, if the 'token' value was 'token' and the address and port were '192.168.0.100:8081', you would use the following line to see the status of the one-wire service.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 hal.onewire.OneWire.GetStatus
```

The response should be in this format:

```bash
{
  "enabled": false,
  "deviceCount": 0,
  "masterPath": "/sys/bus/w1/devices",
  "activeDevices": []
}
```

If not enabled, use the following line to do so:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 hal.onewire.OneWire.Enable
{
  "success": true,
  "message": "enabled"
}
```

To list all available devices, use the following line:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 hal.onewire.OneWire.ListDevices
{
  "devices": [
    {
      "deviceId": "10-000803675ae2",
      "familyCode": "10",
      "serialNumber": "000803675ae2",
      "deviceType": "DS18S20 (Temperature Sensor)",
      "available": true,
      "lastSeen": "2026-01-20T13:24:43Z",
      "properties": {}
    }
  ]
}
```

Additionally, you can use the 'gRPC Clicker' VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

### Application Deployment

#### Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```docker
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o onewire-grpc ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/onewire-grpc .

CMD ["./onewire-grpc"]

```

Open a terminal and paste the following commands:

```bash
docker build -t onewire-grpc-app .
docker save -o onewire-grpc-app.tar onewire-grpc-app:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

#### Dockerfile Breakdown

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' and 'go.sum' files to the directory and downloads all needed files.

```docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```docker
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0
```

The application is built as 'onewire-grpc'.

```docker
RUN go build -o onewire-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'onewire-grpc' application is started.

```docker
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/onewire-grpc .

CMD ["./onewire-grpc"]
```

#### Deploying to Portainer

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#2-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```docker
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## Node-RED gRPC Example

```c
Download Example Code from embedded file: "grpc-one-wire.json"
```

The 'gRPC' node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For 'gRPC' node installation, import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

This section describes usage of the one-wire service on the TDC-E. A Node-RED gRPC application was created to demonstrate the following one-wire functionalities:
- enable / disable
- list one-wire devices
- getting status
- getting hotplug status
- getting device information
- reading temperature (all devices, single device)
- reading data
- streaming device events (adding / removing devices)

For implementation, the following nodes are used:
- 'inject' node
- 'gRPC call' node
- 'debug' node.

The 'gRPC' node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- 'Server':'unix:///var/run/hal/hal.sock'
- a 'Proto File' is provided 

To test out any of the listed functionalities, use the 'inject' node. The result should appear in the 'debug' node. 

See a screenshot of the application below:

![One-Wire Node-RED Example](static/img/node-red-one-wire.png)

## Lua Example (Not Supported)

One-Wire interface is **not supported** in Lua.


<!-- Source: docs/getting-started/installing-accessing-removing-applications.md -->


# Installing, Accessing and Removing Applications

This section discusses the setup of applications on the TDC-E device. The TDC-E device has some predefined applications that can be downloaded and used in simple steps. Those applications include:

- 'Grafana'
- 'Influxdb'
- 'Mariadb'
- 'Mosquitto'
- 'Nodered'
- 'Portainer'
- 'Postgresql'
- 'Redis'
- 'Vscode'
- 'Sftpgo'
- 'Loki'
- 'Fluentbit'
- 'Sensor Data Gate (SDG)'
- 'File Browser'

An example of installing 'Portainer' will be provided.

## Example: Setting Up Portainer

This section describes how to install, access and uninstall an application using the 'Portainer' service as an example.

### Installing Services

To set up Portainer, go to the [TDC-E Home Page](http://192.168.0.100/#/), then select the **System** tab on the left panel. Select **Applications**.

The **Applications** page lists services that can be installed to the TDC-E device by simply clicking on the **Install** button. To install Portainer, find the service in the given list and click **Install**.

> **ℹ️ Info**
>
> To install the listed applications, make sure the TDC-E device is connected to a network.


![Installing Applications](static/img/install-applications.png)

The chosen service will be installed shortly. Once the download and installation is complete, the service will automatically be running. 

### Accessing Services

To access the newly installed service, select the 'Open' button. This will take you to the home page of the installed service. For example, [Portainer](https://192.168.0.100:9443/) is located on 'https://192.168.0.100:9443/'.

![Portainer](static/img/portainer-sample.png)

> **ℹ️ Info**
>
> Some services, like Portainer, require setting up a username and password before proceeding to the application itself.


### Uninstalling Services

To uninstall an application from your device, select the **System** tab on the left panel. Select **Applications**. Select **Uninstall** for the application you want to remove and wait for the process to finish. This will remove the application and all user data on it from the device.

## Uploading Custom Applications

A custom application can also be deployed to the device by selecting the 'Upload custom application' button on **System > Applications**. This prompts the user for a '.tar', '.tar.gz', '.tar.bz2' or '.tar.xz' file.

The file should contain the following directories and files:

- '/compose' (required)
- '/images' (required)
- '/configs'
- 'logo.png'
- 'config.json'
- 'VERSION'

### Application File Directory Explanation

- **/compose (required)**

This directory contains at least one **docker-compose.yml** file, including other files neccessary in the same working directory for docker-compose (e.g .env files). The Docker compose container and service names **must be unique** among all applications on one system.

- **/images (required)**

This directory contains tarball files (.tar|.tar.gz|.tar.bz2|.tar.xz) with all the necessary docker images.

- **/configs** (optional)

This directory contains config files that will be placed into installation directory. These files can be refered to by using environment variable '$APP_INSTALLATION_DIR' and subdirectory configs.

See a 'mosquitto.conf' placed into '/configs/mosquitto/config' directory example below:

```yaml
volumes:
  - mosquitto-data:/mosquitto/data
  - mosquitto-log:/mosquitto/log
  - $APP_INSTALLATION_DIR/configs/mosquitto/config/mosquitto.conf:/mosquitto/
config/mosquitto.conf
```

- **logo.png** (optional)

Logo image that will be shown on device web GUI.

- **config.json** (optional)

This is the application configuration file in 'json' format with following schema:

```json
{
  "autoStart": true,
  "uiEndpoint": 18000,
  "name": "My App",
  "description": "Application serves as an example on how to ...",
  "homepage": "www.sick.com"
}
```

The 'autoStart' parameter specifies if application should start immediately after installation.

The 'uiEndpoint' parameter specifies on which network endpoint application GUI will be available.

- **VERSION** (optional)

The textual file that contains version string of the application that will be reported by the API.


<!-- Source: docs/getting-started/docker-authorization-policy-scope.md -->


# Docker Authorization Policy Scope

In this section, the Docker Authz authorization plugin is introduced, and policy restrictions are listed.

## Introduction

To secure the TDC-E device's Docker engine, an authorization plugin is developed according to a security and compliance policy. This policy **restricts** certain operations like elevating privileges, accessing sensitive devices or writing risky configurations. The plugin is implemented as an **OPA (Open Policy Agent)** policy for Docker.

The policy is implemented as a set of rules Docker commands need to adhere to, and blocks any commands that would violate the device's security.

If any of the commands match the deny rules of the plugin, the request is denied and the plugin returns an 'error' message. Otherwise, the Docker API processes the command.

In the following sections, Docker policy restrictions are listed. 

> **📝 Note**
>
> When developing a custom Docker container and image, make sure to adhere to the policy.


## Policy Segments

### Privileged Containers

_Running containers in privileged mode is **denied**._

Example:
```json
"HostConfig": { "Privileged": false }   // Allowed
"HostConfig": { "Privileged": true }    // Denied
```

### Bind Mounts

_Mounting sensitive paths is **denied**. Make sure to use only **allowed bind mounts**!_

**Allowed Bind Paths:**

- '/datafs/operator'
- '/var/run/iolink'
- '/var/run/cc', '/run/cc'
- '/var/run/hal', '/run/hal'
- '/var/run/docker.sock'
- '/var/volatile/tdcx'
- '/etc/tdcx'
- '/etc/os-release'
- '/media'
- '/datafs/tdc-engine'
- '/etc/machine-id'
- '/run/log/journal'
- '/usr/lib/os-release'
- '/dev/mmcblk1', '/dev/mmcblk1p*'
- '/datafs/appruntime'

**Allowed Read-Only Bind Paths:**

- '/etc/tdcx'
- '/etc/os-release'
- '/etc/machine-id'
- '/run/log/journal'
- '/usr/lib/os-release'

> **📝 Note**
>
> All other bind mounts are **denied**.


Example:
```json
"BindMounts": [
    { "Source": "/etc/tdcx", "ReadOnly": true }    // Allowed
    { "Source": "/etc/tdcx", "ReadOnly": false }   // Denied
    { "Source": "/root", "ReadOnly": true }        // Denied
]
```

### Volume Creation

_Volume creation is subjected to the same **bind mount restrictions as containers**._

**Denied volumes:**

- Bind mounts to **non-whitelisted** or **non-read-only** sensitive paths

### Dangerous Capabilities

_Adding dangerous Linux capabilities is **denied**._

**Denied capabilities:**

- 'all'
- 'audit_control'
- 'sys_admin'
- 'sys_module'
- 'sys_ptrace'
- 'syslog'
- 'dac_read_search'

Other capabilities are considered non-dangerous and **allowed**.

Example:
```json
"HostConfig": { "CapAdd": ["NET_ADMIN"] }      // Allowed
"HostConfig": { "CapAdd": ["SYS_ADMIN"] }      // Denied
```

### Host Namespaces

_Using host namespaces is **denied**._

Therefore, the following are **denied**:

- 'UsernsMode: "host"'
- 'PidMode: "host"'
- 'UTSMode: "host"'
- 'CgroupnsMode: "host"'
- 'IPCMode: "host"'

Example:
```json
"HostConfig": { "UsernsMode": "host" }   // Denied
"HostConfig": { "PidMode": "host" }      // Denied
```

### Publishing All Ports

_Publishing all ports is **denied**._

Example:
```json
"HostConfig": { "PublishAllPorts": true }   // Denied
```

### Security Options

_Unconfined Seccomp, AppArmor, or disabled SELinux labeling is **denied**._

**Denied options:**

- 'seccomp=unconfined'
- 'apparmor=unconfined'
- 'label:disable'

Example:
```json
"HostConfig": { "SecurityOpt": ["seccomp=unconfined"] }   // Denied
```

### Sensitive Sysctls

_Sensitive sysctls are **denied**._

The following are **denied**:
- 'kernel.core_pattern'
- any sysctl starting with 'net.'

Example:
```json
"HostConfig": { "Sysctls": { "kernel.core_pattern": "..." } }   // Denied
"HostConfig": { "Sysctls": { "net.ipv4.ip_forward": "1" } }     // Denied
```

### Device Access

_Access to host devices is **restricted**._

Only the following devices are **allowed:**

- '/dev/serial/by-id/'
- '/dev/mmcblk1', '/dev/mmcblk1p*'

> **📝 Note**
>
> All other devices are **denied**.
> 'DeviceCgroupRules' and 'DeviceRequests' usage is **denied**.


Examples:
```json
"HostConfig": {
    "Devices": [{ "PathOnHost": "/dev/serial/by-id/tdcx-serial" }]   // Allowed
    "Devices": [{ "PathOnHost": "/dev/sda" }]                        // Denied
    "DeviceCgroupRules": ["c 42:* rmw"]                              // Denied
    "DeviceRequests": [{ ... }]                                      // Denied
}
```

### Copying Content into Containers

_Copying content from the host into containers using the archive API is **denied**._

**Denied PUT requests:**
- PUT requests to '/containers/{id}/archive'

If any deny rule listed above matches, the Docker command will fail.

## FAQ

**Q | In some cases my docker container does not start. Looking at logs following error messages can be seen: 'Failed to allocate and map portaddress already in useError starting userland proxy'.**

A | Ensure that no other Docker containers are using the same port. 

If no conflicts are found, the issue may be related to a known Docker Engine issue. Restarting the device should resolve the problem and allow the container to start normally.

---

**Q | I've made changes to TDC-E SSH work environment and now I have troubles working with it.**

A | To reset SSH workspace issue the following operation: 'container-reset-userworkspace' via REST API '/api/v1/system/administration/operation'.

Make sure to confirm the operation with '/api/v1/system/administration/operation/confirm'. Access the Control Center API UI via the page Resources → Documentation → Control Center.

---

**Q | I've made changes to AppEngine either by deploying apps or manipulating its variables and now I have troubles working with it.**

A | To reset AppEngine issue an operation 'container-reset-appengine' via REST API '/api/v1/system/administration/operation'.

Make sure to confirm the operation with '/api/v1/system/administration/operation/confirm'. Access the Control Center API UI via the page Resources → Documentation → Control Center.

---

**Q | After updating to TDC-E FW 1.4.0 it's no longer possible to reach IO-Link REST API on host port 9000 from inside of container.**

A | From version 1.4.0 IO-Link API is available on localhost endpoint '127.0.0.1:19005'. If your container is using 'network=host' then this API is reachable via '127.0.0.1:19005'.

If your container is using a 'bridge network', then add a parameter '--add-host' ("extra_host" in docker compose) '"host.docker.internal:host-gateway"'. 
This way, IO-Link API will be reachable on endpoint 'host.docker.internal:19005'.

---

**Q | After updating to TDC-E FW 1.4.0 I cannot deploy my docker applications anymore. Error mentions 'authorization denied by plugin'.**

A | From version 1.4.0 Docker API authorization is activated which enforces a policy for every Docker API operation.


<!-- Source: docs/getting-started/working-with-docker.md -->


# Working With Docker

This page offers insights into using Docker services on this device, including how to work with Portainer and how to set up access to host services from within containers.

Docker is a platform that simplifies building, sharing, running, and packaging applications and their dependencies into lightweight, portable containers. These containers can run consistently across different environments.

Learn more about Docker [here](https://www.docker.com/).

This device supports Docker for both development and deployment purposes.

> **ℹ️ Info**
>
> Make sure to refer to [Docker Authorization Policy Scope](/getting-started/docker-authorization-policy-scope) when working with Docker. Created containers also have to adhere to the policy.


## Creating Docker Containers

In this section, creating Docker containers is discussed. This includes information about providing access to host devices, mounting volumes, exposing ports, and specifying capabilities.

A Docker container is a **lightweight, isolated runtime environment** that packages an application with everything it needs to run - code, libraries, system tools, configuration, and environment settings. Because all dependencies are included, the application **runs the same way on any system** that supports Docker.

Containers are created from **Docker images**, which typically start with a small **base image** (e.g., Alpine, Debian, or an application‑specific base like Nginx). Additional layers add the application and its dependencies. When the Docker Engine runs the image, it creates a container as a set of isolated processes on the host system, using a layered filesystem and sharing the host’s kernel.

Docker containers can be created using either the **Docker CLI**, **Docker compose**, or **container management tools** (e.g. Portainer). When creating a container, the following needs to be decided:

- Which image the container should run

- How the container should start and behave

- Which host resources (network, storage, devices) it can access

### Basic Container Creation

The simplest way to start a container is using the 'docker run' command:

```bash
docker run IMAGE
```

For example, the following CLI command will run the [alpine](https://hub.docker.com/_/alpine) Docker image.

```bash
docker run -it alpine sh
```

This will:

- Pull the (latest) Alpine image from online registry if it is not present locally

- Create a container from the image

- Start the container with default settings

- Open an interactive shell inside the container

> **ℹ️ Info**
>
> The default online registry for pulling images is [Docker Hub](https://hub.docker.com/). Pulling images requires an **internet connection**.


### Commonly Used Options

In practice, containers are usually started with additional options to configure their runtime environment. See some examples below.

**Run in background and give the container a name:**

```bash
docker run -d --name my-alpine alpine sh
```

This runs the 'alpine' image in the background '(-d option, short for daemon)', and assigns 'my-alpine' as the name of the container.

**Mount a host directory into the container:**

By default, a Docker container has its own isolated filesystem and cannot access files on the host. However, applications often need to read host‑provided files or persist data outside the container. Docker provides two options:

- Volumes – storage handled by Docker
- Bind mounts – direct host‑to‑container directory mappings

Both allow shared access to directories between host and container. This ensures that data created inside the container persists after shutdown, as it is stored on the host rather than in the container’s temporary filesystem.

```bash
docker run -it \
  -v /datafs/operator:/home/root/data \
  alpine sh
```

> **📝 Note**
>
> See [Device-specific Paths](/getting-started/working-with-docker#13-device-specific-paths) below for mapping correct paths to your container.


This starts the alpine container in an interactive shell with data from **host** location '/datafs/operator' mounted to **container** location '/home/root/data'. Files created in the container directory '/home/root/data' will persist on the host.

**Expose a network port:**

```bash
docker run -d \
  --name my-nginx \
  -p 8080:80 \
  nginx
```

This starts an Nginx web server and exposes it on the host at port 8080, forwarding traffic to port 80 inside the container.

**Grant access to a host device or additional capabilities:**

```bash
docker run -it \
  --device /dev/serial/by-id/tdcx-serial:/dev/ttyUSB0 \
  --cap-add NET_ADMIN \
  alpine sh
```

This runs the Alpine shell with:

- the 'NET_ADMIN' capability (allowing network configuration)
- the **host** device '/dev/serial/by-id/tdcx-serial' mapped to the **container** as '/dev/ttyUSB0'.

> **📝 Note**
>
> Access to host devices, privileged mode, and additional capabilities may be restricted by the device’s Docker authorization policy.
>
> See [Device-specific Paths](/getting-started/working-with-docker#13-device-specific-paths) below for mapping correct paths to your container.


**Docker Compose**

The same options can be applied to containers using **Docker Compose** and writing a 'docker-compose.yml' file. Compose provides a declarative and reusable way to define container setup.

For more information go to [Docker Compose Docs](https://docs.docker.com/compose/).

### Device-specific Paths

Some host paths and devices (e.g. SERIAL devices) are exposed at **device-specific locations**.

When mapping volumes or granting device access to a container, always use the paths that correspond to your device.

See a list of **host devices and mount points** for adding to your container below:

- '/dev/ttymxc5' - SERIAL 1 (RS232)
- '/dev/ttymxc1' - SERIAL 2 (RS422/485)
- '/dev/ttymxc6' - GNSS
- '/sys/bus/w1/devices' - 1-WIRE
- '/media/' - external storage

## Host Service Setup

When running Docker containers, you may need to access services running on the host machine (e.g. a containerized application accessing databases or APIs running on the host). Docker provides a special **DNS name** for this purpose inside containers:

```json
host.docker.internal
```

This hostname resolves to the internal IP address of the host machine **from inside the container**, thus allowing containers to communicate with host services without hardcoding IP addresses.

### Example Usage

A 'host.docker.internal' usage is demonstrated by starting a 'Python3' HTTP server, and running an 'alpine' container on the TDC-E device.

A simple HTML page titled 'hello.html' is created. In the same directory, a Python server is started:

```bash
python3 -m http.server <port>
```

This should yield the following output:

```bash
python3 -m http.server 8000
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...
```

Start an alpine container with the 'host.docker.internal' host, and map the gateway dynamically using 'host-gateway'. Add the 'curl' command, then curl the host server for the HTML file.

```bash
docker run -it --rm --add-host=host.docker.internal:host-gateway alpine sh
apk add curl

curl http://host.docker.internal:8000/hello.html
```

The result is shown below.

```html
<!-- hello.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Hello</title>
</head>
<body>
  <h1>Hello, World! </h1>
  <p>This is served from a Python3 HTTP server.</p>
</body>
</html>
```

## Working with Portainer

In this section, working with Portainer is discussed. The section provides an insight into accessing the Portainer service after application installation, deploying images, and building a container for an image.

Portainer is a lightweight web-based GUI for managing Docker. It lets you control containers, images, networks, and volumes without using Docker CLI commands. Portainer itself runs as a Docker container and is given access to the Docker Engine so it can manage other containers on the host.

### Accessing Portainer

Provided you have installed [Portainer](https://docs.portainer.io/) on the TDC-X device, it can be accessed on [this page](https://192.168.0.100:9443/).

For help installing Portainer on your device, refer to [Installing, Accessing and Removing Applications](/getting-started/installing-accessing-removing-applications).

### Deploying an Image on Portainer

To be able to deploy an application on the TDC-E device, the application needs an environment it can run in. This environment is called an 'image'. 

To upload an image, an image file is needed. There are multiple ways of obtaining an image. In this example, we build an application using the 'docker build' command, and the image is saved using 'docker save'.

> **ℹ️ Info**
>
> Check the system architecture before building your application. If the target architecture differs from your development environment, make sure to cross-compile accordingly.


```bash
docker build -t img-tag /path/to/dockerfile
docker save -o output.tar img-tag:latest
```

Replace the 'tag', 'path/to/dockerfile' and 'output' parameters accordingly.

After building the image, go to the [Portainer Dashboard](https://192.168.0.100:9443/#!/2/docker/dashboard) and select the 'Images' option, or find the [Images](https://192.168.0.100:9443/#!/2/docker/images) option on the left panel.

The '.tar' file can be uploaded by selecting the **Import** option.

![Importing an Image](static/img/port-1-image-tar.png)

Select a valid '.tar' file and give the image a fitting name. In this example, a 'dio-grpc-app.tar' file was selected, and the image was named 'dio-grpc'.

![Uploading an Image](static/img/portainer-2-upload-img.png)

Select **Upload** and wait for the image to be uploaded to your device. Once it is on the device, Portainer will show the image and Dockerfile details which specify the ID, size, creation date, build, environment, command and layers of the image. 

### Creating a Container from an Image

An image is needed for a container to be run. In other words, an image cannot do anything on its own as it needs a container to run. To provide a container to an imported image, go to the [Containers](https://192.168.0.100:9443/#!/2/docker/containers) tab on Portainer, where all current containers are listed.

![Containers](static/img/port-3-build-container.png)

To add a container, select the option **Add container**. Provide a name for you container, and the image uploaded to Portainer in the last step.

![Setting Up a Container](static/img/port-4-set-con.png)

One can also set the following parameters:

- Registry
- Always pull the image
- Create a container webhook
- Publish all exposed ports to random host ports
- Port Mapping
- Access control
- Auto Remove

You can also set **Advanced container settings**, which allows the user to set additional environment parameters like:

- Commands and logging
- Volumes
- Network
- Env
- Labels
- Restart policy
- Runtime and resources
- Capabilities

Once all options are set, select **Deploy the container**. The container will now start running.


<!-- Source: docs/getting-started/grpc-usage.md -->


# gRPC Usage

This section provides insight into how to use the HAL Service gRPC API. Setting up 'grpcurl', server authorization, 'gRPC Clicker' and accessing HAL services is discussed.

## 'grpcurl' Installation

The 'grpcurl' utility provides a means to access gRPC services via shell command. It is used to send a request to the TDC-E device with the specified HAL service port. To get the latest version of the 'grpcurl' utility, run the following line in the terminal:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

To be able to access the utility, the 'bin path' needs to be exported.

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

To test 'grpcurl' availability, use the following line:

```bash
grpcurl --version
```

## Authentication

To access TDC-E' HAL Services, user authentication is needed. This is done by fetching an authentication token from the device. Use the 'curl' command to authenticate the TDC-E user. See an example below.

```bash
curl -s -X POST http://{DEVICE_IP}/auth/login -H 'Accept: application/json' -
H 'Content-Type: application/json' -d '{"username":"{USERNAME}","password":
"{PASSWORD}","realm":"{REALM}"}'
```

Replace the '{DEVICE_IP}', '{USERNAME}', '{PASSWORD}' and '{REALM}' with your actual values. For example, if the values of the listed variables are the following:

<GrpcAuthentication />

The 'curl' command will be the following:

<GrpcCurlExamle />

This will return a token response from the server which will be used for further authentication. 

> **📝 Note**
>
> To export the value of your token, you can install a service like 'jq'. See the example of exporting a token with 'jq' below.


```bash
export CC_TOKEN=$(curl -s -X POST http://{DEVICE_IP}/auth/login -H 'Accept:
 application/json' -H 'Content-Type: application/json' -d '{"username":
"{USERNAME}","password":"{PASSWORD}","realm":"{REALM}"}' | jq -r '.token')
```

> **ℹ️ Info**
>
> The token is reset after its expiration time is out, and after each reboot of the device.


## Working with 'grpcurl'

To send gRPC requests to the TDC-E device, an authentication token is needed which is added to the gRPC curl. In other words, each call to the server needs to be authenticated. This is done by adding a header to the 'grpcurl' command.

In the following examples, the '{DEVICE_IP}' and '{GRPC_SERVER_PORT}' will be set to '192.168.0.100' and '8081' respectively. The HAL Service port is set to '8081' by default, but can be configured in the **System Server Settings**.

To list all HAL services via the server reflection, use the following command:

```bash
grpcurl -H 'Authorization: Bearer {token}' -plaintext {DEVICE_IP}:
{GRPC_SERVER_PORT} list
```

See the example below:

```
grpcurl -H 'Authorization: Bearer {token}' -plaintext 192.168.0.100:8081 list
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
hal.boardsupport.BoardSupport
hal.can.Can
hal.digitalIO.DigitalIO
hal.temperaturesensor.TemperatureSensor
```

To list all functionalities of a single service, the following command is used:

```bash
grpcurl -H 'Authorization: Bearer {token}' -plaintext {DEVICE_IP}:
{GRPC_SERVER_PORT} list {PACKAGE}.{SERVICE}
```

For example, we want to list all functionalities of the 'DigitalIO' service from the 'hal.digitalio' package. This is done the following way:

```json
hal.digitalio.DigitalIO.Attach
hal.digitalio.DigitalIO.ListDevices
hal.digitalio.DigitalIO.Read
hal.digitalio.DigitalIO.SetDirection
```

The utility can also be used to access the listed services. For example, one can list all DigitalIO devices using the following command:

```json
{
"devices": [
{
"name": "DIO_F",
"type": "BIDIRECTIONAL",
"direction": "IN"
},
{
"name": "LP_A",
"type": "OUTPUT",
"direction": "OUT"
},
{
"name": "DIO_B",
"type": "BIDIRECTIONAL",
"direction": "IN"
},
{
"name": "DIO_A",
"type": "BIDIRECTIONAL",
"direction": "IN"
},
{
"name": "LP_B",
"type": "OUTPUT",
"direction": "OUT"
},
{
"name": "DIO_C",
"type": "BIDIRECTIONAL",
"direction": "IN"
},
{
"name": "DIO_D",
"type": "BIDIRECTIONAL",
"direction": "IN"
},
{
"name": "DIO_E",
"type": "BIDIRECTIONAL",
"direction": "IN"
}
]
}
```

To read the value of a digitalIO device, the following command can be used:

```json
{
"state": "LOW"
```

To write the direction of a DIO device, use the following command:

Reading the state of 'DIO_A' now will yield the following result:

```json
"state": "HIGH"
```

## gRPC Clicker

The 'gRPC Clicker' extension is a VSCode extension which can be installed via the _Extensions Marketplace_. It is used to connect to a gRPC server and to request services from it. It has a Graphical User Interface to make it easier for the user to interact with the server.

### Installing gRPC Clicker

To install gRPC Clicker, open VSCode and go to the **Extensions** tab. Search for _gRPC Clicker_. Click _Install_.

![gRPC Clicker Extension](static/img/grpc-clicker.PNG)
 
After installation, the tool should create a new tab which can be accessed. To add a new gRPC server, select the **plus** icon in the **Proto Groups** section. A **Proto schema assistant** will open.

The following parameters should be set up:

- 'Name' - Name of your Proto Group
- 'Address' - Device IP and gRPC server port
- 'Custom flags' - Authorization header

For example, set the name to grpc-TDC-E, and the address and port to {selectByDevice({ SIM2000: "192.168.1.1:8081", default: "192.168.0.100:8081" })}. Custom flags need to be set as authentication is needed to request services from the gRPC TDC-E server. The custom flag has the following structure:

```h
-H 'Authorization: Bearer TOKEN'
```

> **📝 Note**
>
> Make sure to generate an authentication token and paste it instead of the 'TOKEN' value!


The group is now set and can list all gRPC services currently available on your TDC-E device. 


### Working with gRPC Clicker

To list all DIO devices, for example, click on the 'hal.digitalio.DigitalIO' subgroup and expand it. The left panel is used for setting parameters that are needed for the request, like when reading or writing a value. The right panel returns the server result.

Select 'ListDevices', which doesn't require any input parameters. You can send the request by selecting the **Send** button, and a response will be generated on the right side of the screen, listing all available devices.

![DigitalIO interfaces list](static/img/grpc-list-dio.jpg)

To read the state of a device, select the 'Read' subgroup and insert the name of the DIO device to read from on the left panel.

![Reading DigitalIO state](static/img/grpc-read-dio.jpg)

To write the direction of a DIO device, select 'SetDirection'. Set the 'name' and 'direction' variables.

![Writing DigitalIO direction](static/img/grpc-write-dio.jpg)

This sets the direction of the DigitalIO, and setting the device state to 'HIGH' via the 'Write' call will now be enabled.

## Accessing HAL Service from Unix Socket

For internal usage, the TDC-E has a Unix socket which allows accessing available HAL Services. The HAL socket on the TDC-E device works with a TCP endpoint like 'grpcurl', so the command is generally used in the same way. 

> **📝 Note**
>
> Install 'grpcurl' in your workspace to access socket content or use the examples provided in the **Code Samples** section.


For example, listing all HAL service commands is done using the following line:

```bash
grpcurl -emit-defaults -plaintext -unix -authority "localhost" "unix:///var/
run/hal/hal.sock" list
```

Instead of a server address, 'grpcurl' now uses the '-unix' flag and the path to the HAL socket. Since the HAL socket is run locally on the TDC-E device, no authentication form is required.

## Additional Notes

For Node-RED examples used in this documentation, a **custom version** of the [node-red-contrib-grpc node](https://flows.nodered.org/node/node-red-contrib-grpc) was implemented. Using the Palette's version in Node-RED will **not allow TDC-E user identification and/or accessing the UNIX socket**! 

Make sure to import the following node into your Node-RED project:

```c
node-red-contrib-grpc v1.2.7 from embedded file: "node-red-contrib-grpc-1.2.7.
tgz"
```


<!-- Source: docs/getting-started/networking.md -->

---
title: Networking
sidebar_position: 104
---



# Networking

## How to Setup DHCP or Static IP

Most industrial computers come with predefined networking settings.

TDC-E comes with two predefined interfaces: 'eth1' and 'eth2'. 'eth2' and 'wlan0' are set to DHCP; 'eth1' is Static.

```
IP address: 192.168.0.100
Subnet mask: 255.255.255.0
Default gateway: 0.0.0.0
```

### Interfaces

- Ethernet 1 - 'eth1'
- Ethernet 2 - 'eth2'
- WLAN - 'wlan0'

![TDC-E Interfaces](static/img/network-image-005.png)

Interfaces mode can be 'Static' or 'DHCP'.

Every inteface has fallback modes that can be set.

See: [DHCP section](#13-dhcp) 


### Static Mode

To change the IP address of a device, change the corresponding values in the shown fields and click apply:

![Networking static mode](static/img/network-image-001.png)

--- 

For example, if the IP pool of your router’s DHCP server is '192.168.0.2' - '192.168.0.254'​, you should enter an IP address within this range, like '192.168.0.75'.

![Networking change IP address](static/img/network-image-002.png)

> **📝 Note**
>
> The IP address of a device is now changed, and the adapter settings on your local PC must be changed accordingly.


### DHCP

To change default 'static' mode to 'DHCP', click in the field  'Addressing mode' and change option from static to DHCP.

![Networking change mode](static/img/network-image-003.png)

When DHCP is selected, every interface has a fallback mode that can be a static IP address or repeated DHCP.

The device will try to get interface networking data from the available DHCP server.

In case DHCP data is not available, the device will set its interface networking data to a statically defined IP address.

![Networking dhcp mode](static/img/network-image-004.png)

> **📝 Note**
>
> The IP address of a device is now changed, and the adapter settings on your local PC must be changed accordingly.


## WLAN Connectivity

In this section, wireless network interface is discussed.

![Wireless network interface](static/img/network-image-006.png)

> **ℹ️ Info**
>
> In order to enable wireless network interface country code must be set.


### Add New Wireless Network

Make sure WLAN interface is enabled, then click on a plus button to add a new network connection.

Enter network credentials and confirm.

![Wireless add new connection](static/img/network-image-007.png)

When new network is connected, status will be "Connected".

![Wireless add connection status](static/img/network-image-008.png)

### Scan for Wireless Networks

To view scanned WLAN networks periodically, click on the 'Periodic scan' button.

> **⚠️ Disclaimer: WLAN Performance Impact**
>
> Periodic scanning remains active while the **WLAN page is open**, which may **affect system performance**.



## Gateway Priority

Gateway priority card shows list of available interfaces and current gateway priority order.

![Gateway priority](static/img/network-image-009.png)

To change the priority, select only one interface item and reorder the selected item.

When order is set, confirm by clicking apply.

> **ℹ️ Info**
>
> When no gateway is set for the current interface in the 'Interfaces' menu, the defined order for this interface is not taken into account.


![Gateway priority - change order](static/img/network-image-012.png)


### DNS

The currently used DNS servers list shows a combination of user-defined DNS servers and all DHCP predefined server lists.

The user-defined DNS server list takes precedence.

Creation order of user defined DNS servers defines the priority of used DNS servers.

![Gateway - DNS](static/img/network-image-010.png)

## Modem Connectivity

Modem and modem connection must be enabled.
When modem connection is disabled, modem can be used for sending and receiving SMS data.

To establish a modem internet connection, APN profile data must be set.

APN (Access Point Name of your current provider) is a mandatory value and must be set.

Depending on your provider, username and/or password are not required.

The dial string field is used in some special cases where modem initialization is required by the carrier.

![Modem connectivity](static/img/network-image-011.png)

> **📝 Note**
>
> In case where multiple interfaces are utilized and the modem interface has internet access, DNS priority must be set accordingly.


## Configuration of Firewall Rules

Firewall section enables you to set all interfaces as private or public.

Public interfaces can be accessed from the Internet.

![Firewall](static/img/network-image-013.png)


<!-- Source: docs/code-samples/dio-examples.md -->


# Digital Input / Output Examples

In this section, the Digital Inputs/Outputs of the TDC-E device are discussed. Programming examples are given and thoroughly explained.

## Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. From it, the list of DIO devices, their input / output directions and states can be read and set. To that end, a simple Go application demonstrating said usage is created.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/dio.proto
```

This installs needed tools, and generates files for the gRPC service.

First, dial options are set to dial a UNIX socket:

```go
dialOptions := []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return net.Dial("unix", socketPath)
	}),
}
```

The application then creates a gRPC client that connects to the service unix socket.

```go
conn, err := grpc.NewClient("passthrough:///unixsocket", dialOptions...)
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
defer conn.Close()

client := pb.NewDigitalIOClient(conn)
ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(nil))
```

A DigitalIO instance is created which is used to access DigitalIO data. It's passed the created client and context, alongside the requested DIO device. All DIOs are listed, and the requested DIO is read repeately, and state is changed.

```go
digitalIO := digitalio.NewDigitalIO(DIODevice, client, ctx)
digitalIO.ListAllDIOs()

for {
	state := digitalIO.ReadDio()
	digitalIO.ChangeDioState(state)
	time.Sleep(2 * time.Second)
}
```

Reading a DigitalIO is done like this:

```go
func (t *DigitalIO) ReadDio() string {
	req := &protos.DigitalIOReadRequest{
		Name: t.Device,
	}
	res, err := t.Client.Read(t.Ctx, req)
	if err != nil {
		log.Fatalf("could not read: %v", err)
	}
	state := res.State.String()
	fmt.Printf("The state of %s is: %s\n", req.Name, state)
	return state
}
```

The application sleeps for 2 seconds after every read and state change.

### Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the DigitalIO service is the following:

```bash
/**
 * DigitalIO Service.
 *
 * Service that enables a control of the Digital IO features of device
 */
syntax = "proto3";

package hal.digitalio;

import "google/protobuf/empty.proto";

option go_package = "./protos;protos";

enum IOType {
  INPUT = 0;
  OUTPUT = 1;
  BIDIRECTIONAL = 2;
}

enum IODirection {
  IN = 0;
  OUT = 1;
}

enum IOState {
  ERROR = 0;
  LOW = 1;
  HIGH = 2;
}

message IODevice {
  string name = 1;
  IOType type = 2;
  IODirection direction = 3;
}

/**
 * Represents the ListDevices response data.
 */
message DigitalIOListDeviceResponse {
  repeated IODevice devices = 1;
}

/**
 * Represents the SetDirection request data.
 */
message DigitalIOSetDirectionRequest {
  string name = 1;
  IODirection direction = 2;
}

/**
 * Represents the Read request data.
 */
message DigitalIOReadRequest {
  string name = 1;
}

/**
 * Represents the Read response data.
 */
message DigitalIOReadResponse {
  IOState state = 1;
}

/**
 * Represents the Write request data.
 */
message DigitalIOWriteRequest {
  string name = 1;
  IOState state = 2;
}

/**
 * Represents the Attach response data.
 */
message DigitalIOAttachResponse {
  string name = 1;
  IOState state = 2;
  int32 error = 3;
  string timestamp = 4;
}

/**
 * Service exposing DigitalIO functions.
 */
service DigitalIO {
  /// Used to retrieve all available digital IO devices.
  rpc ListDevices(google.protobuf.
Empty) returns (DigitalIOListDeviceResponse) {}

  /// Used to set pin direction value of a particular gpio pin.
  rpc SetDirection(DigitalIOSetDirectionRequest) returns (google.protobuf.
Empty) {}

  /// Used to read value of a particular gpio pin.
  rpc Read(DigitalIOReadRequest) returns (DigitalIOReadResponse) {}

  /// Used to write value of a particular gpio pin.
  rpc Write(DigitalIOWriteRequest) returns (google.protobuf.Empty) {}

  /// Used to stream input events form device.
  rpc Attach(google.protobuf.Empty) returns (stream DigitalIOAttachResponse) {}
}
```

</details>

To list all available Digital Input/Output services, use 'grpcurl', which is an open-source utility for accessing gRPC services via the shell. For help setting up the 'grpcurl' command, see [gRPC Usage](/getting-started/grpc-usage).

To list all available Digital Input/Output calls, use the following line:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -
plaintext <device_ip>:<grpc_server_port> list hal.digitalio.DigitalIO
```

The 'token' field is the fetched TDC-E authorization token. For help fetching this token, refer to [gRPC Usage](/getting-started/grpc-usage). The 'device-ip:grpc_server_port' is the TDC-E IP address and the gRPC serving port. For example, if the 'token' value was 'token' and the address and port were '192.168.0.100:8081', you would use the following line to list all available Digital Input/Output services.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 list hal.digitalio.DigitalIO
```

The response should be in this format:

```log
hal.digitalio.DigitalIO.Attach
hal.digitalio.DigitalIO.ListDevices
hal.digitalio.DigitalIO.Read
hal.digitalio.DigitalIO.SetDirection
hal.digitalio.DigitalIO.Write
```

Additionally, you can use the 'gRPC Clicker' VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

### Application Deployment

#### Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```dockerfile
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o dio-grpc ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dio-grpc .

CMD ["./dio-grpc"]
```

Open a terminal and paste the following commands:

```bash
docker build -t dio-grpc-app .
docker save -o dio-grpc-app.tar dio-grpc-app:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

#### Dockerfile Breakdown

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' and 'go.sum' files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```dockerfile
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
CGO_ENABLED=0
```

The application is built as 'dio-grpc'.

```dockerfile
RUN go build -o dio-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'dio-grpc' application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dio-grpc .

CMD ["./dio-grpc"]
```

#### Deploying to Portainer

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## Node-RED gRPC Example

```c
Download Example Code from embedded file: "grpc-dio.json"
```

The 'gRPC' node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For 'gRPC' node installation, import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

This section describes the creation and usage of a Node-RED gRPC example using digital input / output devices. The application makes listing all DIO devices, reading DIO states, writing the DIO state and streaming changes in the DIO devices possible. 

For implementation, the following nodes are used:

- 'inject' node
- 'gRPC call' node
- 'debug' node.

The 'gRPC' node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- 'Server':'unix:///var/run/hal/hal.sock'
- a 'Proto File' is provided 

To test out any of the listed functionalities, use the 'inject' node. The result should appear in the 'debug' node. 

See a screenshot of the application below:

![DIO Node-RED Example](static/img/dio-node-red.png)

## Lua Example

```c
Download Example Code from embedded file: "dio.lua"
```

A Lua application is provided as a DIO usage example. The example demonstrates setting a DIO device as output or input, setting an output state and reading DI states.

The script prints the engine version at the start of the script. It then creates a 'DIO AO' output and sets the value of the output to 'HIGH'.

```lua
dioAO = Connector.DigitalOut.create('DIO_AO')
dioAO:set(true);
print("Set DO A output to HIGH.")
```
To set a DO output to 'LOW', uncomment the following lines:

```lua
dioAO:set(false);
print("Set DO A output to LOW.")
```

The script then sets three DIO devices to inputs ('DIO_BI', 'DIO_CI' and 'DIO_DI'). The state of said inputs is then read using the following function:

```lua
function PrintDIOStateIn(dio, name)
  local state = dio:get()
  print(string.format("%s state: %s", name, state))
end
```

Another way of reading the state is registering the DI 'onChangeStamped', printing the DI state on event.

```lua
dioDI = Connector.DigitalIn.create('DIO_DI')
function Print_DIODI()
    local state = dioDI:get()
    print(string.format("State of DIO_DI is: %s", state))
end

Connector.DigitalIn.register(dioDI, "OnChangeStamped", Print_DIODI)
```

See the result of the application run below.

```log
[15:55:53.721: INFO: AppEngine] Starting app: 'dio' (priority: LOW)
Engine version: 0.8.3
Set DO A output to HIGH.
DIO BI state: false
DIO CI state: false
DIO_DI state: false
```


<!-- Source: docs/code-samples/ain-examples.md -->


# Analog Input Examples

In this section, the Analog Inputs of the TDC-E device are discussed. Programming examples are given and thoroughly explained.

## Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. From it, the list of analog inputs and their current values can be read. The application also has an example of how to change the analog input mode.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/ain.proto
```

This installs needed tools, and generates files for the gRPC service.

First, dial options are set to dial a UNIX socket:

```go
conn, err := grpc.NewClient(
	"unix://"+socketPath,
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
	log.Fatalf("failed to connect to Unix socket: %v", err)
}
defer conn.Close()
```

The application then creates a gRPC client that connects to the service unix socket.

```go
conn, err := grpc.NewClient("passthrough:///unixsocket", dialOptions...)
if err != nil {
	log.Fatalf("failed to connect to Unix socket: %v", err)
}
defer conn.Close()

client := protos.NewAnalogINClient(conn)
ctx = metadata.NewOutgoingContext(ctx, metadata.New(nil))
```

An AnalogIn instance is created which is then used to access AnalogIn data. It is passed a list of AnalogIn devices to read, and the created client and context. AnalogIn devices are listed and they are enabled. The program enters an inifnite loop that reads AnalogIn values, until the user exits the application.

```go
ain, err := analogin.NewAnalogIn(AIDevices, client, ctx)
if err != nil {
	log.Fatalf("failed to create AnalogIN: %v", err)
}

ain.ListAINDevices()
ain.SetState(true)

done := make(chan struct{})
go func() {
	for {
		select {
		case <-ctx.Done():
			close(done)
			return
		default:
			ain.ReadAINVal()
			time.Sleep(1 * time.Second)
		}
	}
}()

 <-sigChan
fmt.Println("\nShutting down...")
ain.SetState(false)
cancel()
<-done
```

> **ℹ️ Info**
>
> To use the AnalogIn service, channels first need to be enabled! After usage, disable the channels!


To set the channel state, the following code is used:

```go
func (t *AnalogIN) SetState(state bool) {
	for i := range t.AIDevices {
		req := &protos.AnalogInSetStateRequest{
			Name:   t.AIDevices[i],
			Enable: state,
		}
		res, err := t.Client.SetState(t.Ctx, req)
		if err != nil {
			log.Printf("failed to set state to %v: %v", state, err)
			return
		}
		log.Printf("SetState response: %v\n", res.GetMessage())
	}
}
```

Reading AnalogInput values is done like this:

```go
func (t *AnalogIN) ReadAINVal() {
	for i := range t.AIDevices {
		req := &protos.AnalogInReadRequest{
			Name: t.AIDevices[i],
		}
		res, err := t.Client.Read(t.Ctx, req)
		if err != nil {
			log.Fatalf("could not read: %v", err)
		}
		log.Printf("%s Values:\nadcRaw: %d\nconverted: %f\nunit: %s\n", t.
AIDevices[i], res.AdcRaw, res.Converted, res.Unit)
	}
}
```

The response is printed to the terminal and contains the following values:

- 'adcRaw'
- 'Converted'
- 'Unit'

The program listens to a SIGINT and SIGTERM signal from the user, and will exit upon receiving it.

### Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the AnalogIn service is the following:

```bash
/**
 * AnalogIn service.
 *
 * Service that enables a control of the AnalogIn features of device
 */
syntax = "proto3";

package hal.analogin;

import "google/protobuf/empty.proto";

option go_package = "./protos;protos";

enum AnalogInMode {
  Voltage = 0;
  Current = 1;
}

enum AnalogInUnits {
  V = 0;
  mA = 1;
}

message AnalogInDevice {
  string name = 1;
  AnalogInMode mode = 2;
}

/**
 * Represents the ListDevice response data.
 */
message AnalogInListDeviceResponse {
  repeated AnalogInDevice devices = 1;
}

/**
 * Represents the Read request data.
 */
message AnalogInReadRequest {
  string name = 1;
}

/**
 * Represents the Read response data.
 */
message AnalogInReadResponse {
  uint32 adcRaw = 1;
  float converted = 2;
  AnalogInUnits unit = 3;
}

/**
 * Represents the SetMeasureMode request data.
 */
message AnalogInSetMeasureModeRequest {
  string name = 1;
  AnalogInMode mode = 2;
}

message AnalogInAttachResponse {
  string channel = 1;
  string status = 2;
  string timestamp = 3;
}

/**
 * Represents the request to set the AnalogIn service state for a specific chan
nel.
 */
message AnalogInSetStateRequest {
  string name = 1;  // name of the channel to set the state for
  bool enable = 2;  // true = ON, false = OFF
}

/**
 * Represents the response to the set AnalogIn service state.
 */
message AnalogInSetStateResponse {
  bool success = 1;
  string message = 2;
}

/**
 * Represents the request to get the AnalogIn service state for a specific chan
nel.
 */
message AnalogInGetStateRequest {
  string name = 1; // channel name like "channel0"
}

/**
 * Represents the response to the AnalogIn service state.
 */
message AnalogInGetStateResponse {
  bool enabled = 1;        // True if the channel is being actively polled
  string message = 2;      // status message
}

/**
 * Service exposing AnalogIn functions.
 */
service AnalogIN {
  /// Used to retrieve the current state of the analog input channel.
  rpc GetState (AnalogInGetStateRequest) returns (AnalogInGetStateResponse) {}

  /// Used to set the analog input channel (ON/OFF).
  rpc SetState(AnalogInSetStateRequest) returns (AnalogInSetStateResponse) {}

  /// Used to retrieve all available analog input devices.
  rpc ListDevices(google.protobuf.
Empty) returns (AnalogInListDeviceResponse) {}

  /// Used to read value of a particular analog input channel.
  rpc Read(AnalogInReadRequest) returns (AnalogInReadResponse) {}
  
  /// Used to set measure mode of a particular analog input channel.
  rpc SetMeasureMode(AnalogInSetMeasureModeRequest) returns (google.protobuf.
Empty) {}

  /// Used to monitor for overcurrent events
  rpc Attach(google.protobuf.Empty) returns (stream AnalogInAttachResponse) {}
}

```
</details>

To list all available Analog Input services, use 'grpcurl', which is an open-source utility for accessing gRPC services via the shell. For help setting up the 'grpcurl' command, refer to [gRPC Usage](/getting-started/grpc-usage).

To list all available Analog Input calls, use the following line:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -
plaintext <device_ip>:<grpc_server_port> list hal.analogin.AnalogIN
```

The 'token' field is the fetched TDC-E authorization token. For help fetching this token, refer to [gRPC Usage](/getting-started/grpc-usage). The 'device-ip:grpc_server_port' is the TDC-E IP address and the gRPC serving port. For example, if the 'token' value was 'token' and the address and port were '192.168.0.100:8081', you would use the following line to list all available Analog Input services.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 list hal.analogin.AnalogIN
```

The response should be in this format:

```log
hal.analogin.AnalogIN.Attach
hal.analogin.AnalogIN.ListDevices
hal.analogin.AnalogIN.Read
hal.analogin.AnalogIN.GetState
hal.analogin.AnalogIN.SetState
hal.analogin.AnalogIN.SetMeasureMode
```

Additionally, you can use the 'gRPC Clicker' VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).


### Application Deployment

#### Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```docker
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o ain-grpc ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ain-grpc .

CMD ["./ain-grpc"]
```

Open a terminal and paste the following commands:

```bash
docker build -t ain-grpc-app .
docker save -o ain-grpc-app.tar ain-grpc-app:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

#### Dockerfile Breakdown

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' and 'go.sum' files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```dockerfile
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
CGO_ENABLED=0
```

The application is built as 'ain-grpc'.

```dockerfile
RUN go build -o ain-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'ain-grpc' application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ain-grpc .

CMD ["./ain-grpc"]
```

#### Deploying to Portainer

To deploy the application to the TDC-E device, [Portainer](https://192.168.0.100:9443/) can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## Node-RED gRPC Example

```c
Download Example Code from embedded file: "grpc-ain.json"
```

The 'gRPC' node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For 'gRPC' node installation, import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

This section describes the creation and usage of a Node-RED gRPC example using analog inputs. The application makes listing all AI devices, reading AIN values, changing AIN mode and streaming changes in the AI devices possible. 

For implementation, the following nodes are used:

- 'inject' node
- 'gRPC call' node
- 'debug' node.

The 'gRPC' node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- 'Server':'unix:///var/run/hal/hal.sock'
- a 'Proto File' is provided 

To test out any of the listed functionalities, use the 'inject' node. The result should appear in the 'debug' node. 

See a screenshot of the application below:

![AIN Node-RED Example](static/img/AIN-node-red.png)

## Lua Example


See the example result of the application run below:


<!-- Source: docs/code-samples/temperature-sensor-examples.md -->


# Temperature Sensor Examples

In this section, temperature sensors on the TDC-E device are discussed. Programming examples are given and thoroughly explained.

## Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that lists all temperature sensors and prints their current value is created.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/temperature-sensor-service.
proto
```

This installs needed tools, and generates files for the gRPC service. The application also requires an access token to connect to the gRPC service the TDC-E provides. Please refer to [gRPC Usage](/getting-started/grpc-usage) for instructions on how to generate an access token for the TDC-E.

Once you've obtained the token, navigate to 'grpc-temperature/pkg/auth/token.json'. Paste your generated token in the 'access_token' field. This token is then used to create a 'context' which will be used to create gRPC calls.

The main application first creates a new gRPC client that connects to the service unix socket. It does so by using the generated gRPC Go file specification. 

```go
conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.
NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
```

To demonstrate the listed functionalities, the application lists all available temperature sensors on your TDC-E device, and then the temperature value of each specific sensor measured in '°C' is printed.

To list all temperature devices, an empty 'pb' request is sent to the server. The client specifies the function 'ListDevices' to receive all devices from the TDC-E device, which are then printed to the user.

```go
req := &emptypb.Empty{}
res, err := client.ListDevices(context.Background(), req)
```

Afterwards, each sensor's temperature is read. 

```go
func readTemperature(ctx context.Context, client protos.
TemperatureSensorClient, device string) {
	req := pb.TemperatureSensorReadRequest{
		Name: device,
	}
	res, err := client.Read(ctx, &req)
	if err != nil {
		log.Fatalf("couldn't read from temperature sensor: %v", err)
	}
	fmt.Printf("%s temperature: %.2f %s\n", device, res.Value, res.Unit)
}
```

This function creates a proto request for reading the temperature sensor from the given device name and prints the float result in °C.


### Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the TemperatureSensor service is the following:

```bash
/**
 * Temperature sensor service.
 *
 * Service that provides device temperature sensor reading
 */
syntax = "proto3";

package hal.temperaturesensor;

import "google/protobuf/empty.proto";

option go_package = "./protos;protos";


message TemperatureSensorDevice {
  string name = 1;
}

/**
 * Represents the ListDevice response data.
 */
message TemperatureSensorListDeviceResponse {
  repeated TemperatureSensorDevice devices = 1;
}

message TemperatureSensorReadRequest {
  string name = 1;
}

message TemperatureSensorReadResponse {
  float value = 1;
  string unit = 2;
}

/**
 * Service exposing TemperatureSensor functions.
 */
service TemperatureSensor {
  /// Used to retrieve all available temperature sensors
  rpc ListDevices(google.protobuf.
Empty) returns (TemperatureSensorListDeviceResponse) {}

  /// Used to read value of a particular temperature sensor.
  rpc Read(TemperatureSensorReadRequest) returns (TemperatureSensorReadResponse
) {}
}
```

</details>

To list all available temperature sensor services, use 'grpcurl', which is an open-source utility for accessing gRPC services via the shell. For help setting up the 'grpcurl' command, refer to [gRPC Usage](/getting-started/grpc-usage).

Use the following line for listing them: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -
plaintext <device_ip>:<grpc_server_port> list hal.temperaturesensor.
TemperatureSensor
```

The 'token' field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage). The 'device-ip:grpc_server_port' is the TDC-E IP address and the gRPC serving port. For example, if the 'token' value was 'token' and the address and port were '192.168.0.100:8081', you would use the following line to list all available temperature sensor devices.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 list hal.temperaturesensor.TemperatureSensor
```

The response should be in this format:

```log
hal.temperaturesensor.TemperatureSensor.ListDevices
hal.temperaturesensor.TemperatureSensor.Read
```

Additionally, you can use the 'gRPC Clicker' VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

### Application Deployment

#### Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```dockerfile
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o temp-grpc ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/temp-grpc .

CMD ["./temp-grpc"]
```

Open a terminal and paste the following commands:

```bash
docker build -t temp-grpc-app .
docker save -o temp-grpc-app.tar temp-grpc-app:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

#### Dockerfile Breakdown

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' and 'go.sum' files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```dockerfile
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
CGO_ENABLED=0
```

The application is built as 'temp-grpc'.

```dockerfile
RUN go build -o temp-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'temp-grpc' application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/temp-grpc .

CMD ["./temp-grpc"]
```

#### Deploying to Portainer

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## Node-RED gRPC Example

```c
Download Example Code from embedded file: "grpc-temp.json"
```

The 'gRPC' node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For 'gRPC' node installation, import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

This section describes usage of the temperature sensors on the TDC-E device. A Node-RED gRPC application was created to demonstrate listing all temperature sensor devices and reading temperature values from said devices.

For implementation, the following nodes are used:

- 'inject' node
- 'gRPC call' node
- 'debug' node.

The 'gRPC' node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- 'Server':'unix:///var/run/hal/hal.sock'
- a 'Proto File' is provided 

To test out any of the listed functionalities, use the 'inject' node. The result should appear in the 'debug' node. 

See a screenshot of the application below:

![AIN Node-RED Example](static/img/TEMP-node-red.png)

## Lua Example


The 'extractTemperature' function checks whether the passed variable is a monitor. If not, the function returns the 'nil' value. Otherwise, it uses the monitor's 'get()' function to fetch the temperature value of the device.

```lua
local function extractTemperature(temperatureMonitor)
  if not temperatureMonitor then
    return nil
  end
  return temperatureMonitor:get()
end
```

See the result of the application run below.


<!-- Source: docs/code-samples/can-examples.md -->


# Controller Area Network Examples

In this section, the Controller Area Network (CAN) interface of the TDC-E device and examples of its usage are discussed. Programming examples are given and thoroughly explained.

## Setting up CAN Device

In this section, setting up the CAN device is discussed. Using the dedicated CAN HAL service, setup of the following parameters is possible:
 
 - Transceiver Power
 - Termination
 - Interface Mapping (namespace)

In the following sections, setting up the transceiver power, termination, and CAN namespace is discussed. Additionally, reading CAN statistics is described. This is done using the CAN HAL Service 'hal.can.Can'. Examples using 'grpcurl' are given below. For more information about using gRPC services, refer to [gRPC Usage](/getting-started/grpc-usage).

### Setting Up CAN Transceiver Power

The CAN HAL service provides a means to set up transceiver power for the CAN interface. The HAL service 'hal.can.Can.SetTransceiverPower' is used. An example of setting the transceiver power of the connected CAN device on is given below.

```bash
grpcurl -d '{"interfaceName":"CAN1","powerOn":true}' -H 'Authorization:
 Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/SetTransceiverPower
{}
```

Checking changes to the CAN state is done by using the 'hal.can.Can.GetTransceiverPower' HAL service.

```bash
grpcurl -d '{"interfaceName":"CAN1"}' -H 'Authorization: Bearer token' -
plaintext 192.168.0.100:8081 hal.can.Can/GetTransceiverPower
{
  "powerOn": true
}
```

### Setting Up CAN Termination

The CAN HAL service provides a means to set up CAN termination. The HAL service 'hal.can.Can.SetTermination' is used. An example of setting the termination of the connected CAN device on is given below.

```bash
grpcurl -d '{"interfaceName":"CAN1","enableTermination":true}' -
H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/
SetTermination
{}
```

Checking changes to the termination is done by using the 'hal.can.Can.GetTermination' HAL service.

```bash
grpcurl -d '{"interfaceName":"CAN1"}' -H 'Authorization: Bearer token' -
plaintext 192.168.0.100:8081 hal.can.Can/GetTermination
{
  "terminationEnabled": true
}
```

### Setting Up CAN Interface Mapping

When working with the CAN interface, one can use different namespaces. The CAN interface is then bound to a namespace and can be accessed from there. Examples include:
- AppEngine
- Host
- Any other running container

To set the namespace of the CAN interface, the HAL service 'hal.can.Can.SetInterfaceToContainer' is used. Possible values for the 'dockerContainerName' parameter include:
- 'app-engine' - maps the CAN interface to AppEngine
- empty string - maps the CAN interface to the host
- 'container-name' - maps the CAN interface to a running Docker container

See an example of exposing the CAN interface to the host below.

```bash
grpcurl -d '{"interfaceName":"CAN1","dockerContainerName":""}' -
H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/
SetInterfaceToContainer
{}
```

To see the changes to the interface mapping, the 'hal.can.Can.GetInterfaceTocontainerMapping' service is used.

```bash
grpcurl -d '{}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:
8081 hal.can.Can/GetInterfaceToContainerMapping
{
  "items": [
    {
      "interfaceName": "CAN1",
      "dockerContainerName": ""
    }
  ]
}
```

> **📝 Note**
>
> Mapping the CAN interface results in exclusive access to the interface.


Mapping the CAN interface to AppEngine will render access to the interface exclusive to AppEngine. Example usage can be see in the [Lua Example](#4-lua-example) below. Mapping the interface to the host, the CAN interface is visible only to the host. Example usage can be seen in the [Go Example](#2-go-example) below.

### Viewing CAN Statistics

To see statistics of your CAN device, the HAL service 'hal.can.Can.GetStatistics' is used. See an example of its usage below.

```bash
grpcurl -d '{"interfaceName":"CAN1"}' -H 'Authorization: Bearer token' -
plaintext 192.168.0.100:8081 hal.can.Can/GetStatistics
{
  "RxPackets": "1118",
  "TxPackets": "92",
  "RxBytes": "4783",
  "TxBytes": "736",
  "RxErrors": "0",
  "TxErrors": "0",
  "RxDropped": "0",
  "TxDropped": "1",
  "Multicast": "0",
  "Collisions": "0",
  "RxLengthErrors": "0",
  "RxOverErrors": "0",
  "RxCrcErrors": "0",
  "RxFrameErrors": "0",
  "RxFifoErrors": "0",
  "RxMissedErrors": "0",
  "TxAbortedErrors": "1",
  "TxCarrierErrors": "0",
  "TxFifoErrors": "0",
  "TxHeartbeatErrors": "0",
  "TxWindowErrors": "0",
  "RxCompressed": "0",
  "TxCompressed": "0"
}
```

### Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the CAN service is the following:

```bash
syntax = "proto3";

package hal.can;

import "google/protobuf/empty.proto";

option go_package = "./grpc;grpc";

enum BusStates {
  ERROR_ACTIVE = 0;
  ERROR_WARNING = 1;
  ERROR_PASSIVE = 2;
  BUS_OFF = 3;
  STOPPED = 4;
  SLEEPING = 5;
}

// Request to set transceiver power
message SetTransceiverPowerRequest { 
  string interfaceName = 1; // e.g., "CAN1" 
  bool powerOn = 2; // true to power on, false to power off
}

// Request to enable or disable termination
message SetTerminationRequest {
  string interfaceName = 1; // e.g., "CAN1"
  bool enableTermination = 2; // true to enable, false to disable
}

// Request to list available CAN interfaces
message ListInterfacesRequest {}

// Response containing the list of CAN interfaces
message ListInterfacesResponse {
  repeated string interfaces = 1; // List of available interface names (e.g.,
 ["CAN1", "CAN2"])
}

// Genereric request containing the name of the interface
message GetInterfaceNameRequest {
  string interfaceName = 1;  // e.g., "CAN1"
}

// Response containing power status of a CAN transceiver
message GetTransceiverPowerResponse {
  bool powerOn = 1; // true if enabled, otherwise false
}

// Response containing status of a CAN termination pin
message GetTerminationResponse {
  bool terminationEnabled = 1; // true if enabled, otherwise false
}

// Response containing the statistics of a CAN interface
message GetStatisticsResponse {
	uint64 RxPackets = 1;
	uint64 TxPackets = 2;
	uint64 RxBytes = 3;
	uint64 TxBytes = 4;
	uint64 RxErrors = 5;
	uint64 TxErrors = 6;
	uint64 RxDropped = 7;
	uint64 TxDropped = 8;
	uint64 Multicast = 9;
	uint64 Collisions = 10;
	uint64 RxLengthErrors = 11;
	uint64 RxOverErrors = 12;
	uint64 RxCrcErrors = 13;
	uint64 RxFrameErrors = 14;
	uint64 RxFifoErrors = 15;
	uint64 RxMissedErrors = 16;
	uint64 TxAbortedErrors = 17;
	uint64 TxCarrierErrors = 18;
	uint64 TxFifoErrors = 19;
	uint64 TxHeartbeatErrors = 20;
	uint64 TxWindowErrors = 21;
	uint64 RxCompressed = 22;
	uint64 TxCompressed = 23;
}

// Response containing the bitrate of a CAN interface
message GetBitrateResponse {
  uint32 bitrate = 1; // Number of transmitted frames
}

// Response containing the bitrate of a CAN interface
message GetBusStateResponse {
  BusStates state = 1; // Number of transmitted frames
}

// Request containing parameters for SetInterfaceToContainer method
message SetInterfaceToContainerRequest {
  string interfaceName = 1;  // e.g., "CAN1"
  string dockerContainerName = 2;  // e.g., "app-engine",
 leave empty to move to host
}

// Represents a mapping of interface to docker container
message InterfaceToContainerMapping {
  string interfaceName = 1;  // e.g., "CAN1"
  string dockerContainerName = 2;  // e.g., "app-engine",
 if empty string then interface is mapped to host
}

// Response contains mapping of interfaces to docker containers
message GetInterfaceToContainerMappingResponse {
  repeated InterfaceToContainerMapping items = 1; //
 List of interface to docker container mappings 
}

/**
 * Service exposing CAN functions.
 */
service Can {
  // Gets the state of transceiver power
  rpc GetTransceiverPower(GetInterfaceNameRequest) returns (GetTransceiverPower
Response);

  // Gets the state of termination
  rpc GetTermination(GetInterfaceNameRequest) returns (GetTerminationResponse);

  // Sets the transceiver power (on/off)
  rpc SetTransceiverPower(SetTransceiverPowerRequest) returns (google.protobuf.
Empty) {}

  // Enables or disables CAN termination
  rpc SetTermination(SetTerminationRequest) returns (google.protobuf.Empty) {}

  // Lists available CAN interfaces
  rpc ListInterfaces(ListInterfacesRequest) returns (ListInterfacesResponse) {}

  // Gets the CAN interface statistics
  rpc GetStatistics(GetInterfaceNameRequest) returns (GetStatisticsResponse) {}

  // Gets the CAN interface bitrate
  rpc GetBitrate(GetInterfaceNameRequest) returns (GetBitrateResponse) {}

  // Gets the CAN interface state
  rpc GetBusState(GetInterfaceNameRequest) returns (GetBusStateResponse) {}

  // Moves specified interface into namespace of a specified docker container
  rpc SetInterfaceToContainer(SetInterfaceToContainerRequest) returns (google.
protobuf.Empty) {}

  // Returns mapping of interfaces onto docker containers
  rpc GetInterfaceToContainerMapping(google.protobuf.
Empty) returns (GetInterfaceToContainerMappingResponse) {}
}
```

</details>

## Go Example


A Go application is provided as a CAN usage example. The application creates a CAN bus, binds it to a CAN port, and sends and receives CAN data simultaneously. For application testing, the 'Kvaser Leaf Light v2' device was used. For generating and testing data, Kvaser's 'CanKing' application and drivers were used.

### Application Implementation

The application is implemented using the [\"canbus\" Go package](https://pkg.go.dev/github.com/go-daq/canbus@v0.2.0) to communicate with the CAN bus. 

> **📝 Note**
>
> Check your CAN usage namespace, and check if the CAN link is up. Check the [previous section](#13-setting-up-can-interface-mapping) to see CAN namespace assignment. If the CAN link is already enabled, skip the next step, and run the application normally.


If the CAN link is down, whilst running the application, add a '--setup=true' flag to the starting arguments, which will set up the CAN interface according to specifications in the '/pkg/canSetup' file. The setup sets the CAN bitrate, powers on the transceiver and brings the interface link up.

```bash
./can --setup=true
```

Otherwise, run the application normally. Firstly, a new CAN bus connection is established using the following lines:

```go
can, err := canbus.New()
if err != nil {
	log.Fatalf("Failed to open CAN socket: %v", err)
}
defer can.Close()
```

The bus is then bound to the 'can0' interface.

```go
err = can.Bind("can0")
if err != nil {
	log.Fatalf("Cannot bind to can socket: %s", err)
}
```

Two goroutines are started simultaneously. One sends sample data from the CAN device, while the other listens to received data and prints the data to the console. The goroutines are synced using a wait group that waits for all routines to finish execution. See both goroutines explained below.

The goroutine code snippet for sending data is shown below.

```go
data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
frame := canbus.Frame{
	ID:   0x123,
	Data: data,
}

for {
	if _, err := can.Send(frame); err != nil {
		log.Printf("Failed to send CAN frame: %v", err)
	} else {
		fmt.Println("Sent CAN frame:", frame)
	}

	time.Sleep(1 * time.Second)
}
```

The specified code creates a data structure which contains the ID and data to be sent by a message, and attempts to send the CAN frame to the CAN device each second.

See how the goroutine receives CAN data below. A frame variable is used which listens to the CAN device, and the device logs all received data as it's sent.

```go
frame, err := can.Recv()
if err != nil {
  log.Printf("Failed to receive CAN frame: %v", err)
} else {
	log.Println("Received CAN frame:", frame)
}
```

See a screenshot of sending and receiving CAN data below.

![Sending Data from CAN](static/img/cansendGo.PNG)

An example print from the application can be seen below.

```log
Received CAN frame: {1122 [220 3 0 0 0] SFF}
Received CAN frame: {586 [221 3 0] SFF}
Received CAN frame: {1792 [222 3 0 0 0 0 0 0] SFF}
Received CAN frame: {538 [223 3 0 0] SFF}
Received CAN frame: {511 [225] SFF}
Received CAN frame: {931 [226 3 0 0] SFF}
Received CAN frame: {1904 [227 3 0 0 0 0] SFF}
Received CAN frame: {653 [] SFF}
Received CAN frame: {1520 [229] SFF}
Received CAN frame: {1565 [230 3 0 0 0 0 0] SFF}
Received CAN frame: {854 [231 3 0 0 0] SFF}
Sent CAN frame: {291 [1 2 3 4 5 6 7 8] SFF}
Sent CAN frame: {291 [1 2 3 4 5 6 7 8] SFF}
Sent CAN frame: {291 [1 2 3 4 5 6 7 8] SFF}
```

### Application Deployment

#### Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```dockerfile
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o can ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/can .

CMD ["./can"]
```

Open a terminal and paste the following commands:

```bash
docker build -t can-app .
docker save -o can-app.tar can-app:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

#### Dockerfile Breakdown

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' and 'go.sum' files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```dockerfile
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
CGO_ENABLED=0
```

The application is built as 'can'.

```dockerfile
RUN go build -o can ./cmd/main.go
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'can' application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/can .

CMD ["./can"]
```

#### Deploying to Portainer

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

## NodeRED Example

```c
Download Example Code from embedded file: "can.json"
```

The following NodeRED example will utilize the NodeRED palette [\"node-red-contrib-socketcan\"](https://flows.nodered.org/node/node-red-contrib-socketcan) to communicate with the CAN bus. 

To start with the example, download NodeRED from the Applications tab. For more information, refer to [Installing, Accessing and Removing Applications](../getting-started/installing-accessing-removing-applications).


### NodeRED Setup

Next open NodeRED's web user interface and in the top-right corner click on the menu icon next to the "Deploy" button, and select "Manage palette". Next click on "Install" card and search for  "node-red-contrib-socketcan" and install the palette.

![NodeRED manage palette](static/img/can_nodered_socketcan_install.png)

Import the CAN example code by pressing 'Ctrl + i' or by right clicking a flow → 'Insert' → 'Import'.

### Flow Usage & Breakdown

After successfully importing the example flows, you should have a workspace similar to that of image below:

![NodeRED example flows A](static/img/can_nodered_flows_a.png)

The example is composed of two flows, a '#1 Step - CAN Setup' and a '#2 Step - CAN Sender' flow. We  Use flow #1 to setup the CAN device on a container level, and we use flow #2 to demonstrate the usage of said CAN device.

'#1 Step - CAN Setup' flow is composed of two parts, a automatic flowchart (top section) and a manual debugging flow. Use the automatic flow to setup everything for CAN by clicking on the blue square on the 'Setup CAN!' node. All of the red nodes in this flow run one of the following shell commands:

- 'ip link' - lists all network devices
- 'ip link set can0 [up/down]' - enables/disables the can0 device 
- 'ip link set can0 type can bitrate 500000' - sets the bitrate to 500kbps, all other devices need the same bitrate to communicate
- 'apk update && apk add iproute2' - get standard implementation of iproute2

> **📝 Note**
>
> By default, the 'nodered' container already comes with the version of 'ìproute2', but it's minimal in features and doesn't have the ability to manage CAN devices out of the box. For this reason we are installing the standard implementation using 'apk'.


![NodeRED example flows B](static/img/can_nodered_flows_b.png)

'#2 Step - CAN Sender' is composed of 3 sections, a looping CAN sending flowchart (contains 'Loop Trigger' node), two manual CAN sending flowcharts (inject nodes that start with _Send..._) and a CAN dump flowchart (contains 'Pretty Print CAN Frame' node). You can freely use the inject nodes to try out the CAN functionality.
Look inside the manual sections inject nodes on possible ways to send payload to 'socketcan-in' nodes or consult the palettes documentation section [Sending CAN frames](https://flows.nodered.org/node/node-red-contrib-socketcan).

> **📝 Note**
>
> The green indicator under the 'socketcan-*' can be a bit misleading, since it can show as connected even if the CAN device is down. Verify if the device is active in flow #1 by triggering the 'ip link' node.


After sending packets with any of the two CAN sending flowcharts, NodeREDs debug ouput and 'candump' CLI tools output (from the 'can-utils' package, as a alternative for 'CanKing' application) should look like this:

![NodeRED CAN Outputs](static/img/can_nodered_output.png)

### Socketcan Node Configuration

If the device contains different CAN devices from can0, you can configure the socketcan palette nodes by editing both lime nodes with names 'socketcan-*' and under the 'Interface' field add can0, can1 or other CAN bus devices. For example if you wish to use the virtual device vcan0, you can configure the node like this:

![NodeRED socketcan configure](static/img/can_nodered_socketcan_configure.png)

## Lua Example


A Lua application is provided as a CAN usage example. The application creates a CAN handler, opens the handler and sends or receives CAN data. For application testing, the 'Kvaser Leaf Light v2' device was used. For generating and testing data, Kvaser's 'CanKing' application and drivers were used.

The application first initializes the CAN socket. The 'CAN1' handle is created, and the baud rate is set to '500000', and the handle is opened.

```lua
Handle = CANSocket.create('CAN1')
Handle:setBaudRate(500000)
Handle:open()
```

For printing received data and logging errors, two local functions were implemented.

```lua
--@handleOnReceive(id:int,buffer:binary)
local function handleOnReceive(id, buffer)
  print("Received data.")
  print(id)
end

--@handleOnReceive(errorCode:int)
local function handleOnError(errorCode)
  Log.info("Error code " .. errorCode)
end

CANSocket.register(Handle, "OnReceive", handleOnReceive)
CANSocket.register(Handle, "OnError", handleOnError)
```

Reading is tested with the 'CanKing' application. See a screenshot of sending CAN data below.

![Sending Data from CAN](static/img/cansend.PNG)

To send data to the device, three test IDs and messages are created. The script then sends the defined data in an infinite loop.

```lua
Ids = {20, 21, 22}
Msgs = {"\x41\x42\x43\x44\x45\x46\x47\x48", "\x51\x52\x53\x54\x55\x56\x57\x58",
 "\x61\x62\x63\x64\x65\x66\x67\x68"}
while true do
  Handle:transmit(Ids, Msgs)
  Script.sleep(1000)
end
```

See the result of the application run below.
```log
[12:52:10.131: INFO: AppEngine] Starting app: 'can' (priority: LOW)
Received data.
776
Received data.
779
Received data.
693
Received data.
857
```


<!-- Source: docs/code-samples/gnss-examples.md -->


# GNSS Examples


In this section, GNSS on the TDC-E device is discussed. Programming examples are given and thoroughly explained.

## Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that checks device availability and fetches GNSS data is provided.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/gnss.proto
```

This installs needed tools, and generates files for the gRPC service. The application also requires an access token to connect to the gRPC service the TDC-E provides. Please refer to [gRPC Usage](/getting-started/grpc-usage) for instructions on how to generate an access token for the TDC-E.

Once you've obtained the token, navigate to 'grpc-gnss/pkg/auth/token.json'. Paste your generated token in the 'access_token' field. This token is then used to create a 'context' which will be used to create gRPC calls.

The application then creates a gRPC client that connects to the service unix socket. It does so by using the generated gRPC Go file specification.

```go
conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.
NewCredentials()))
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

### Proto File

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
    double latitude = 1;     //
 The latitude of the GNSS position in decimal degrees.
    double longitude = 2;    //
 The longitude of the GNSS position in decimal degrees.
    double altitude = 3;     // The altitude of the GNSS position in meters.
    string timestamp = 4;    //
 The timestamp of the GNSS record in ISO 8601 format.
    double speed_kph = 5;    // The speed of the object in kilometers per hour.
    double speed_mph = 6;    // The speed of the object in miles per hour.
    double course = 7;       // The direction of travel in degrees,
 where 0 is north.
    double satellites = 8;   //
 The number of satellites used to determine the position.
    double fix_type = 9;     // Indicates the type of GNSS fix (e.g.,
 1 for no fix, 2 for 2D fix, 3 for 3D fix).
    double fix_quality = 10;  // Indicates the GNSS fix quality (e.g.,
 0 for invalid, 1 for GPS fix, 2 for DGPS fix).
    double hdop = 11;        // The Horizontal Dilution of Precision,
 a measure of the GNSS signal's accuracy.
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
    rpc GetDeviceStatus(google.protobuf.
Empty) returns (DeviceManagementResponse) {}
  
    // Provides a way do disable/enable receiver and antenna
    rpc SetDeviceStatus(DeviceManagementRequest) returns (google.protobuf.
Empty) {}
  
    // Provides a way to set preferred constellations
    rpc SetConstellations(GnssConstellationsSettings) returns (google.protobuf.
Empty) {}
  
    // Provides a way to get current preferred constellations
    rpc GetConstellations(google.protobuf.
Empty) returns (GnssConstellationsSettings) {}
  
    // Provides a way to get a list of available constellations
    rpc ListAvailableConstellations(google.protobuf.
Empty) returns (GnssConstellations) {}
  
    // Provides a way to set NMEA output frequency
    rpc SetNmeaOutputFrequency(NMEAFrequencySettings) returns (google.protobuf.
Empty) {}
  
    // Provides a way to set NMEA output frequency
    rpc GetNmeaOutputFrequency(google.protobuf.
Empty) returns (NMEAFrequencySettings) {}
  
    // Provides a way to get a list of available NMEA output frequencies
    rpc ListNmeaOutputFrequencies(google.protobuf.
Empty) returns (NMEAFrequenciesList) {}
  }

```

</details>

To get the GNSS device status, use 'grpcurl', which is an open-source utility for accessing gRPC services via the shell. For help setting up the 'grpcurl' command, refer to [gRPC Usage](/getting-started/grpc-usage).

Use the following line: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -
plaintext <device_ip>:<grpc_server_port> hal.gnss.Gnss.GetDeviceStatus
```

The 'token' field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage). The 'device-ip:grpc_server_port' is the TDC-E IP address and the gRPC serving port. For example, if the 'token' value was 'token' and the address and port were '192.168.0.100:8081', you would use the following line to see the device status of the GNSS service.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -
plaintext 192.168.0.100:8081 hal.gnss.Gnss.GetDeviceStatus
```

The response should be in this format:

```log
"enableReciever": true,
"enableAntenna": true,
"sessionActive": true,
"rebootRequired": false
```

Additionally, you can use the 'gRPC Clicker' VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

### Application Deployment

#### Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a 'Dockerfile' is created. The file is shown below.

```docker
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o gnss-grpc ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gnss-grpc .

CMD ["./gnss-grpc"]
```

Open a terminal and paste the following commands:

```bash
docker build -t gnss-grpc-app .
docker save -o gnss-grpc-app.tar gnss-grpc-app:latest
```

This will build the docker container and save the application as a '.tar' file which can be used for Portainer upload.

#### Dockerfile Breakdown

The 'Dockerfile' first creates a build image that is used to build the Go application. It sets the working directory as '/app', then copies the 'go.mod' and 'go.sum' files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the 'arm32hf' architecture and is based on 'Linux'. With this in mind, the application is set to the following:

```dockerfile
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
CGO_ENABLED=0
```

The application is built as 'gnss-grpc'.

```dockerfile
RUN go build -o gnss-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest 'alpine' version for a smaller image size. The working directory is once again set to '/app', and the application is copied. The last line specifies that, upon deployment, the 'gnss-grpc' application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gnss-grpc .

CMD ["./gnss-grpc"]
```

#### Deploying to Portainer

To deploy the application to the TDC-E device, 'Portainer' can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## Node-RED gRPC Example

```c
Download Example Code from embedded file: "grpc-gnss.json"
```

The 'gRPC' node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For 'gRPC' node installation, import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

This section describes usage of the the GNSS device on the TDC-E. A Node-RED gRPC application was created to demonstrate reading from the GNSS device.

For implementation, the following nodes are used:

- 'inject' node
- 'gRPC call' node
- 'debug' node.

The 'gRPC' node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- 'Server':'unix:///var/run/hal/hal.sock'
- a 'Proto File' is provided 

To test out any of the listed functionalities, use the 'inject' node. The result should appear in the 'debug' node. 

See a screenshot of the application below:

![GNSS Node-RED Example](static/img/node-red-gnss.png)

## Lua Example

```c
Download Example Code from embedded file: "gnss.lua"
```

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


<!-- Source: docs/code-samples/networking.md -->

---
title: Networking
sidebar_position: 208
---


# Networking

Networking REST API samples in Lua are prepared and shown below.

> **📝 Note**
>
> While the provided sample is written in **Lua**, the functionality is supported across multiple programming languages, including **Go and Node-RED**.


## Lua Examples

In this section Lua code examples will be shown.

### Control Center Admin Token

To get a required administrator token, see the snippet bellow how to request it through the TDC login interface.

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/auth/login"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("POST")
-- Prepare credentials.
request:addHeader("Accept", "application/json")
request:setContentType("application/json")
request:setContentBuffer('{"username":"admin","password":"myadminpassword",
"realm":"admin"}')
-- Execute request.
local response = client:execute(request)
local response_code = response:getStatusCode()
local response_body = response:getContent()
-- Extract token
local json = JSON.create()
local json_body = json:parseFromString(response_body) 
-- Token   
local login_token = json_body:getValue("/token")   
```
A JWT token is returned that needs to be included in every API request.

 - Example:
 ```jwt
eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiMTcyLjE4LjAuMiIsImVtYWlsIjoiYWRtaW5Ac2ljay5jb20iLCJleHAiOjE3MzY3NzQ5OTQsImlhdCI6MTczNjc3MTM5NCwiaXNzIjoiaHR0cDovLzE5Mi4xNjguMC4xMDAvYXV0aC9sb2dpbiIsImp0aSI6IjFab0dsR0xTOFJvQ3kwVEJkUHhhNFZ5V1RkeTlIWVpUQUJ3dnBEazciLCJuYmYiOjE3MzY3NzXjnxjdnsajkfbajiI6ImFkbWluIiwicmVhbG0iOiJhZG1pbiIsInJvbGVzIjpbImF1dGhwL3VzZXIiLCJjY2F1dGgvYWRtaW4iXSwic3ViIjoiYWRtaW4ifQ.pNikIbs9JRzkaRI3gU9GXqyXYn2Z-de9Na22bi8NJe7sGo-pzUTV61RipNWEnYOxn66x-snzie_Wy_17fdxOtg
 ```

```c
Download Example Code from embedded file: "networking_cc_token.lua"
```

> **ℹ️ Info**
>
> Please download required helper script: ```c
> Util from embedded file: "networking_util.lua"
> ```les/networking_util.lua" download>Util</a>.


### 1.2. Interface List

Get a list of all network interfaces.

The list is returned as a JSON array containing the name and all interface prop
erties.

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/interfaces/"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")
-- Execute.
local response = client:execute(request)
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

 - Response

```json
[
  {
    "isConfigurable": true,
    "name": "eth1",
    "settings": {
      "defaultGateway": "192.168.0.1",
      "dhcpFallbackAddress": "",
      "dhcpFallbackGateway": "",
      "dhcpFallbackMode": 0,
      "enabled": true,
      "mac": "f8:dc:7a:b0:05:1e",
      "staticAddress": "192.168.0.100/24",
      "useDhcp": false
    },
    "status": {
      "address": "192.168.0.100/24",
      "defaultGateway": "192.168.0.1",
      "enabled": true,
      "state": 6
    },
    "type": 1
  }
]
```

```c
Download Example Code from embedded file: "networking_interface_list.lua"
```

### 1.3. Interface Settings

#### 1.3.1. Get Network Interface Settings

Get the settings of an interface by name.

```url
http://192.168.0.100/api/v1/network/interfaces/{name}/settings
```

```lua
-- Define endpoint
local endpoint = "http://192.168.0.100/api/v1/network/interfaces/eth1/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)  
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

- Response

```json
{
  "defaultGateway": "192.168.0.1",
  "dhcpFallbackAddress": "",
  "dhcpFallbackGateway": "",
  "dhcpFallbackMode": 0,
  "enabled": true,
  "mac": "f8:dc:7a:b0:05:1e",
  "staticAddress": "192.168.0.100/24",
  "useDhcp": false
}
```

```c
Download Example Code from embedded file: "networking_interface_settings.lua"
```

### 1.4. Set Network Interface Mode

Set DHCP or static mode for a network interface.

When settings are accepted and applied, status code `204` is returned.

> **📝 Note**
>
> Value of all true or false type variables must be boolean.


#### 1.4.1. Set DHCP Mode

To set a DHCP mode for an interface,
 JSON with settings must be sent using the PUT HTTP method.

> **📝 Note**
>
> When setting DHCP mode, a static fallback address must be set.


- JSON data

```json
{
  "enabled":true, 
  "useDhcp":true, 
  "staticAddress":"", 
  "defaultGateway": "", 
  "dhcpFallbackAddress":"192.168.2.1/24", 
  "dhcpFallbackGateway":"0.0.0.0", 
  "dhcpFallbackMode":0
}
```

- LUA example code

```lua
-- Define endpoint
local endpoint = "http://192.168.0.100/api/v1/network/interfaces/eth2/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("PUT")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Prepare setting data.
request:setContentType("application/json")
request:setContentBuffer('{"enabled":true, "useDhcp":true, "staticAddress":"", "defaultGateway": "", "dhcpFallbackAddress":"192.168.2.1/24", "dhcpFallbackGateway":"0.0.0.0", "dhcpFallbackMode":0}')

-- Execute.
local response = client:execute(request)  
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

#### 1.4.2. Set Static IP Address

- JSON data

```json
{
  "enabled":true,
  "useDhcp":false,
  "staticAddress":"192.168.2.100/24",
  "defaultGateway": "192.168.2.1",
  "dhcpFallbackAddress":"",
  "dhcpFallbackGateway":"",
  "dhcpFallbackMode":0
}
```
- LUA example code

```lua
-- Define endpoint
local endpoint = "http://192.168.0.100/api/v1/network/interfaces/eth2/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("PUT")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Prepare setting data.
request:setContentType("application/json")
request:setContentBuffer('{"enabled":true, "useDhcp":false, "staticAddress":"192.168.2.100/24", "defaultGateway": "192.168.2.1", "dhcpFallbackAddress":"", "dhcpFallbackGateway":"", "dhcpFallbackMode":0}')

-- Execute.
local response = client:execute(request)  
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

```c
Download Example Code from embedded file: "networking_interface_mode.lua"
```

#### 1.4.3. Set an Fallback IP Address

static IP, see sections: 

 - [Set DHCP mode](/code-samples/networking#141-set-dhcp-mode)
 - [Set Static IP address](/code-samples/networking#142-set-static-ip-address)


### 1.5. Setup WLAN Interface

Check sections below for further WLAN setup details.

```c
Download Example Code from embedded file: "networking_wlan_interface.lua"
```

#### 1.5.1. Get WLAN Interface Name

Get TDC Wireless interface name.

 - LUA example code

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/wlan"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```
- Response

```json
[
  "wlan0"
]
```

#### 1.5.2. WLAN Interface Details

To get details of a specific WLAN interface, the WLAN name must be known.

 - LUA example code

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/wlan/wlan0"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```
 - Response

```json
{
  "name": "wlan0",
  "settings": {
    "apMode": {
      "channel": 1,
      "encryption": "",
      "hidden": false,
      "restricted": false,
      "ssid": "TDCE-Next"
    },
    "countryCode": "HR",
    "enabled": true,
    "mode": "station",
    "stationMode": {
      "scanningEnabled": true
    }
  },
  "status": {
    "connectedSsid": "my-network1-name",
    "signalLevel": 96,
    "state": "connected"
  }
}
```

#### 1.5.3. WLAN Network List

Get a list of all available networks (saved or discovered) for one WLAN interfa
ce.

The list of available networks can be used when adding a new network.

```lua
local endpoint = "http://192.168.0.100/api/v1/network/wlan/wlan0/networks"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

 - Response

```json
[
  {
    "actualConnectionState": "connected",
    "bssid": "0a:17:68:88:5d:0d",
    "channel": 0,
    "events": [],
    "flags": [],
    "isInRange": true,
    "quality": 96,
    "reconnect": false,
    "ssid": "my-network1-name",
    "stored": true,
    "strength": -52,
    "username": ""
  },
  {
    "actualConnectionState": "disconnected",
    "bssid": "d4:01:ff:01:f5:aa",
    "channel": 0,
    "events": [],
    "flags": [],
    "isInRange": true,
    "quality": 64,
    "reconnect": false,
    "ssid": "other-network1",
    "stored": false,
    "strength": -68,
    "username": ""
  },
  ...
```

#### 1.5.4. Add Wireless Network

Connection to a new network requires a JSON body to be sent.

This same example is used when modifying an existing network connection.

```json
{
  "connectionState": "connected",
  "passphrase": "wifipassword",
  "reconnect": true,
  "username": ""
}
```
Connection state of a new network can be predefined.

 - `connectionState` should be `connected`.

```
 - unknown
 - connected
 - connecting
 - disconnected
 - disconnecting
```

 - LUA example code

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/wlan/wlan0/networks/my-network1-name"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("PUT")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Prepare setting data.
request:setContentType("application/json")
request:setContentBuffer('{"connectionState": "connected","passphrase": "networkpassword","reconnect": true,"username": ""}')

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

#### 1.5.5. Delete Wireless Network

Removing previously defined wireless network.

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/wlan/wlan0/networks/my-network1-name"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("DELETE")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```
### 1.6. Setup Gateway Priority

Get or set gateway priority list.

```c
Download Example Code from embedded file: "networking_gateway.lua"
```

#### 1.6.1. Get Gateway Priority List

Request current gateway priority list.

```lua
-- Define endpoint
local endpoint = "http://192.168.0.100/api/v1/network/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)  
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

 - Response

```json
{
  "defaultGatewayPriority": [
    "wwan0",
    "eth1",
    "eth2",
    "wlan0"
  ],
  "networkDriver": "ControlCenter"
}
```

#### 1.6.2. Set Gateway Priority List

Priority list must be sent in the same form as current json response.

- Content JSON

```json
{
  "defaultGatewayPriority": [
    "eth1",
    "eth2",
    "wlan0",
    "wwan0"
  ],
  "networkDriver": "ControlCenter"
}
```
- LUA example code

```lua
-- Define endpoint
local endpoint = "http://192.168.0.100/api/v1/network/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("PUT")

-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Prepare setting data.
request:setContentType("application/json")
request:setContentBuffer('{"defaultGatewayPriority": ["eth2","eth1","wlan0","wwan0"], "networkDriver": "ControlCenter"}')

-- Execute.
local response = client:execute(request)  
-- Get status code & json response.
local response_code = response:getStatusCode()
local response_body = response:getContent()
```
- See: [Get gateway priority list](/code-samples/networking#161-get-gateway-
priority-list)


### 1.7. Modem

Modem status and setup examples.

```c
Download Example Code from embedded file: "networking_modem.lua"
```

#### 1.7.1. Get Modem Interface

Get list of all modem devices on a system.

- LUA example code

```lua
-- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/modem"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

- Response

```json
[
  "wwan0"
]
```

#### 1.7.2. Modem Status

Get all modem settings and statuses.

> **📝 Note**
>
> Settings section will not contain password data.


 - Lua example

```lua
  -- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/modem/wwan0"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

 - Response

```json
{
  "name": "wwan0",
  "settings": {
    "apn": "internet.ht.hr",
    "apnPassword": "",
    "apnUsername": "",
    "connectionEnabled": true,
    "dialString": "",
    "enabled": true,
    "persistConnection": true
  },
  "status": {
    "ccid": "8938599202301824969",
    "currentAccessTechnologies": [
      "lte"
    ],
    "currentCapabilities": [
      "gsmUmts",
      "lte"
    ],
    "dataConnectionInfo": {
      "attempts": 1,
      "connected": true,
      "connectionError": "",
      "dns": [
        "195.29.247.161",
        "195.29.247.162"
      ],
      "downlinkSpeed": 0,
      "duration": 8940,
      "failedAttempts": 0,
      "gateway": "10.153.202.125",
      "ipCidr": "10.153.202.126/30",
      "rxBytes": 1383,
      "startDate": "2025-01-14T07:13:05Z",
      "totalDuration": 8940,
      "totalRxBytes": 1383,
      "totalTxBytes": 678,
      "txBytes": 678,
      "uplinkSpeed": 0
    },
    "imei": "869283050806985",
    "imsi": "219019912652496",
    "lockSimStatus": "simPin2",
    "modemState": "connected",
    "operatorIdentifier": "21901",
    "operatorName": "HT HR",
    "powerState": "on",
    "signalStrength": 100,
    "simActive": true,
    "stateFailedReason": "none"
  }
}
```

#### 1.7.3. Read Modem Setting

Get modem settings.

> **📝 Note**
>
> Settings section will not contain password data.


 - Lua example

```lua
  -- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/modem/wwan0/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("GET")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

 - Response

```json
{
  "apn": "internet.ht.hr",
  "apnPassword": "",
  "apnUsername": "",
  "connectionEnabled": true,
  "dialString": "",
  "enabled": true,
  "persistConnection": true
}
```

#### 1.7.4. Add Modem Setting

Add modem settings.

To add modem connection setting JSON data is required.

```json
{
  "apn": "string",
  "enabled": true,
  "connectionEnabled": true,
  "apnPassword": "",
  "apnUsername": "",  
  "dialString": "",  
  "persistConnection": true
  }
```
 - Lua example

```lua
  -- Define endpoint.
local endpoint = "http://192.168.0.100/api/v1/network/modem/wwan0/settings"
-- Create client.
local client = HTTPClient.create()
-- Create request.
local request = HTTPClient.Request.create()
request:setURL(endpoint)
request:setMethod("PUT")
-- Prepare credentials.
request:addHeader("accept", "application/json")
request:addHeader("Authorization", "Bearer " .. "TOKEN")

-- Prepare setting data.
request:setContentType("application/json")
request:setContentBuffer('{"apn": "internet.ht.hr","enabled": true,"connectionEnabled": true,"apnPassword": "","apnUsername": "","dialString": "","persistConnection": true}')

-- Execute.
local response = client:execute(request)
-- Get status code & json response. 
local response_code = response:getStatusCode()
local response_body = response:getContent()
```

Code `204` is expected when success.



<!-- Source: docs/code-samples/wpan-examples.md -->


# Wireless Personal Area Network Examples

In this section,
 the Wireless Personal Area Network (WPAN) service on the TDC-E device
 is discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example


In this section, a Go gRPC application is created and documented.
 A gRPC client is created from a Proto file to match the gRPC server the {Devic
eName()} device is serving. An application that checks device availability,
 connects to a device, and reads data from it is provided.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.
org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

To implement a gRPC client,
 a Proto file matching the gRPC server's specifications is needed.
 This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

GO111MODULE=on go install github.com/bufbuild/buf/cmd/buf@v1.34.0

buf build pkg/
buf generate pkg/
```

This installs needed `Protoc` files, and generates files for the gRPC service.

First, dial options are set to dial a UNIX socket:

```go
dialOptions := []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return net.Dial("unix", socketPath)
	}),
}
```

The application then creates a gRPC client that connects to the service unix so
cket.

```go
conn, err := grpc.NewClient("passthrough:///unixsocket", dialOptions...)
if err != nil {
	log.Fatalf("failed to connect to Unix socket: %v", err)
}
defer conn.Close()

client := protos.NewAWpanClient(conn)
ctx = metadata.NewOutgoingContext(ctx, metadata.New(nil))
```

A Wpan instance is created. It's passed the created client and context. 

Providing WPAN discovery is on, WPAN devices to connect to are listed.
 To see how to turn the **device discovery** process on,
 go to the [next section](#12-proto-file).
device doesn't appear on the list.

Set the `deviceMacAddress` variable to check it's device state,
 to connect with the device, and to attach it to the HDI Input.


```go
wpanDevice := wpan.NewWpan(ctx, client)
wpanDevice.ListWPANDevices()

time.Sleep(time.Second * 2)

if wpanDevice.CheckDevice(deviceMacAddress) {
	wpanDevice.ConnectDevice(deviceMacAddress)
	wpanDevice.AttachToHDIInput(deviceMacAddress)
}
```

If the device MAC address is not discovered with `CheckDevice`,
 the application prints an appropriate log and exits. Otherwise,
 the Go application proceeds to pair with and connect to the device. 

> **ℹ️ Info**
>
> Pairing is a stream in which the pairing status is possibly preceded by displ
> aying a `pass key`. Pairing usually precedes connecting to the device,
> but it's also possible that pairing is handled automatically when trying to es
> tablish a connection.


Connecting is done like this:

```go
func (t *Wpan) ConnectDevice(macAddress string) {
	req := protos.DeviceConnection{
		Address: macAddress,
	}
	_, err := t.Client.Connect(t.Ctx, &req)
	if err != nil {
		log.Fatalf("Failed to connect %s: %v", macAddress, err)
	}
}
```

the WPAN device to the console. 

```go
func (t *Wpan) AttachToHDIInput(macAddress string) {
	req := protos.DeviceAddress{
		Address: macAddress,
	}
	stream, err := t.Client.AttachToHIDInput(t.Ctx, &req)
	if err != nil {
		log.Fatalf("Failed to pair %s: %v", macAddress, err)
	}
	for {
		inputEvent, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Stream closed.")
				break
			}
			log.Fatalf("Error receiving input event: %v", err)
		}
		log.Printf("Received input event: %+v\n", inputEvent)
	}
}
```

See an example print of the data stream after connecting to a device below:

```log
2025/01/10 19:09:43 ----
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:747679000} type:4 code:4 value:458835
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:747679000} type:1 code:69 value:1
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:747679000}
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:758890000} type:17 value:1
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:758890000} type:4 code:4 value:458835
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:758890000} type:1 code:69
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:758890000}
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:770142000} type:4 code:4 value:458835
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:770142000} type:1 code:69 value:1
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:770142000}
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:780145000} type:17
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:780145000} type:4 code:4 value:458835
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:780145000} type:1 code:69
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:780145000}
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:781398000} type:4 code:4 value:458784
2025/01/10 19:09:55 Received input event: time:{sec:1736536195 nsec:781398000} type:1 code:4 value:1
```

### 1.2. Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the WPAN service is the following:

```bash
/**
 * WPAN service.
 *
 * Service that enables WPAN functionalities using bluez.
 */
syntax = "proto3";

package hal.wpan;
 
import "google/protobuf/empty.proto";
 
option go_package = "./protos;protos";

/**
 * Represents WPAN service status.
 */
message ServiceStatusData {
   bool enabled = 1; 
}

/**
 * Represents type of the scan.
 */
enum ScanType {
   Auto = 0;
   BrEdr = 1;
   Le = 2;
}

/**
 * Represents discovery filters.
 */
message DiscoveryFilters {
   repeated string uuids = 1;
   optional int32 rssi = 2;
   optional int32 pathloss = 3;
   optional ScanType transport = 4;
   optional bool duplicateData = 5;
   optional bool discoverable = 6;
   optional string pattern = 7;
}

/**
 * Represents device data.
 */
message DeviceData {
   optional string alias = 1;
   optional string address = 2;
   repeated string uuids = 3;
   optional int32 rssi = 4;
   optional bool paired = 5;
   optional bool bonded = 6;
   optional bool connected = 7;
   optional bool trusted = 8;
   optional bool blocked = 9;
   optional bool servicesResolved = 10;
}

/**
 * Represents devices.
 */
message Devices {
   repeated DeviceData devices = 1;
}

/**
 * Represents device address.
 */
message DeviceAddress {
   // Device MAC address.
   string address = 1;

   // Optional device name.
   optional string name = 2;
}

/**
 * Represents device connection request.
 */
message DeviceConnection {
   string address = 1;
   optional string profile = 2;
}

/**
 * Represents type of pairing process display.
 */
enum DisplayType {
   Pincode = 0;
   Passkey = 1;
}

/**
 * Represents device pairing process display.
 */
message PairingProcessDisplay {
   DisplayType type = 1;
   string content = 2;
}

/**
 * Represents device pairing status.
 */
message PairingStatus {
   bool paired = 1;
   string errMsg = 2;
}

/**
 * Represents device pairing response.
 */
message PairingResponse {
   oneof response {
      PairingStatus status = 1;
      PairingProcessDisplay pairing = 2;
   }
}

/**
 * Represents input HID device basic information.
 */
message HidDeviceData {
   string address = 1;
   string name = 2;
}

/**
 * Represents array of input HID devices.
 */
message HidDevices {
   repeated HidDeviceData hidDevices = 1;
}

/**
 * Represents syscall UNIX time value.
 */
message SyscallTime {
   int64 sec = 1;
   int64 nsec = 2;
}

/**
 * Represents input event from input HID device.
 * HID input events based on:
 *  - https://raw.githubusercontent.com/torvalds/linux/v6.7/include/uapi/linux/input.h
 *  - https://raw.githubusercontent.com/torvalds/linux/v6.7/include/uapi/linux/input-event-codes.h
 */
message InputEvent {
   // Time in seconds since epoch at which event occurred.
   SyscallTime time = 1;

   // Event type.
   int32 type = 2;

   // Event code related to the event type.
   int32 code = 3;

   // Event value related to the event type.
   int32 value = 4;
}

/**
 * Represents scanned device event.
 */
enum ScannedDeviceEvent {
   Added = 0;
   Removed = 1;
   Updated = 2;
}

/**
 * Represents extended device data.
 */
message ExtendedDeviceData {
   // Base device data.
   DeviceData device = 1;

   // Advertised transmitted power level.
   optional int32 txPower = 2;

   // Manufacturer specific advertisement data.
   message ManufacturerData {
      // manufacturer ID
      int32 ID = 1;

      // manufacturer data
      bytes data = 2;
   }
   repeated ManufacturerData manufacturerData = 3;

   // Service advertisement data.
   message ServiceData {
      // service UUID
      string uuid = 1;

      // service data
      bytes data = 2;
   }
   repeated ServiceData serviceData = 4;

   // Advertising data flags of the remote device.
   optional bytes advertisingFlags = 5;

   // Advertising data of the remote device.
   message AdvertisingData {
      // AD type
      int32 adType = 1;

      // advertising data
      bytes data = 2;
   }
   repeated AdvertisingData advertisingData = 6;
}
 
/**
 * Represents extended device data with scanned device event.
 */
message ExtendedDeviceDataEvent {
   // Scanned device event.
   ScannedDeviceEvent event = 1;

   // Scanned extended device data.
   ExtendedDeviceData device = 2;
}

/**
 * Represents HID input events description for connected device.
 */
message HIDDescription {
   // Vendor information.
   optional int32 vendor = 1;

   // Product information.
   optional int32 product = 2;

   // Version information.
   optional int32 version = 3;

   // Represents available device evdev code list based on evdev type.
   message EvdevTypeCodes {
      // Evdev type:
      int32 evType = 1;

      // Evdev codes list.
      repeated int32 evCodes = 2;
   }
   repeated EvdevTypeCodes typeCodes = 4;

   // Represents details on ABS input types.
   message EvdevAbsInfo {
      // ABS value.
      int32 value = 1;

      // ABS minimum value.
      int32 minimum = 2;

      // ABS maximum value.
      int32 maximum = 3;

      // ABS fuzz value.
      optional int32 fuzz = 4;

      // ABS flat value.
      optional int32 flat = 5;

      // ABS resolution value.
      optional int32 resolution = 6;
   }

   // Represents details on ABS input types based on evdev code.
   message EvdevCodeAbsInfos {
      // Evdev code:
      int32 evCode = 1;

      // ABS input type info.
     EvdevAbsInfo absInfo = 2;
   }
   repeated EvdevCodeAbsInfos codeAbsInfos = 5;

   // Represent evdev properties.
   repeated int32 evProperties = 6;
}

/**
 * Service exposing WPAN functions.
 */
service WPAN {
  // Used to enable/disable WPAN on device.
  rpc ToggleService(ServiceStatusData) returns (google.protobuf.Empty) {}

  // Used to get WPAN service status.
  rpc ServiceStatus(google.protobuf.Empty) returns (ServiceStatusData) {}

  // Used to start WPAN discovery.
  rpc StartDiscovery(DiscoveryFilters) returns (google.protobuf.Empty) {}

  // Used to stop WPAN discovery.
  rpc StopDiscovery(google.protobuf.Empty) returns (google.protobuf.Empty) {}

  // Used to start advanced WPAN discovery.
  rpc StartDiscoveryStream(DiscoveryFilters) returns (stream ExtendedDeviceDataEvent) {}

  // Used to retrieve found (and stored) devices.
  rpc GetDevices(google.protobuf.Empty) returns (Devices) {}

  // Used to retrieve device with specified address.
  rpc GetDevice(DeviceAddress) returns (DeviceData) {}

  // Used to remove all devices.
  rpc RemoveDevices(google.protobuf.Empty) returns (google.protobuf.Empty) {}

  // Used to remove device with specified address.
  rpc RemoveDevice(DeviceAddress) returns (google.protobuf.Empty) {}

  // Used to pair to device with specified address.
  rpc Pair(DeviceAddress) returns (stream PairingResponse) {}

  // Used to cancel pairing to device with specified address.
  rpc CancelPairing(DeviceAddress) returns (google.protobuf.Empty) {}

  // Used to connect with device with specified address.
  rpc Connect(DeviceConnection) returns (google.protobuf.Empty) {}

  // Used to disconnect from device with specified address.
  rpc Disconnect(DeviceConnection) returns (google.protobuf.Empty) {}

  // Used to get all HID inputs connected via WPAN HID profile.
  rpc GetHIDInputs(google.protobuf.Empty) returns (HidDevices) {}

  // Used to get HID description for a specific HID input.
  rpc GetHIDDescription(DeviceAddress) returns (HIDDescription) {}

  // Used to attach to and listen for input events from HID input.
  rpc AttachToHIDInput(DeviceAddress) returns (stream InputEvent) {}
}
```

</details>

#### 1.2.1. gRPC services

To get the WPAN service status, use `grpcurl`, which is an open-
source utility for accessing gRPC services via the shell.
 For help setting up the `grpcurl` command, refer to [gRPC Usage](/getting-
started/grpc-usage).

Use the following line: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> hal.wpan.WPAN.ServiceStatus
```

The `token` field is the fetched TDC-E authorization token.
 For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage).
 The `device-ip:
grpc_server_port` is the TDC-E IP address and the gRPC serving port.
 For example,
 if the `token` value was `token` and the address and port were `192.168.0.100:
8081`,
 you would use the following line to see the device status of the WPAN service.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.wpan.WPAN.ServiceStatus
```

The response should be in this format:

```json
{
  "enabled": true
}
```

If not enabled, use the following line to do so:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext -d '{"enabled":true}' 192.168.0.100:8081 hal.wpan.WPAN.ToggleService
{}
```

To start the discovery process of WPAN devices,
 use the following `grpcurl` command:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext -d '{"uuids":["00001101-0000-1000-8000-00805f9b34fb"]}' 192.168.0.100:8081 hal.wpan.WPAN.StartDiscovery
{}
```

This ensures the discovery process starts early.
 This particular command searches for all devices which use the `00001101-0000-
1000-8000-00805f9b34fb` UUID, which serves as a filter for WPAN devices.
 To stop the discovery process, use the following command:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.wpan.WPAN.StopDiscovery
{}
```

Additionally,
 you can use the `gRPC Clicker` VSCode extension for working with gRPC services
. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-
usage).

> **ℹ️ Info**
>
> If experiencing gRPC service timeouts when accessing a WPAN device,
> you can resolve this by adding a `max_time` parameter to the gRPC call.
> Below is an example of how to set the timeout to 10 seconds.
> Adjust the timeout as needed.


```bash
grpcurl -d '{"address":"MAC_ADDR"}' -H 'Authorization: Bearer token' -plaintext -max-time 10 192.168.0.100:8081 hal.wpan.WPAN/GetDevice
```

### 1.3. Application Deployment

#### 1.3.1. Dockerfile

To deploy the application,
 a Go container should be created and deployed to the TDC-E.
 To that end, a `Dockerfile` is created. The file is shown below.

```docker
# build image for Go app
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# setting environment
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0

RUN go build -o wpan-grpc ./cmd/main.go

# runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wpan-grpc .

CMD ["./wpan-grpc"]
```

Open a terminal and paste the following commands:

```bash
docker build -t wpan-grpc-app .
docker save -o wpan-grpc-app.tar wpan-grpc-app:latest
```

This will build the docker container and save the application as a `.
tar` file which can be used for Portainer upload.

#### 1.3.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go appli
cation. It sets the working directory as `/app`, then copies the `go.
mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-
E device has the `arm32hf` architecture and is based on `Linux`.
 With this in mind, the application is set to the following:

```dockerfile
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
CGO_ENABLED=0
```

The application is built as `wpan-grpc`.

```dockerfile
RUN go build -o wpan-grpc ./cmd/main.go
```

Finally,
 a runtime image is created from the latest `alpine` version for a smaller imag
e size. The working directory is once again set to `/app`,
 and the application is copied. The last line specifies that, upon deployment,
 the `wpan-grpc` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wpan-grpc .

CMD ["./wpan-grpc"]
```

#### 1.3.3. Deploying to Portainer

To deploy the application to the TDC-E device,
 `Portainer` can be used. To see instructions on the process,
 refer to [Working with Portainer](/getting-started/working-with-docker#3-
working-with-portainer). As soon as the image and container are set up,
 the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see re
> sults!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## 2. Node-RED gRPC Example

```c
Download Example Code from embedded file: "grpc-wpan.json"
```

The `gRPC` node set is not part of the initial Node-
RED package list and will have to be installed to the Palette.
 For `gRPC` node installation,
 import the following file in the **Manage Palette** section:
```c
Download Node from embedded file: "node-red-contrib-grpc-1.2.7.tgz"
```

This section describes usage of the WPAN service on the TDC-E. A Node-
RED gRPC application was created to demonstrate the following WPAN functionalit
ies:
- getting / setting service status
- starting / stopping discovery 
- listing devices
- listing device parameters
- pairing device
- connecting device
- attaching HID inputs from device

For implementation, the following nodes are used:

- `inject` node
- `gRPC call` node
- `debug` node.

The `gRPC` node server needs to be set properly to be able to connect to the gR
PC server and interpret server results correctly.
 The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node.
 The result should appear in the `debug` node. 

See a screenshot of the application below:

![wpan Node-RED Example](static/img/node-red-wpan.png)

## 3. Lua Example (Not Supported)

WPAN interface is **not supported** in Lua.


<!-- Source: docs/app-engine-docs/app-engine-docs.md -->

<>
    <iframe class="docs-custom-content" src="../../app-engine-docs/"/>
</>


