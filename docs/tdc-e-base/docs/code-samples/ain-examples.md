---
title: Analog Input Examples
sidebar_position: 202
---

import Admonition from '@theme/Admonition';


# Analog Input Examples

In this section, the Analog Inputs of the TDC-E device are discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. From it, the list of analog inputs and their current values can be read. The application also has an example of how to change the analog input mode.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

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
		log.Printf("%s Values:\nadcRaw: %d\nconverted: %f\nunit: %s\n", t.AIDevices[i], res.AdcRaw, res.Converted, res.Unit)
	}
}
```

The response is printed to the terminal and contains the following values:

- `adcRaw`
- `Converted`
- `Unit`

The program listens to a SIGINT and SIGTERM signal from the user, and will exit upon receiving it.

### 1.2. Proto File

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
 * Represents the request to set the AnalogIn service state for a specific channel.
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
 * Represents the request to get the AnalogIn service state for a specific channel.
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
  rpc ListDevices(google.protobuf.Empty) returns (AnalogInListDeviceResponse) {}

  /// Used to read value of a particular analog input channel.
  rpc Read(AnalogInReadRequest) returns (AnalogInReadResponse) {}
  
  /// Used to set measure mode of a particular analog input channel.
  rpc SetMeasureMode(AnalogInSetMeasureModeRequest) returns (google.protobuf.Empty) {}

  /// Used to monitor for overcurrent events
  rpc Attach(google.protobuf.Empty) returns (stream AnalogInAttachResponse) {}
}

```
</details>

To list all available Analog Input services, use `grpcurl`, which is an open-source utility for accessing gRPC services via the shell. For help setting up the `grpcurl` command, refer to [gRPC Usage](../getting-started/grpc-usage).

To list all available Analog Input calls, use the following line:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> list hal.analogin.AnalogIN
```

The `token` field is the fetched TDC-E authorization token. For help fetching this token, refer to [gRPC Usage](../getting-started/grpc-usage). The `device-ip:grpc_server_port` is the TDC-E IP address and the gRPC serving port. For example, if the `token` value was `token` and the address and port were `192.168.0.100:8081`, you would use the following line to list all available Analog Input services.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 list hal.analogin.AnalogIN
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

Additionally, you can use the `gRPC Clicker` VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](../getting-started/grpc-usage).


### 1.3. Application Deployment

#### 1.2.1. Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.


Open a terminal and paste the following commands:

```bash
docker build -t ain-grpc-app .
docker save -o ain-grpc-app.tar ain-grpc-app:latest
```

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 1.2.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```


The application is built as `ain-grpc`.

```dockerfile
RUN go build -o ain-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `ain-grpc` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ain-grpc .

CMD ["./ain-grpc"]
```

#### 1.2.3. Deploying to Portainer

To deploy the application to the TDC-E device, [Portainer](https://192.168.0.100:9443/) can be used. To see instructions on the process, refer to [Working with Portainer](../getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## 2. Node-RED gRPC Example

<a href="../code/node-red-examples/grpc-ain.json" download>Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>Download Node</a>

This section describes the creation and usage of a Node-RED gRPC example using analog inputs. The application makes listing all AI devices, reading AIN values, changing AIN mode and streaming changes in the AI devices possible. 

For implementation, the following nodes are used:

- `inject` node
- `gRPC call` node
- `debug` node.

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

See a screenshot of the application below:

![AIN Node-RED Example](/img/AIN-node-red.png)

## 3. Lua Example


See the example result of the application run below:


