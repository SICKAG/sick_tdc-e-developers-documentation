FROM jakezp/arm32v7-mqtt

COPY mosquitto.conf /mosquitto/config/mosquitto.conf

# Create the password file and hash the password
RUN echo "user1:password" > /tmp/password_file && \
    mosquitto_passwd -U /tmp/password_file && \
    cp /tmp/password_file /mosquitto/config/password_file && \
    cp /tmp/password_file /mosquitto/config/local_password_file && \
    rm /tmp/password_file

# Copy the password file back to the local filesystem
VOLUME /mosquitto/config

CMD ["mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]
