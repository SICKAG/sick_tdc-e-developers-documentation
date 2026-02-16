---
title: Controller Area Network Examples
sidebar_position: 206
---

import Admonition from '@theme/Admonition';


# Controller Area Network Examples

In this section, the Controller Area Network (CAN) interface of the TDC-E device and examples of its usage are discussed. Programming examples are given and thoroughly explained.

## 1. Setting up CAN Device

In this section, setting up the CAN device is discussed. Using the dedicated CAN HAL service, setup of the following parameters is possible:
 
 - Transceiver Power
 - Termination
 - Interface Mapping (namespace)

In the following sections, setting up the transceiver power, termination, and CAN namespace is discussed. Additionally, reading CAN statistics is described. This is done using the CAN HAL Service `hal.can.Can`. Examples using `grpcurl` are given below. For more information about using gRPC services, refer to [gRPC Usage](/getting-started/grpc-usage).

### 1.1. Setting Up CAN Transceiver Power

The CAN HAL service provides a means to set up transceiver power for the CAN interface. The HAL service `hal.can.Can.SetTransceiverPower` is used. An example of setting the transceiver power of the connected CAN device on is given below.

```bash
grpcurl -d '{"interfaceName":"CAN1","powerOn":true}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/SetTransceiverPower
{}
```

Checking changes to the CAN state is done by using the `hal.can.Can.GetTransceiverPower` HAL service.

```bash
grpcurl -d '{"interfaceName":"CAN1"}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/GetTransceiverPower
{
  "powerOn": true
}
```

### 1.2. Setting Up CAN Termination

The CAN HAL service provides a means to set up CAN termination. The HAL service `hal.can.Can.SetTermination` is used. An example of setting the termination of the connected CAN device on is given below.

```bash
grpcurl -d '{"interfaceName":"CAN1","enableTermination":true}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/SetTermination
{}
```

Checking changes to the termination is done by using the `hal.can.Can.GetTermination` HAL service.

```bash
grpcurl -d '{"interfaceName":"CAN1"}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/GetTermination
{
  "terminationEnabled": true
}
```

### 1.3. Setting Up CAN Interface Mapping

When working with the CAN interface, one can use different namespaces. The CAN interface is then bound to a namespace and can be accessed from there. Examples include:
- AppEngine
- Host
- Any other running container

To set the namespace of the CAN interface, the HAL service `hal.can.Can.SetInterfaceToContainer` is used. Possible values for the `dockerContainerName` parameter include:
- `app-engine` - maps the CAN interface to AppEngine
- empty string - maps the CAN interface to the host
- `container-name` - maps the CAN interface to a running Docker container

See an example of exposing the CAN interface to the host below.

```bash
grpcurl -d '{"interfaceName":"CAN1","dockerContainerName":""}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/SetInterfaceToContainer
{}
```

To see the changes to the interface mapping, the `hal.can.Can.GetInterfaceTocontainerMapping` service is used.

```bash
grpcurl -d '{}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/GetInterfaceToContainerMapping
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

### 1.4. Viewing CAN Statistics

To see statistics of your CAN device, the HAL service `hal.can.Can.GetStatistics` is used. See an example of its usage below.

```bash
grpcurl -d '{"interfaceName":"CAN1"}' -H 'Authorization: Bearer token' -plaintext 192.168.0.100:8081 hal.can.Can/GetStatistics
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

### 1.5. Proto File

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
  repeated string interfaces = 1; // List of available interface names (e.g., ["CAN1", "CAN2"])
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
  string dockerContainerName = 2;  // e.g., "app-engine", leave empty to move to host
}

// Represents a mapping of interface to docker container
message InterfaceToContainerMapping {
  string interfaceName = 1;  // e.g., "CAN1"
  string dockerContainerName = 2;  // e.g., "app-engine", if empty string then interface is mapped to host
}

// Response contains mapping of interfaces to docker containers
message GetInterfaceToContainerMappingResponse {
  repeated InterfaceToContainerMapping items = 1; // List of interface to docker container mappings 
}

/**
 * Service exposing CAN functions.
 */
service Can {
  // Gets the state of transceiver power
  rpc GetTransceiverPower(GetInterfaceNameRequest) returns (GetTransceiverPowerResponse);

  // Gets the state of termination
  rpc GetTermination(GetInterfaceNameRequest) returns (GetTerminationResponse);

  // Sets the transceiver power (on/off)
  rpc SetTransceiverPower(SetTransceiverPowerRequest) returns (google.protobuf.Empty) {}

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
  rpc SetInterfaceToContainer(SetInterfaceToContainerRequest) returns (google.protobuf.Empty) {}

  // Returns mapping of interfaces onto docker containers
  rpc GetInterfaceToContainerMapping(google.protobuf.Empty) returns (GetInterfaceToContainerMappingResponse) {}
}
```

</details>

## 2. Go Example


A Go application is provided as a CAN usage example. The application creates a CAN bus, binds it to a CAN port, and sends and receives CAN data simultaneously. For application testing, the `Kvaser Leaf Light v2` device was used. For generating and testing data, Kvaser's `CanKing` application and drivers were used.

### 2.1. Application Implementation

The application is implemented using the [\"canbus\" Go package](https://pkg.go.dev/github.com/go-daq/canbus@v0.2.0) to communicate with the CAN bus. 

> **📝 Note**
>
> Check your CAN usage namespace, and check if the CAN link is up. Check the [previous section](#13-setting-up-can-interface-mapping) to see CAN namespace assignment. If the CAN link is already enabled, skip the next step, and run the application normally.


If the CAN link is down, whilst running the application, add a `--setup=true` flag to the starting arguments, which will set up the CAN interface according to specifications in the `/pkg/canSetup` file. The setup sets the CAN bitrate, powers on the transceiver and brings the interface link up.

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

The bus is then bound to the `can0` interface.

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

![Sending Data from CAN](/img/cansendGo.PNG)

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

### 2.2. Application Deployment

#### 2.2.1. Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.


Open a terminal and paste the following commands:

```bash
docker build -t can-app .
docker save -o can-app.tar can-app:latest
```

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 2.2.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```


