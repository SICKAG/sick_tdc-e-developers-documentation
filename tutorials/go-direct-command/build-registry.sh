#!/bin/bash

docker build --no-cache -t registry-name:1.0.0 .
docker push registry-name:1.0.0

sleep 30