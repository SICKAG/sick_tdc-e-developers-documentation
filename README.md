# SICK TDC-E Developer's Documentation

## Welcome!

This is the [SICK TDC-E Developer's Documentation](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki), aimed at customers, users and developers alike to explore the TDC-E's options, capabilities and programmability. We want the user to get to know the TDC-E device so that its wide variety of options could be used to their full extent. To that end, the Developer's Documentation servers as a guide which contains the following file structure:

* TDC-E Configuration
  * [Network Configuration](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/TDC%E2%80%90E-Network-Configuration)
  * [Interface Container Configuration](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/TDC%E2%80%90E-Interface-Configuration)
* Getting Started
  * [Environment Setup](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup)
    * [Setting up Node-RED Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#1-setting-up-node-red)
    * [Setting up Python Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment#2-setting-up-python-application-environment)
    * [Setting up C# Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#3-setting-up-c-application-environment)
    * [Setting up Go Application Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#4-setting-up-go-application-environment)
    * [Setting up MySql Environment](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Environment-Setup#5-setting-up-mysql-environment)
  * [Image Building and Composing](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Getting-Started-%E2%80%90-Build-and-Compose)
* Development Documentation
  * [Examples Using the Hardware Interface](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Development-Documentation-%E2%80%90-Examples-Using-the-Hardware-Interface)
  * [Examples Without the Hardware Interface](https://github.com/SICKAG/sick_tdc-e-developers-documentation/wiki/Development-Documentation-%E2%80%90-Examples-Without-Using-the-Hardware-Interface)
* [Examples](https://github.com/SICKAG/sick_tdc-e-developers-documentation/tree/getting_started/examples)

## TDC-E Configuration
The TDC-E Configuration section describes how to configure parts of the TDC-E so as to be able to utilize it to its full extent. The configuration is split into two parts: network and interface configuration. Each one describes how to connect your TDC-E with the desired service.

## Getting Started
Getting Started is a section which mostly focuses on deploying applications and environment creation for the TDC-E device. It also focuses on installing and setting services via _Dockerfile_ and _docker_compose.yml_. 

## Development Documentation
The Development Documentation is the part of the documentation that describes code examples split into two main categories: one that focuses on applications that are implemented to use TDC-E's hardware interface, while the other part focuses on describing applications that do not use it to access needed data. To that end, this section focuses on the hows and whys of the programming experience.

## Examples
The examples, divided into sections pertaining the interface the service the application uses, are meant to be a help for starting application development with said service. Most examples are made with either Python or Go, with an additional Node-RED application as an alternative.


