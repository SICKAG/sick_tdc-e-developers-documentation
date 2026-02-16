---
title: Installing, Accessing and Removing Applications
sidebar_position: 100
---

import Admonition from '@theme/Admonition';


# Installing, Accessing and Removing Applications

This section discusses the setup of applications on the TDC-E device. The TDC-E device has some predefined applications that can be downloaded and used in simple steps. Those applications include:

- `Grafana`
- `Influxdb`
- `Mariadb`
- `Mosquitto`
- `Nodered`
- `Portainer`
- `Postgresql`
- `Redis`
- `Vscode`
- `Sftpgo`
- `Loki`
- `Fluentbit`
- `Sensor Data Gate (SDG)`
- `File Browser`

An example of installing `Portainer` will be provided.

## 1. Example: Setting Up Portainer

This section describes how to install, access and uninstall an application using the `Portainer` service as an example.

### 1.1. Installing Services

To set up Portainer, go to the [TDC-E Home Page](http://192.168.0.100/#/), then select the **System** tab on the left panel. Select **Applications**.

The **Applications** page lists services that can be installed to the TDC-E device by simply clicking on the **Install** button. To install Portainer, find the service in the given list and click **Install**.

> **ℹ️ Info**
>
> To install the listed applications, make sure the TDC-E device is connected to a network.


![Installing Applications](/img/install-applications.png)

The chosen service will be installed shortly. Once the download and installation is complete, the service will automatically be running. 

### 1.2. Accessing Services

To access the newly installed service, select the `Open` button. This will take you to the home page of the installed service. For example, [Portainer](https://192.168.0.100:9443/) is located on `https://192.168.0.100:9443/`.

![Portainer](/img/portainer-sample.png)

> **ℹ️ Info**
>
> Some services, like Portainer, require setting up a username and password before proceeding to the application itself.


### 1.3. Uninstalling Services

To uninstall an application from your device, select the **System** tab on the left panel. Select **Applications**. Select **Uninstall** for the application you want to remove and wait for the process to finish. This will remove the application and all user data on it from the device.

## 2. Uploading Custom Applications

A custom application can also be deployed to the device by selecting the `Upload custom application` button on **System > Applications**. This prompts the user for a `.tar`, `.tar.gz`, `.tar.bz2` or `.tar.xz` file.

The file should contain the following directories and files:

- `/compose` (required)
- `/images` (required)
- `/configs`
- `logo.png`
- `config.json`
- `VERSION`

### 2.1. Application File Directory Explanation

- **/compose (required)**

This directory contains at least one **docker-compose.yml** file, including other files neccessary in the same working directory for docker-compose (e.g .env files). The Docker compose container and service names **must be unique** among all applications on one system.

- **/images (required)**

This directory contains tarball files (.tar|.tar.gz|.tar.bz2|.tar.xz) with all the necessary docker images.

- **/configs** (optional)

This directory contains config files that will be placed into installation directory. These files can be refered to by using environment variable `$APP_INSTALLATION_DIR` and subdirectory configs.

See a `mosquitto.conf` placed into `/configs/mosquitto/config` directory example below:

```yaml
volumes:
  - mosquitto-data:/mosquitto/data
  - mosquitto-log:/mosquitto/log
  - $APP_INSTALLATION_DIR/configs/mosquitto/config/mosquitto.conf:/mosquitto/config/mosquitto.conf
```

- **logo.png** (optional)

Logo image that will be shown on device web GUI.

- **config.json** (optional)

This is the application configuration file in `json` format with following schema:

```json
{
  "autoStart": true,
  "uiEndpoint": 18000,
  "name": "My App",
  "description": "Application serves as an example on how to ...",
  "homepage": "www.sick.com"
}
```

The `autoStart` parameter specifies if application should start immediately after installation.

The `uiEndpoint` parameter specifies on which network endpoint application GUI will be available.

- **VERSION** (optional)

The textual file that contains version string of the application that will be reported by the API.
