---
title: Serial Examples
sidebar_position: 3
---

import Admonition from '@theme/Admonition';


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


## 1. Go Example

This section handles setting up a serial device using gRPC, and creating serial Go applications. 

The first application is made and tested with the `Leaf Wetness Sensor`, which operates with a `Baudrate` of `9600`, using the communication protocol `MODBUS`, and with `RS485`. For this application, an `Isolated Soil Sensor` was used.

The second application was created by simulating a `RS422` device using two TDC-E devices and connecting their serial interfaces. The reading and writing functionalities are implemented as separate applications. The `RS232` device was created using loopback mode by connecting TX and RX wires.

### 1.1. Setting up Serial Device

In this subsection, setting up the serial device is discussed. Using the dedicated serial HAL service, setup of the following parameters is possible:
 
 - Mode
 - Transceiver Power
 - Slew Rate

In the following subsections, setting up the mode of a serial device is discussed. This is done using the Serial HAL Service `hal.serial.Serial`. Examples using `grpcurl` are given below. For more information about using gRPC services, refer to [gRPC Usage](/getting-started/grpc-usage).

#### 1.1.1. Setting Up Serial Mode

The Serial HAL service allows setting up the serial mode. The TDC-E device supports three modes across its two serial interfaces:

 - RS232 (SERIAL_1 only)
 - RS422 (SERIAL_2 only)
 - RS485 (SERIAL_2 only)

To set up your serial mode, use the `hal.serial.Serial.SetMode` service. Examples are given below.

**RS232 (SERIAL_1):**

```
grpcurl -emit-defaults -H 'Authorization: Bearer {token}' -plaintext -d '{"interfaceName":"SERIAL_1","mode":"RS232"}' 192.168.0.100:8081 hal.serial.Serial.SetMode
{}
```

**RS422 or RS485 (SERIAL_2):**

```
grpcurl -emit-defaults -H 'Authorization: Bearer {token}' -plaintext -d '{"interfaceName":"SERIAL_2","mode":"RS422"}' 192.168.0.100:8081 hal.serial.Serial.SetMode
{}
```

To check changes to the serial mode, use the HAL service `hal.serial.Serial.GetMode`.

```
grpcurl -emit-defaults -H 'Authorization: Bearer {token}' -plaintext -d '{"interfaceName":"SERIAL_2"}' 192.168.0.100:8081 hal.serial.Serial.GetMode
{
  "mode": "RS422"
}
```

#### 1.1.2. Viewing Serial Statistics

The Serial HAL service provides a means to view serial statistics. The HAL service `hal.serial.Serial.GetStatistics` is used. An example of seeing serial statistics for the RS422/RS485 interface (SERIAL_2) is given below.

```
grpcurl -d '{"interfaceName":"SERIAL_2"}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.serial.Serial.GetStatistics
{
  "txCount": "1050",
  "rxCount": "982"
}
```

For RS232 interface (SERIAL_1), use `"interfaceName":"SERIAL_1"`.

#### 1.1.3. Proto File

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
  repeated string slewRates = 1;   // Available slew rates (HIGH_SPEED and SLOW_SPEED)
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
  repeated string interfaces = 1; // List of available serial interface names (e.g., ["serial1", "serial2"])
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
  bool terminationEnabled = 1; // true if termination is enabled, otherwise false
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
  repeated string alternativeDevPaths = 3; // Alternative device paths in /dev/serial/by-id directory
  string description = 4; // Device description (dynamic devices might have generic simple description)
}

// List of serial interfaces with detailed information
message DetailedSerialDevices {
  repeated DetailedSerialDeviceData serialDevices = 1; // list of serial devices with detailed information
}

