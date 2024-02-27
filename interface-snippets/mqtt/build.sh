#!/bin/bash

docker build --no-cache -t mosquitto-confed .
docker save -o mqttconfed.tar mosquitto-confed

sleep 30