The application is built as `can`.

```dockerfile
RUN go build -o can ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `can` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/can .

CMD ["./can"]
```

#### 2.2.3. Deploying to Portainer

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

## 3. NodeRED Example

<a href="../code/node-red-examples/can.json" download>Download Example Code</a>

The following NodeRED example will utilize the NodeRED palette [\"node-red-contrib-socketcan\"](https://flows.nodered.org/node/node-red-contrib-socketcan) to communicate with the CAN bus. 

To start with the example, download NodeRED from the Applications tab. For more information, refer to [Installing, Accessing and Removing Applications](../getting-started/installing-accessing-removing-applications).


### 3.1. NodeRED Setup

Next open NodeRED's web user interface and in the top-right corner click on the menu icon next to the "Deploy" button, and select "Manage palette". Next click on "Install" card and search for  "node-red-contrib-socketcan" and install the palette.

![NodeRED manage palette](/img/can_nodered_socketcan_install.png)

Import the CAN example code by pressing `Ctrl + i` or by right clicking a flow → `Insert` → `Import`.

### 3.2. Flow Usage & Breakdown

After successfully importing the example flows, you should have a workspace similar to that of image below:

![NodeRED example flows A](/img/can_nodered_flows_a.png)

The example is composed of two flows, a `#1 Step - CAN Setup` and a `#2 Step - CAN Sender` flow. We  Use flow #1 to setup the CAN device on a container level, and we use flow #2 to demonstrate the usage of said CAN device.

`#1 Step - CAN Setup` flow is composed of two parts, a automatic flowchart (top section) and a manual debugging flow. Use the automatic flow to setup everything for CAN by clicking on the blue square on the `Setup CAN!` node. All of the red nodes in this flow run one of the following shell commands:

- `ip link` - lists all network devices
- `ip link set can0 [up/down]` - enables/disables the can0 device 
- `ip link set can0 type can bitrate 500000` - sets the bitrate to 500kbps, all other devices need the same bitrate to communicate
- `apk update && apk add iproute2` - get standard implementation of iproute2

> **📝 Note**
>
> By default, the `nodered` container already comes with the version of `ìproute2`, but it's minimal in features and doesn't have the ability to manage CAN devices out of the box. For this reason we are installing the standard implementation using `apk`.


![NodeRED example flows B](/img/can_nodered_flows_b.png)

`#2 Step - CAN Sender` is composed of 3 sections, a looping CAN sending flowchart (contains `Loop Trigger` node), two manual CAN sending flowcharts (inject nodes that start with _Send..._) and a CAN dump flowchart (contains `Pretty Print CAN Frame` node). You can freely use the inject nodes to try out the CAN functionality.
Look inside the manual sections inject nodes on possible ways to send payload to `socketcan-in` nodes or consult the palettes documentation section [Sending CAN frames](https://flows.nodered.org/node/node-red-contrib-socketcan).

> **📝 Note**
>
> The green indicator under the `socketcan-*` can be a bit misleading, since it can show as connected even if the CAN device is down. Verify if the device is active in flow #1 by triggering the `ip link` node.


After sending packets with any of the two CAN sending flowcharts, NodeREDs debug ouput and `candump` CLI tools output (from the `can-utils` package, as a alternative for `CanKing` application) should look like this:

![NodeRED CAN Outputs](/img/can_nodered_output.png)

### 3.3. Socketcan Node Configuration

If the device contains different CAN devices from can0, you can configure the socketcan palette nodes by editing both lime nodes with names `socketcan-*` and under the `Interface` field add can0, can1 or other CAN bus devices. For example if you wish to use the virtual device vcan0, you can configure the node like this:

![NodeRED socketcan configure](/img/can_nodered_socketcan_configure.png)

## 4. Lua Example


A Lua application is provided as a CAN usage example. The application creates a CAN handler, opens the handler and sends or receives CAN data. For application testing, the `Kvaser Leaf Light v2` device was used. For generating and testing data, Kvaser's `CanKing` application and drivers were used.


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

Reading is tested with the `CanKing` application. See a screenshot of sending CAN data below.

![Sending Data from CAN](/img/cansend.PNG)

To send data to the device, three test IDs and messages are created. The script then sends the defined data in an infinite loop.

```lua
Ids = {20, 21, 22}
Msgs = {"\x41\x42\x43\x44\x45\x46\x47\x48", "\x51\x52\x53\x54\x55\x56\x57\x58", "\x61\x62\x63\x64\x65\x66\x67\x68"}
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
