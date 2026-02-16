---
title: Networking
sidebar_position: 208
---

import Admonition from '@theme/Admonition';

# Networking

Networking REST API samples in Lua are prepared and shown below.

> **📝 Note**
>
> While the provided sample is written in **Lua**, the functionality is supported across multiple programming languages, including **Go and Node-RED**.


## 1. Lua Examples

In this section Lua code examples will be shown.

### 1.1. Control Center Admin Token

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
request:setContentBuffer('{"username":"admin","password":"myadminpassword","realm":"admin"}')
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

<a href="../code/lua-examples/networking_cc_token.lua" download>Download Example Code</a>

> **ℹ️ Info**
>
> Please download required helper script: <a href="../code/lua-examples/networking_util.lua" download>Util</a>.


### 1.2. Interface List

Get a list of all network interfaces.

The list is returned as a JSON array containing the name and all interface properties.

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

<a href="../code/lua-examples/networking_interface_list.lua" download>Download Example Code</a>

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

<a href="../code/lua-examples/networking_interface_settings.lua" download>Download Example Code</a>

### 1.4. Set Network Interface Mode

Set DHCP or static mode for a network interface.

When settings are accepted and applied, status code `204` is returned.

> **📝 Note**
>
> Value of all true or false type variables must be boolean.


#### 1.4.1. Set DHCP Mode

To set a DHCP mode for an interface, JSON with settings must be sent using the PUT HTTP method.

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

<a href="../code/lua-examples/networking_interface_mode.lua" download>Download Example Code</a>

#### 1.4.3. Set an Fallback IP Address

Fallback address can be set using the same code snippet as for setting dhcp or static IP, see sections: 

 - [Set DHCP mode](./networking#141-set-dhcp-mode)
 - [Set Static IP address](./networking#142-set-static-ip-address)


### 1.5. Setup WLAN Interface

Check sections below for further WLAN setup details.

<a href="../code/lua-examples/networking_wlan_interface.lua" download>Download Example Code</a>

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

Get a list of all available networks (saved or discovered) for one WLAN interface.

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

<a href="../code/lua-examples/networking_gateway.lua" download>Download Example Code</a>

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
- See: [Get gateway priority list](./networking#161-get-gateway-priority-list)


### 1.7. Modem

Modem status and setup examples.

<a href="../code/lua-examples/networking_modem.lua" download>Download Example Code</a>

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

