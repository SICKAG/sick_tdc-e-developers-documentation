# SICK TDC-E Developer's Documentation

## 📘 TDC-E BaseOS Documentation (Current)

**👉 [View the Complete TDC-E BaseOS Developer Documentation](https://sickag.github.io/sick_tdc-e-developers-documentation/)**

The new TDC-E device with BaseOS documentation is available on GitHub Pages and includes:
- **Getting Started** - Installation, gRPC setup, Docker workflows
- **Code Samples** - Examples for CAN, DIO, GNSS, Serial, Networking, and more
- **AppEngine Documentation** - Lua crowns and capabilities
- **Interface Snippets** - Ready-to-use code for hardware interfaces

---

## 📙 TDC-E L4M Documentation (Legacy)

This is the [SICK TDC-E L4M Developer's Documentation](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki), aimed at customers and developers using the legacy L4M firmware. The wiki contains the following file structure:

* TDC-E Configuration
  * [L4M Configurations and Connections](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/TDC%E2%80%90E-L4M-Configuration-and-Connections)
  * [Network Configuration](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/TDC%E2%80%90E-Network-Configuration)
  * [Interface Configuration](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/TDC%E2%80%90E-Interface-Configuration)
* Getting Started
  * [Environment Setup](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup)
    * [Setting up Node-RED Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#1-setting-up-node-red)
    * [Setting up Python Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#2-setting-up-python-application-environment)
    * [Setting up C# Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#3-setting-up-c-application-environment)
    * [Setting up Go Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#4-setting-up-go-application-environment)
    * [Setting up MySql Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#5-setting-up-mysql-environment)
    * [Setting up MQTT (Mosquitto) Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#6-setting-up-mosquitto-environment)
  * [Image Building and Composing](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Build-and-Compose)
* Development Documentation
  * [Snippets Using the Hardware Interface](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/%5BHW-Interface%5D-Development-Documentation)
  * [Snippets Without the Hardware Interface](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/%5BNo-HW-Interface%5D-Development-Documentation)
  * [Project Examples Documentation](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/%5BExamples%5D-Development-Documentation)
  * [Tutorials Documentation](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/%5BTutorials%5D-Development-Documentation)
* [Interface Snippets](https://github.com/SICKAG/sick_tdc-e-developers-documentation/tree/main/interface-snippets)
* [Examples](https://github.com/SICKAG/sick_tdc-e-developers-documentation/tree/main/examples)
* [Tutorials](https://github.com/SICKAG/sick_tdc-e-developers-documentation/tree/main/tutorials)
* [Useful Links](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Useful-Links)
* [Docs](https://github.com/SICKAG/sick_tdc-e-developers-documentation/tree/main/docs)

---

## 📂 Repository Structure

This repository contains:

- **[docs/tdc-e-base/](docs/tdc-e-base/)** - BaseOS documentation source (published to GitHub Pages)
- **[interface-snippets/](interface-snippets/)** - Code examples by interface type
- **[examples/](examples/)** - Complete project examples
- **[tutorials/](tutorials/)** - Step-by-step tutorials
- **[Wiki](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki)** - Legacy L4M firmware documentation

---

## 🚀 Quick Links

- **BaseOS Documentation:** https://sickag.github.io/sick_tdc-e-developers-documentation/
- **L4M Wiki:** https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki
- **Code Examples:** [examples/](examples/) | [interface-snippets/](interface-snippets/)

---

## TDC-E Configuration
The TDC-E Configuration section describes how to configure parts of the TDC-E so as to be able to utilize it to its full extent. The configuration is split into two parts: network and interface configuration. Each one describes how to connect your TDC-E with the desired service.

## Getting Started
Getting Started is a section which mostly focuses on deploying applications and environment creation for the TDC-E device. It also focuses on installing and setting services via _Dockerfile_ and _docker_compose.yml_. 

## Development Documentation
The Development Documentation is the part of the documentation that describes code examples split into two main categories: one that focuses on applications that are implemented to use TDC-E's hardware interface, while the other part focuses on describing applications that do not use it to access needed data. To that end, this section focuses on the hows and whys of the programming experience. This section also contains documentation of more complex examples that combine multiple service to fulfill functionalities.

## Interface Snippets
The interface snippets, divided into sections pertaining the interface the service the application uses, are meant to be a help for starting application development with said service. Most examples are made with either Python or Go, with an additional Node-RED application as an alternative. 

## Examples
The examples contain more complex project examples that showcase various interfaces, technologies and services to fulfill their functionalities.

## Tutorials
Tutorial examples are codes and environments written specifially for users to be able to download and set up their environment with no further actions needed.

## Useful Links
Useful links contain relevant related articles regarding the TDC-E usage and configuration.

## Docs
Useful documents for additional TDC-E support.

## Notes

* Newer versions of L4M might require minor changes in application logic. Applications refactored for recent releases of the TDC-E device firmware are named in `appname_version_X_X` format, where X is the version of the firmware.
