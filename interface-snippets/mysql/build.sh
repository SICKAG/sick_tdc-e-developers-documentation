#!/bin/bash

docker build -t mysql-img .
docker save -o mysql.tar mysql-img

sleep 30
