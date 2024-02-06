#!/bin/bash

docker build --no-cache -t myapp:1.0.0 .
docker save -o dioled.tar myapp:1.0.0
sleep 30
