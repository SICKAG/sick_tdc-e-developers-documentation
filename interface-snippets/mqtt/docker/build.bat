docker build --no-cache -t mosquitto-confed .
docker save -o mqttconfed.tar mosquitto-confed

timeout /t 30