// SerialInterfaceService definition
service Serial {
  // Lists available serial interfaces
  rpc ListInterfaces(ListInterfacesRequest) returns (ListInterfacesResponse);

  // Gets the available Slew Rates of Serial interface
  rpc GetAvailableSlewRates(google.protobuf.Empty) returns (GetAvailableSlewRatesResponse);

  // Sets the transceiver power (on/off)
  rpc SetTransceiverPower(SetTransceiverPowerRequest) returns (google.protobuf.Empty);

  // Gets the transceiver power (on/off)
  rpc GetTransceiverPower(InterfaceNameRequest) returns (GetTransceiverPowerResponse);

  // Sets the mode of the serial interface (RS485 or RS422)
  rpc SetMode(SetModeRequest) returns (google.protobuf.Empty);

  // Gets the mode of the serial interface (RS485 or RS422)
  rpc GetMode(InterfaceNameRequest) returns (GetModeResponse);

  // Gets the available modes of the serial interface (RS485 and RS422)
  rpc GetAvailableModes(google.protobuf.Empty) returns (GetAvailableModesResponse);

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
  rpc ListDetailedInterfaces(google.protobuf.Empty) returns (DetailedSerialDevices);
}
```

</details>

### 1.2. RS485 Example

In this section, the RS485 application is implemented.

<a href="../code/go-examples/modbus-serial-sensor-tdce.zip" download>Download Example Code</a>

#### 1.2.1. Application Implementation

The application uses a single `.go` file to run. The `go.bug.st/serial` package is used to work with the serial port on `/dev/ttymxc1` (SERIAL_2 for RS485). The port mode is set and the communication to the port is opened.

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

A message for fetching the temperature and humidity of the sensor is created. For the `Isolated Soil Sensor`, a message in a specific format is sent to the serial port which prompts the sensor to return the required values.

To send a message to the serial port, the `Write` function is used.

```go
_, err := port.Write(message)
	if err != nil {
		log.Fatalf("Error writing to serial port: %v", err)
	}
```

The program then sleeps for a second so that the sensor has enough time to process the message and send data back to the host. For reading the serial port data, the `Read` function is used.

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

#### 1.2.2. Application Deployment

This section describes the Go application deployment.

**1.2.2.1. Dockerfile**

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.

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

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

**1.2.2.2. Dockerfile Breakdown**

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` file to the directory, then downloads necessary files.

```docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the `arm32hf` architecture and is based on `Linux`. With this in mind, the application is set to the following:

```docker
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0
```

The application image is built as `modbus-serial`.

```docker
RUN go build -o modbus-serial .
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `serial` application is started.

```docker
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/modbus-serial .

CMD ["./modbus-serial"]
```

**1.2.2.3. Deploying to Portainer**

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). 

As soon as the image and container are set up, the application starts running.


### 1.3. RS232 / RS422 Example

In this section, the RS232 / RS422 application is discussed. The RS232 and RS422 examples are mostly the same, with the exception of using the `/dev/ttymxc5` port for RS232, while RS422 uses the `/dev/ttymxc1` port.

For clarity, example usage is shown using RS422.

Download the sample applications below:

- RS232: <a href="../code/go-examples/serial-rs232-tdce.zip" download>Download Example RS232 Code</a>

- RS422: <a href="../code/go-examples/serial-rs422-tdce.zip" download>Download Example RS422 Code</a>

#### 1.3.1. Application Implementation

Both applications need to connect to the TDC-E device's port responsible for serial communication in order to read data sent via RS422. To achieve this, a function is created to connect to the `/dev/ttymxc1` (SERIAL_2) port using the `go.bug.st/serial` Go package.

This function configures the serial connection, setting the parity, data bits (8), stop bits, and baud rate. It then attempts to open the `/dev/ttymxc1` port. If successful, the port is returned; otherwise, an error is logged, and the application is terminated.

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

To send data through the port, the `writeData` function is used in the writer application:

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

This function accepts the port and a message in byte format. In this example, the application sends a `Hello world!` message every two seconds. If writing to the port fails, an error message is logged.

#### 1.3.2. Application Deployment

This section describes the Go application deployment.

**1.3.2.1. Dockerfile**

To deploy the applications, Go containers should be created and deployed to the TDC-E devices. To that end, a `Dockerfile` is created for both applications. As `Dockerfiles` are identical, this documentation will focus on showing a single `Dockerfile`, but to create two applications, make sure to rename the `serial` tag in the file. The file is shown below.

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

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

**1.3.2.2. Dockerfile Breakdown**

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` file to the directory, then downloads necessary files.

```docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Next, the Go environment is set. This needs to be done as the TDC-E device has the `arm32hf` architecture and is based on `Linux`. With this in mind, the application is set to the following:

```docker
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0
```

The application image is built as `serial`.

```docker
RUN go build -o serial .
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `serial` application is started.

```docker
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/serial .

