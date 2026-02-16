---
title: gRPC Usage
sidebar_position: 103
---

import Admonition from '@theme/Admonition';


# gRPC Usage

This section provides insight into how to use the HAL Service gRPC API. Setting up `grpcurl`, server authorization, `gRPC Clicker` and accessing HAL services is discussed.

## 1. `grpcurl` Installation

The `grpcurl` utility provides a means to access gRPC services via shell command. It is used to send a request to the TDC-E device with the specified HAL service port. To get the latest version of the `grpcurl` utility, run the following line in the terminal:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

To be able to access the utility, the `bin path` needs to be exported.

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

To test `grpcurl` availability, use the following line:

```bash
grpcurl --version
```

## 2. Authentication

To access TDC-E' HAL Services, user authentication is needed. This is done by fetching an authentication token from the device. Use the `curl` command to authenticate the TDC-E user. See an example below.

```bash
curl -s -X POST http://{DEVICE_IP}/auth/login -H 'Accept: application/json' -H 'Content-Type: application/json' -d '{"username":"{USERNAME}","password":"{PASSWORD}","realm":"{REALM}"}'
```

Replace the `{DEVICE_IP}`, `{USERNAME}`, `{PASSWORD}` and `{REALM}` with your actual values. For example, if the values of the listed variables are the following:

<GrpcAuthentication />

The `curl` command will be the following:

<GrpcCurlExamle />

This will return a token response from the server which will be used for further authentication. 

> **📝 Note**
>
> To export the value of your token, you can install a service like `jq`. See the example of exporting a token with `jq` below.


```bash
export CC_TOKEN=$(curl -s -X POST http://{DEVICE_IP}/auth/login -H 'Accept: application/json' -H 'Content-Type: application/json' -d '{"username":"{USERNAME}","password":"{PASSWORD}","realm":"{REALM}"}' | jq -r '.token')
```

> **ℹ️ Info**
>
> The token is reset after its expiration time is out, and after each reboot of the device.


## 3. Working with `grpcurl`

To send gRPC requests to the TDC-E device, an authentication token is needed which is added to the gRPC curl. In other words, each call to the server needs to be authenticated. This is done by adding a header to the `grpcurl` command.

In the following examples, the `{DEVICE_IP}` and `{GRPC_SERVER_PORT}` will be set to `192.168.0.100` and `8081` respectively. The HAL Service port is set to `8081` by default, but can be configured in the **System Server Settings**.

To list all HAL services via the server reflection, use the following command:

```bash
grpcurl -H 'Authorization: Bearer {token}' -plaintext {DEVICE_IP}:{GRPC_SERVER_PORT} list
```


To list all functionalities of a single service, the following command is used:

```bash
grpcurl -H 'Authorization: Bearer {token}' -plaintext {DEVICE_IP}:{GRPC_SERVER_PORT} list {PACKAGE}.{SERVICE}
```


## 4. gRPC Clicker

The `gRPC Clicker` extension is a VSCode extension which can be installed via the _Extensions Marketplace_. It is used to connect to a gRPC server and to request services from it. It has a Graphical User Interface to make it easier for the user to interact with the server.

### 4.1. Installing gRPC Clicker

To install gRPC Clicker, open VSCode and go to the **Extensions** tab. Search for _gRPC Clicker_. Click _Install_.

![gRPC Clicker Extension](/img/grpc-clicker.PNG)
 
After installation, the tool should create a new tab which can be accessed. To add a new gRPC server, select the **plus** icon in the **Proto Groups** section. A **Proto schema assistant** will open.

The following parameters should be set up:

- `Name` - Name of your Proto Group
- `Address` - Device IP and gRPC server port
- `Custom flags` - Authorization header

For example, set the name to grpc-TDC-E, and the address and port to {selectByDevice({ SIM2000: "192.168.1.1:8081", default: "192.168.0.100:8081" })}. Custom flags need to be set as authentication is needed to request services from the gRPC TDC-E server. The custom flag has the following structure:

```h
-H 'Authorization: Bearer TOKEN'
```

> **📝 Note**
>
> Make sure to generate an authentication token and paste it instead of the `TOKEN` value!


The group is now set and can list all gRPC services currently available on your TDC-E device. 


### 4.1. Working with gRPC Clicker


## 5. Accessing HAL Service from Unix Socket

For internal usage, the TDC-E has a Unix socket which allows accessing available HAL Services. The HAL socket on the TDC-E device works with a TCP endpoint like `grpcurl`, so the command is generally used in the same way. 

> **📝 Note**
>
> Install `grpcurl` in your workspace to access socket content or use the examples provided in the **Code Samples** section.


For example, listing all HAL service commands is done using the following line:

```bash
grpcurl -emit-defaults -plaintext -unix -authority "localhost" "unix:///var/run/hal/hal.sock" list
```

Instead of a server address, `grpcurl` now uses the `-unix` flag and the path to the HAL socket. Since the HAL socket is run locally on the TDC-E device, no authentication form is required.

## 6. Additional Notes

For Node-RED examples used in this documentation, a **custom version** of the [node-red-contrib-grpc node](https://flows.nodered.org/node/node-red-contrib-grpc) was implemented. Using the Palette's version in Node-RED will **not allow TDC-E user identification and/or accessing the UNIX socket**! 

Make sure to import the following node into your Node-RED project:

<a href="../code/extras/node-red-contrib-grpc-1.2.7.tgz" download>node-red-contrib-grpc v1.2.7</a>
