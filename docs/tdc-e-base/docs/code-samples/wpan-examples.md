---
title: Wireless Personal Area Network Examples
sidebar_position: 210
---

import Admonition from '@theme/Admonition';


# Wireless Personal Area Network Examples

In this section, the Wireless Personal Area Network (WPAN) service on the TDC-E device is discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that checks device availability, connects to a device, and reads data from it is provided.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

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

The application then creates a gRPC client that connects to the service unix socket.

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

Providing WPAN discovery is on, WPAN devices to connect to are listed. To see how to turn the **device discovery** process on, go to the [next section](#12-proto-file). It's a good practice to refresh the discovery process once in a while if your device doesn't appear on the list.

Set the `deviceMacAddress` variable to check it's device state, to connect with the device, and to attach it to the HDI Input.


```go
wpanDevice := wpan.NewWpan(ctx, client)
wpanDevice.ListWPANDevices()

time.Sleep(time.Second * 2)

if wpanDevice.CheckDevice(deviceMacAddress) {
	wpanDevice.ConnectDevice(deviceMacAddress)
	wpanDevice.AttachToHDIInput(deviceMacAddress)
}
```

If the device MAC address is not discovered with `CheckDevice`, the application prints an appropriate log and exits. Otherwise, the Go application proceeds to pair with and connect to the device. 

> **ℹ️ Info**
>
> Pairing is a stream in which the pairing status is possibly preceded by displaying a `pass key`. Pairing usually precedes connecting to the device, but it's also possible that pairing is handled automatically when trying to establish a connection.


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

The application then enters a streaming mode which prints all data received by the WPAN device to the console. 

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

To get the WPAN service status, use `grpcurl`, which is an open-source utility for accessing gRPC services via the shell. For help setting up the `grpcurl` command, refer to [gRPC Usage](/getting-started/grpc-usage).

Use the following line: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> hal.wpan.WPAN.ServiceStatus
```

The `token` field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage). The `device-ip:grpc_server_port` is the TDC-E IP address and the gRPC serving port. For example, if the `token` value was `token` and the address and port were `192.168.0.100:8081`, you would use the following line to see the device status of the WPAN service.

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

To start the discovery process of WPAN devices, use the following `grpcurl` command:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext -d '{"uuids":["00001101-0000-1000-8000-00805f9b34fb"]}' 192.168.0.100:8081 hal.wpan.WPAN.StartDiscovery
{}
```

This ensures the discovery process starts early. This particular command searches for all devices which use the `00001101-0000-1000-8000-00805f9b34fb` UUID, which serves as a filter for WPAN devices. To stop the discovery process, use the following command:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.wpan.WPAN.StopDiscovery
{}
```

Additionally, you can use the `gRPC Clicker` VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

> **ℹ️ Info**
>
> If experiencing gRPC service timeouts when accessing a WPAN device, you can resolve this by adding a `max_time` parameter to the gRPC call. Below is an example of how to set the timeout to 10 seconds. Adjust the timeout as needed.


```bash
grpcurl -d '{"address":"MAC_ADDR"}' -H 'Authorization: Bearer token' -plaintext -max-time 10 192.168.0.100:8081 hal.wpan.WPAN/GetDevice
```

### 1.3. Application Deployment

#### 1.3.1. Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.


Open a terminal and paste the following commands:

```bash
docker build -t wpan-grpc-app .
docker save -o wpan-grpc-app.tar wpan-grpc-app:latest
```

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 1.3.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```


The application is built as `wpan-grpc`.

```dockerfile
RUN go build -o wpan-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `wpan-grpc` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wpan-grpc .

CMD ["./wpan-grpc"]
```

#### 1.3.3. Deploying to Portainer

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](/getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## 2. Node-RED gRPC Example

<a href="../code/node-red-examples/grpc-wpan.json" download>Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>Download Node</a>

This section describes usage of the WPAN service on the TDC-E. A Node-RED gRPC application was created to demonstrate the following WPAN functionalities:
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

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

See a screenshot of the application below:

![wpan Node-RED Example](/img/node-red-wpan.png)

## 3. Lua Example (Not Supported)

WPAN interface is **not supported** in Lua.