CMD ["./serial"]
```

**1.3.2.3. Deploying to Portainer**

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer).

As soon as the image and container are set up, the application starts running.


## 2. Node-RED Examples 

A NodeRed application is provided as a Serial link usage example. Two scripts are given; one writes data to the TDC-E serial port, while the other reads from this port.

For Serial Out the palette `node-red-node-serialport` is required. Install it over the NodeRed UI.

### 2.1. Serial Node Example

<a href="../code/node-red-examples/serial-inout-tdce.json" download>Download Example Code</a>

#### 2.1.1. Write to Serial-Port

Drag a `Serial Out` node onto your flow.

Double click for settings:
- Add new Serial Port by clicking on the '+' in the Serial Port row
- The Serial Port is on `/dev/ttymxc5` for RS232 on SERIAL_1
- The Serial Port is on `/dev/ttymxc1` for RS422/RS485 on SERIAL_2
- Match the Settings to your Other Serial Device

Deploy the flow. If a "🟩 connected" appears under the `Serial Out` node, NodeRed successfully connected to the serial port.

Add an `Inject` node.

Edit the `Inject` node:
- Change payload to String, with a value of "Hello world!"
- Set Repeat to interval every 5 seconds

Connect the two Nodes and deploy your flow.

#### 2.1.2. Read from Serial-Port

Drag a `Serial In` node onto your flow.

Double click for settings:
- Add new Serial Port by clicking on the '+' in the Serial Port row
- The Serial Port is on `/dev/ttymxc5` for RS232 on SERIAL_1
- The Serial Port is on `/dev/ttymxc1` for RS422/RS485 on SERIAL_2
- Match the Settings to your Other Serial Device

Deploy the flow. If a "🟩 connected" appears under the `Serial Out` node, NodeRed successfully connected to the serial port.

Add a `Debug` node and connect the two Nodes and deploy your flow. Now all incoming serial data gets printed onto you debug window, which you can find in the top right corner.

![Serial Node-RED Example](/img/serial-node-red.png)

### 2.2. Serial gRPC Example

<a href="../code/node-red-examples/grpc-serial-tdce.json" download>Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>Download Node</a>

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

![Serial RS422 Node-RED Example](/img/grpc-serial.png)

## 3. Lua Example

A Lua application is provided as a Serial link usage example. Two scripts are given; one writes data to the TDC-E serial port, while te other reads from this port. For `RS485` communication, a `DIGITUS USB to Serial Adapter` is connected to the TDC-E device for reading and writing data. For `RS422` communication, two TDC-E devices are connected via a serial interface. For `RS232`, loopback is created by connecting the TX and RX wires.

> **⚠️ Warning**
>
> Termination for `SERIAL_1` and `SERIAL_2` cannot be set as it is unsupported! Setting termination will result in configuration failure.


To change between the `RS232`, `RS422` and `RS485` standards, find the following line of code in the `.lua` example and set it accordingly. For example, to set the `RS485` standard, use the following:

```lua
S2:setType("RS485")
```

### 3.1. Writing to the Serial Device

<a href="../code/lua-examples/serial-write-tdce.lua" download>Download Example Code</a>

The first `.lua` script writes data to the serial device. It creates a serial connection to the device. The type of the device is set and baud rate is set to `115200`.

```lua
S2 = SerialCom.create('SER2')
  
S2:setType("RS485")
S2:setBaudRate(115200)
```

A connection is opened and a `Hello world!` message is created. Then, a timer is created to send this message to the RS485 device periodically every 5 seconds. The message is transmitted the following way:

```lua
local Retb = S2:transmit(message)
print("Transmitted " .. Retb .. " bytes.")
```

Note that a timer is implemented instead of a `Sleep` service. This is because `Sleep` will cause all code to wait for the specified time, while timers operate locally, meaning that only the `rs-write.lua` application script will sleep for the specified time. 

The script prints the engine version at the end of the script. See the result of the application run below.

```
[15:52:06.041: INFO: AppEngine] Starting app: 'serial' (priority: LOW)
Transmitted 13 bytes.
Transmitted 13 bytes.
Transmitted 13 bytes.
```

### 3.2. Reading from the Serial Device

<a href="../code/lua-examples/serial-read-tdce.lua" download>Download Example Code</a>

The second `.lua` script reads data from the RS485 device. The script creates a serial connection and sets the baudrate to `115200`, then opens the connection. A `Callback` function is created to read from the device.

A `Callback` function is created to read from the device.

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
