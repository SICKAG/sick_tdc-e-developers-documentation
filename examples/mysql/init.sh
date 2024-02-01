#!/bin/bash

sleep 30;

mysql -e "CREATE DATABASE IF NOT EXISTS gpslocation;";
mysql -e "CREATE DATABASE IF NOT EXISTS diobase;"
mysql -e "CREATE DATABASE IF NOT EXISTS analog_base;"
mysql -e "CREATE TABLE IF NOT EXISTS gpslocation.gpsdata (id integer PRIMARY KEY AUTO_INCREMENT, latitude double, longitude double, time datetime, altitude double, speedKnots double, speedMph double, speedKmh double, course integer, fix integer, numberOfSatellites integer, gpsFixAvailable boolean, hdop double);"
mysql -e "CREATE TABLE IF NOT EXISTS diobase.dios (id integer PRIMARY KEY AUTO_INCREMENT, duration varchar(20), whattime timestamp);"
mysql -e "CREATE TABLE IF NOT EXISTS analog_base.analogi (id integer PRIMARY KEY AUTO_INCREMENT, whattime timestamp, value float);"
