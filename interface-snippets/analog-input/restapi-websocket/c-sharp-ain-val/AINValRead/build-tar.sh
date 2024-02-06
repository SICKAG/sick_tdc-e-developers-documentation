#!/bin/bash

docker build -t file-name .
docker save -o cain.tar file-name

sleep 30