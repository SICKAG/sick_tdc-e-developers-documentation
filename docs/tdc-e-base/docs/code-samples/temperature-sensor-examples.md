---
title: Temperature Sensor Examples
sidebar_position: 205
---

import Admonition from '@theme/Admonition';


# Temperature Sensor Examples

In this section, temperature sensors on the TDC-E device are discussed. Programming examples are given and thoroughly explained.

## 1. Go gRPC Example


In this section, a Go gRPC application is created and documented. A gRPC client is created from a Proto file to match the gRPC server the TDC-E device is serving. An application that lists all temperature sensors and prints their current value is created.

The application uses Golang's [gRPC service](https://pkg.go.dev/google.golang.org/grpc) to fetch and send data to the server.

### 1.1. Application Implementation

To implement a gRPC client, a Proto file matching the gRPC server's specifications is needed. This Proto file used to generate gRPC Go files using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=pkg/pb --go-grpc_out=pkg/pb pkg/temperature-sensor-service.proto
```

This installs needed tools, and generates files for the gRPC service. The application also requires an access token to connect to the gRPC service the TDC-E provides. Please refer to [gRPC Usage](/getting-started/grpc-usage) for instructions on how to generate an access token for the TDC-E.

Once you've obtained the token, navigate to `grpc-temperature/pkg/auth/token.json`. Paste your generated token in the `access_token` field. This token is then used to create a `context` which will be used to create gRPC calls.

The main application first creates a new gRPC client that connects to the service unix socket. It does so by using the generated gRPC Go file specification. 

```go
conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
```

To demonstrate the listed functionalities, the application lists all available temperature sensors on your TDC-E device, and then the temperature value of each specific sensor measured in `°C` is printed.

To list all temperature devices, an empty `pb` request is sent to the server. The client specifies the function `ListDevices` to receive all devices from the TDC-E device, which are then printed to the user.

```go
req := &emptypb.Empty{}
res, err := client.ListDevices(context.Background(), req)
```

Afterwards, each sensor's temperature is read. 

```go
func readTemperature(ctx context.Context, client protos.TemperatureSensorClient, device string) {
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


### 1.2. Proto File

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
  rpc ListDevices(google.protobuf.Empty) returns (TemperatureSensorListDeviceResponse) {}

  /// Used to read value of a particular temperature sensor.
  rpc Read(TemperatureSensorReadRequest) returns (TemperatureSensorReadResponse) {}
}
```

</details>

To list all available temperature sensor services, use `grpcurl`, which is an open-source utility for accessing gRPC services via the shell. For help setting up the `grpcurl` command, refer to [gRPC Usage](/getting-started/grpc-usage).

Use the following line for listing them: 

```bash
grpcurl -expand-headers -H 'Authorization: Bearer <token>' -emit-defaults -plaintext <device_ip>:<grpc_server_port> list hal.temperaturesensor.TemperatureSensor
```

The `token` field is the fetched TDC-E authorization token. For help fetching this token, see [gRPC Usage](/getting-started/grpc-usage). The `device-ip:grpc_server_port` is the TDC-E IP address and the gRPC serving port. For example, if the `token` value was `token` and the address and port were `192.168.0.100:8081`, you would use the following line to list all available temperature sensor devices.

```bash
grpcurl -expand-headers -H 'Authorization: Bearer token' -emit-defaults -plaintext 192.168.0.100:8081 list hal.temperaturesensor.TemperatureSensor
```

The response should be in this format:

```log
hal.temperaturesensor.TemperatureSensor.ListDevices
hal.temperaturesensor.TemperatureSensor.Read
```

Additionally, you can use the `gRPC Clicker` VSCode extension for working with gRPC services. For help setting the service up, refer to [gRPC Usage](/getting-started/grpc-usage).

### 1.3. Application Deployment

#### 1.2.1. Dockerfile

To deploy the application, a Go container should be created and deployed to the TDC-E. To that end, a `Dockerfile` is created. The file is shown below.


Open a terminal and paste the following commands:

```bash
docker build -t temp-grpc-app .
docker save -o temp-grpc-app.tar temp-grpc-app:latest
```

This will build the docker container and save the application as a `.tar` file which can be used for Portainer upload.

#### 1.2.2. Dockerfile Breakdown

The `Dockerfile` first creates a build image that is used to build the Go application. It sets the working directory as `/app`, then copies the `go.mod` and `go.sum` files to the directory and downloads all needed files.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```


The application is built as `temp-grpc`.

```dockerfile
RUN go build -o temp-grpc ./cmd/main.go
```

Finally, a runtime image is created from the latest `alpine` version for a smaller image size. The working directory is once again set to `/app`, and the application is copied. The last line specifies that, upon deployment, the `temp-grpc` application is started.

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/temp-grpc .

CMD ["./temp-grpc"]
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

<a href="../code/node-red-examples/grpc-temp.json" download>Download Example Code</a>

The `gRPC` node set is not part of the initial Node-RED package list and will have to be installed to the Palette. For `gRPC` node installation, import the following file in the **Manage Palette** section:
<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>Download Node</a>

This section describes usage of the temperature sensors on the TDC-E device. A Node-RED gRPC application was created to demonstrate listing all temperature sensor devices and reading temperature values from said devices.

For implementation, the following nodes are used:

- `inject` node
- `gRPC call` node
- `debug` node.

The `gRPC` node server needs to be set properly to be able to connect to the gRPC server and interpret server results correctly. The following configuration is used:

- `Server`:`unix:///var/run/hal/hal.sock`
- a `Proto File` is provided 

To test out any of the listed functionalities, use the `inject` node. The result should appear in the `debug` node. 

See a screenshot of the application below:

![AIN Node-RED Example](/img/TEMP-node-red.png)

## 3. Lua Example


The `extractTemperature` function checks whether the passed variable is a monitor. If not, the function returns the `nil` value. Otherwise, it uses the monitor's `get()` function to fetch the temperature value of the device.

```lua
local function extractTemperature(temperatureMonitor)
  if not temperatureMonitor then
    return nil
  end
  return temperatureMonitor:get()
end
```

See the result of the application run below.


