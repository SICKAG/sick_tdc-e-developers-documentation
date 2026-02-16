---
title: One-Wire Examples
sidebar_position: 11
---

import Admonition from '@theme/Admonition';


# One-Wire Examples

In this section, the One-Wire service on the TDC-E device is discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example

<a href="../code/go-examples/grpc-one-wire.zip" download>Download Example Code</a>

In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that reads status and temperature / data from the device is given.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file is then placed in a separate package, and from it, gRPC Go files are generated using the following commands:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/one-wire-service.proto
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

A ticker channel is created, and the application then enters a loop which continuously reads temperatures from all devices. The loop listens for a stop signal via `ctx.Done()`. On pressing `Ctrl+C`, the applications stops data collection and the one-wire service is disabled.

```go
for {
	select {
	case <-ctx.Done():
		log.Println("Disabling 1-wire service...")
		disableCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

The `Wire.go` package also contains example functions on how to read device data and temperature data from a single device.

```go
func (t *OneWire) ReadTemperature(ctx context.Context, deviceId string) (*gen.TemperatureReading, error) {
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

func (t *OneWire) ReadData(ctx context.Context, deviceId string, offset int32, length int32) (*gen.ReadDataResponse, error) {
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
2026/01/20 13:09:55 1-Wire Status: enabled:true device_count:1 master_path:"/sys/bus/w1/devices" active_devices:"10-000803675ae2"
2026/01/20 13:09:55 ----------------------------------------
2026/01/20 13:09:55 Available 1-Wire Devices: devices:{device_id:"10-000803675ae2" family_code:"10" serial_number:"000803675ae2" device_type:"DS18S20 (Temperature Sensor)" available:true last_seen:"2026-01-20T13:09:55Z"} 
2026/01/20 13:09:55 ----------------------------------------
2026/01/20 13:09:58 1-Wire Temperatures: readings:{device_id:"10-000803675ae2" temperature_celsius:85 timestamp:1768914598 valid:true} conversion_time_ms:879 
2026/01/20 13:09:58 ----------------------------------------
2026/01/20 13:10:00 1-Wire Temperatures: readings:{device_id:"10-000803675ae2" temperature_celsius:85 timestamp:1768914600 valid:true} conversion_time_ms:879 
2026/01/20 13:10:00 ----------------------------------------
2026/01/20 13:10:02 1-Wire Temperatures: readings:{device_id:"10-000803675ae2" temperature_celsius:85 timestamp:1768914602 valid:true} conversion_time_ms:895 
2026/01/20 13:10:02 ----------------------------------------
^C2026/01/20 13:10:03 Cancel signal received, shutting down...
2026/01/20 13:10:03 Disabling 1-wire service...
2026/01/20 13:10:03 1-Wire service disabled: success:true message:"disabled"
2026/01/20 13:10:03 ----------------------------------------
```

### 1.2. Proto File

<details>
<summary>Click to expand Proto file</summary>

The Proto file used for the OneWire service is the following:

```bash
syntax = "proto3";

package hal.onewire;

option go_package = "gitlab.sickcn.net/GBC05/devicebuilder/devices/tdc/services/1-wire.git/internal/service/gen";

service OneWire {
  rpc Enable(EnableRequest) returns (EnableResponse);
  rpc Disable(DisableRequest) returns (DisableResponse);
  rpc GetStatus(StatusRequest) returns (StatusResponse);
  rpc GetHotPlugStatus(HotPlugStatusRequest) returns (HotPlugStatusResponse);

  rpc ListDevices(ListDevicesRequest) returns (ListDevicesResponse);
  rpc GetDeviceInfo(DeviceInfoRequest) returns (DeviceInfo);

  // Real-time device event streaming via netlink uevents (instant hot-plug detection)
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
  string last_seen = 6;  // ISO 8601 formatted timestamp (e.g., "2026-01-12T10:30:45Z")
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
  string event_type_name = 4;  // Human-readable event type (e.g., "DEVICE_ADDED")
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

#### 1.2.1. gRPC services

To get the one-wire service status, use `grpcurl`, which is an open-source utility for accessing gRPC services via the shell. For help setting up the `grpcurl` command, refer to [gRPC Usage](../getting-started/grpc-usage).

Use the following line: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> hal.onewire.OneWire.GetStatus
```

The `token` field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](../getting-started/grpc-usage). The `device-ip:grpc_server_port` is the TDC-E IP address and the gRPC serving port. For example, if the `token` value was `token` and the address and port were `192.168.0.100:8081`, you would use the following line to see the status of the one-wire service.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.onewire.OneWire.GetStatus
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
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.onewire.OneWire.Enable
{
  "success": true,
  "message": "enabled"
}
```

To list all available devices, use the following line:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 hal.onewire.OneWire.ListDevices
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

Additionally, you can use the `gRPC Clicker` VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](../getting-started/grpc-usage).

### 1.3. Application Deployment

#### 1.3.1. Dockerfile

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

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 1.3.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

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

The application is built as `onewire-grpc`.

```docker
RUN go build -o onewire-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `onewire-grpc` application is started.

```docker
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/onewire-grpc .

CMD ["./onewire-grpc"]
```

#### 1.3.3. Deploying to Portainer

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](../getting-started/working-with-docker#2-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```docker
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## 2. Node-RED gRPC Example

<a href="../code/node-red-examples/grpc-one-wire.json" download>Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>Download Node</a>

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
- `inject` node
- `gRPC call` node
- `debug` node.

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

See a screenshot of the application below:

![One-Wire Node-RED Example](/img/node-red-one-wire.png)

## 3. Lua Example (Not Supported)

One-Wire interface is **not supported** in Lua.
