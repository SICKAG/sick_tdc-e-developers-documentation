FROM arm32v7/debian AS build

# Run installs
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    apt-utils \
    nano \
    python3 \
    python3-pip \
    iproute2 \
    python3-dev \
    default-libmysqlclient-dev \
    libssl-dev \
    pkg-config \
    build-essential

# pip install needed libraries with cert
RUN pip3 install requests --break-system-packages

# Create a directory for the app
WORKDIR /app

# Copy the necessary Python files
COPY readDIO.py ./
COPY toggleState.py ./
COPY direct.py ./

# Set the entry point command
CMD ["python3", "readDIO.py"]
