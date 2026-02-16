---
title: Digital Input / Output Examples
sidebar_position: 201
---

import Admonition from '@theme/Admonition';


# Digital Input / Output Examples

In this section, the Digital Inputs/Outputs of the TDC-E device are discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. From it, the list of DIO devices, their input / output directions and states can be read and set. To that end, a simple Go application demonstrating said usage is created.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

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

### 1.2. Proto File

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
  rpc ListDevices(google.protobuf.Empty) returns (DigitalIOListDeviceResponse) {}

  /// Used to set pin direction value of a particular gpio pin.
  rpc SetDirection(DigitalIOSetDirectionRequest) returns (google.protobuf.Empty) {}

  /// Used to read value of a particular gpio pin.
  rpc Read(DigitalIOReadRequest) returns (DigitalIOReadResponse) {}

  /// Used to write value of a particular gpio pin.
  rpc Write(DigitalIOWriteRequest) returns (google.protobuf.Empty) {}

  /// Used to stream input events form device.
  rpc Attach(google.protobuf.Empty) returns (stream DigitalIOAttachResponse) {}
}
```

</details>

To list all available Digital Input/Output services, use `grpcurl`, which is an open-source utility for accessing gRPC services via the shell. For help setting up the `grpcurl` command, see [gRPC Usage](../getting-started/grpc-usage).

To list all available Digital Input/Output calls, use the following line:

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> list hal.digitalio.DigitalIO
```

The `token` field is the fetched TDC-E authorization token. For help fetching this token, refer to [gRPC Usage](../getting-started/grpc-usage). The `device-ip:grpc_server_port` is the TDC-E IP address and the gRPC serving port. For example, if the `token` value was `token` and the address and port were `192.168.0.100:8081`, you would use the following line to list all available Digital Input/Output services.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 list hal.digitalio.DigitalIO
```

The response should be in this format:

```log
hal.digitalio.DigitalIO.Attach
hal.digitalio.DigitalIO.ListDevices
hal.digitalio.DigitalIO.Read
hal.digitalio.DigitalIO.SetDirection
hal.digitalio.DigitalIO.Write
```

Additionally, you can use the `gRPC Clicker` VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](../getting-started/grpc-usage).

### 1.3. Application Deployment

#### 1.2.1. Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.


Open a terminal and paste the following commands:

```bash
docker build -t dio-grpc-app .
docker save -o dio-grpc-app.tar dio-grpc-app:latest
```

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 1.2.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```


The application is built as `dio-grpc`.

```dockerfile
RUN go build -o dio-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `dio-grpc` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dio-grpc .

CMD ["./dio-grpc"]
```

#### 1.2.3. Deploying to Portainer

To deploy the application to the TDC-E device, `Portainer` can be used. To see instructions on the process, refer to [Working with Portainer](../getting-started/working-with-docker#3-working-with-portainer). As soon as the image and container are set up, the application starts running.

> **📝 Note**
>
> Make sure to expose the HAL unix socket to the container to be able to see results!


Bind the following:

```dockerfile
/var/run/hal/hal.sock:/var/run/hal/hal.sock
```

## 2. Node-RED gRPC Example

<a href="../code/node-red-examples/grpc-dio.json" download="grpc-dio.json">Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download="node-red-contrib-grpc-1.2.7.tgz">Download Node</a>

This section describes the creation and usage of a Node-RED gRPC example using digital input / output devices. The application makes listing all DIO devices, reading DIO states, writing the DIO state and streaming changes in the DIO devices possible. 

For implementation, the following nodes are used:

- `inject` node
- `gRPC call` node
- `debug` node.

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

See a screenshot of the application below:

![DIO Node-RED Example](/img/dio-node-red.png)

## 3. Lua Example

<a href="../code/lua-examples/dio.lua" download="dio.lua">Download Example Code</a>

A Lua application is provided as a DIO usage example. The example demonstrates setting a DIO device as output or input, setting an output state and reading DI states.

The script prints the engine version at the start of the script. It then creates a `DIO AO` output and sets the value of the output to `HIGH`.

```lua
dioAO = Connector.DigitalOut.create('DIO_AO')
dioAO:set(true);
print("Set DO A output to HIGH.")
```
To set a DO output to `LOW`, uncomment the following lines:

```lua
dioAO:set(false);
print("Set DO A output to LOW.")
```

The script then sets three DIO devices to inputs (`DIO_BI`, `DIO_CI` and `DIO_DI`). The state of said inputs is then read using the following function:

```lua
function PrintDIOStateIn(dio, name)
  local state = dio:get()
  print(string.format("%s state: %s", name, state))
end
```

Another way of reading the state is registering the DI `onChangeStamped`, printing the DI state on event.

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
