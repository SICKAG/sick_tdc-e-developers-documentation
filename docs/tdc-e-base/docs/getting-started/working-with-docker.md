---
title: Working With Docker
sidebar_position: 102
---

import Admonition from '@theme/Admonition';


# Working With Docker

This page offers insights into using Docker services on this device, including how to work with Portainer and how to set up access to host services from within containers.

Docker is a platform that simplifies building, sharing, running, and packaging applications and their dependencies into lightweight, portable containers. These containers can run consistently across different environments.

Learn more about Docker [here](https://www.docker.com/).

This device supports Docker for both development and deployment purposes.

> **ℹ️ Info**
>
> Make sure to refer to [Docker Authorization Policy Scope](/getting-started/docker-authorization-policy-scope) when working with Docker. Created containers also have to adhere to the policy.


## 1. Creating Docker Containers

In this section, creating Docker containers is discussed. This includes information about providing access to host devices, mounting volumes, exposing ports, and specifying capabilities.

A Docker container is a **lightweight, isolated runtime environment** that packages an application with everything it needs to run - code, libraries, system tools, configuration, and environment settings. Because all dependencies are included, the application **runs the same way on any system** that supports Docker.

Containers are created from **Docker images**, which typically start with a small **base image** (e.g., Alpine, Debian, or an application‑specific base like Nginx). Additional layers add the application and its dependencies. When the Docker Engine runs the image, it creates a container as a set of isolated processes on the host system, using a layered filesystem and sharing the host’s kernel.

Docker containers can be created using either the **Docker CLI**, **Docker compose**, or **container management tools** (e.g. Portainer). When creating a container, the following needs to be decided:

- Which image the container should run

- How the container should start and behave

- Which host resources (network, storage, devices) it can access

### 1.1. Basic Container Creation

The simplest way to start a container is using the `docker run` command:

```bash
docker run IMAGE
```

For example, the following CLI command will run the [alpine](https://hub.docker.com/_/alpine) Docker image.

```bash
docker run -it alpine sh
```

This will:

- Pull the (latest) Alpine image from online registry if it is not present locally

- Create a container from the image

- Start the container with default settings

- Open an interactive shell inside the container

> **ℹ️ Info**
>
> The default online registry for pulling images is [Docker Hub](https://hub.docker.com/). Pulling images requires an **internet connection**.


### 1.2. Commonly Used Options

In practice, containers are usually started with additional options to configure their runtime environment. See some examples below.

**Run in background and give the container a name:**

```bash
docker run -d --name my-alpine alpine sh
```

This runs the `alpine` image in the background `(-d option, short for daemon)`, and assigns `my-alpine` as the name of the container.

**Mount a host directory into the container:**

By default, a Docker container has its own isolated filesystem and cannot access files on the host. However, applications often need to read host‑provided files or persist data outside the container. Docker provides two options:

- Volumes – storage handled by Docker
- Bind mounts – direct host‑to‑container directory mappings

Both allow shared access to directories between host and container. This ensures that data created inside the container persists after shutdown, as it is stored on the host rather than in the container’s temporary filesystem.

```bash
docker run -it \
  -v /datafs/operator:/home/root/data \
  alpine sh
```

> **📝 Note**
>
> See [Device-specific Paths](/getting-started/working-with-docker#13-device-specific-paths) below for mapping correct paths to your container.


This starts the alpine container in an interactive shell with data from **host** location `/datafs/operator` mounted to **container** location `/home/root/data`. Files created in the container directory `/home/root/data` will persist on the host.

**Expose a network port:**

```bash
docker run -d \
  --name my-nginx \
  -p 8080:80 \
  nginx
```

This starts an Nginx web server and exposes it on the host at port 8080, forwarding traffic to port 80 inside the container.

**Grant access to a host device or additional capabilities:**

```bash
docker run -it \
  --device /dev/serial/by-id/tdcx-serial:/dev/ttyUSB0 \
  --cap-add NET_ADMIN \
  alpine sh
```

This runs the Alpine shell with:

- the `NET_ADMIN` capability (allowing network configuration)
- the **host** device `/dev/serial/by-id/tdcx-serial` mapped to the **container** as `/dev/ttyUSB0`.

> **📝 Note**
>
> Access to host devices, privileged mode, and additional capabilities may be restricted by the device’s Docker authorization policy.
>
> See [Device-specific Paths](/getting-started/working-with-docker#13-device-specific-paths) below for mapping correct paths to your container.


**Docker Compose**

The same options can be applied to containers using **Docker Compose** and writing a `docker-compose.yml` file. Compose provides a declarative and reusable way to define container setup.

For more information go to [Docker Compose Docs](https://docs.docker.com/compose/).

### 1.3. Device-specific Paths

Some host paths and devices (e.g. SERIAL devices) are exposed at **device-specific locations**.

When mapping volumes or granting device access to a container, always use the paths that correspond to your device.

See a list of **host devices and mount points** for adding to your container below:


## 2. Host Service Setup

When running Docker containers, you may need to access services running on the host machine (e.g. a containerized application accessing databases or APIs running on the host). Docker provides a special **DNS name** for this purpose inside containers:

```json
host.docker.internal
```

This hostname resolves to the internal IP address of the host machine **from inside the container**, thus allowing containers to communicate with host services without hardcoding IP addresses.

### 2.1. Example Usage

A `host.docker.internal` usage is demonstrated by starting a `Python3` HTTP server, and running an `alpine` container on the TDC-E device.

A simple HTML page titled `hello.html` is created. In the same directory, a Python server is started:

```bash
python3 -m http.server <port>
```

This should yield the following output:

```bash
python3 -m http.server 8000
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...
```

Start an alpine container with the `host.docker.internal` host, and map the gateway dynamically using `host-gateway`. Add the `curl` command, then curl the host server for the HTML file.

```bash
docker run -it --rm --add-host=host.docker.internal:host-gateway alpine sh
apk add curl

curl http://host.docker.internal:8000/hello.html
```

The result is shown below.

```html
<!-- hello.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Hello</title>
</head>
<body>
  <h1>Hello, World! 🐳</h1>
  <p>This is served from a Python3 HTTP server.</p>
</body>
</html>
```

## 3. Working with Portainer

In this section, working with Portainer is discussed. The section provides an insight into accessing the Portainer service after application installation, deploying images, and building a container for an image.

Portainer is a lightweight web-based GUI for managing Docker. It lets you control containers, images, networks, and volumes without using Docker CLI commands. Portainer itself runs as a Docker container and is given access to the Docker Engine so it can manage other containers on the host.

### 3.1. Accessing Portainer


### 3.2. Deploying an Image on Portainer

To be able to deploy an application on the TDC-E device, the application needs an environment it can run in. This environment is called an `image`. 

To upload an image, an image file is needed. There are multiple ways of obtaining an image. In this example, we build an application using the `docker build` command, and the image is saved using `docker save`.

> **ℹ️ Info**
>
> Check the system architecture before building your application. If the target architecture differs from your development environment, make sure to cross-compile accordingly.


```bash
docker build -t img-tag /path/to/dockerfile
docker save -o output.tar img-tag:latest
```

Replace the `tag`, `path/to/dockerfile` and `output` parameters accordingly.

After building the image, go to the [Portainer Dashboard](https://192.168.0.100:9443/#!/2/docker/dashboard) and select the `Images` option, or find the [Images](https://192.168.0.100:9443/#!/2/docker/images) option on the left panel.

The `.tar` file can be uploaded by selecting the **Import** option.

![Importing an Image](/img/port-1-image-tar.png)

Select a valid `.tar` file and give the image a fitting name. In this example, a `dio-grpc-app.tar` file was selected, and the image was named `dio-grpc`.

![Uploading an Image](/img/portainer-2-upload-img.png)

Select **Upload** and wait for the image to be uploaded to your device. Once it is on the device, Portainer will show the image and Dockerfile details which specify the ID, size, creation date, build, environment, command and layers of the image. 

### 3.3. Creating a Container from an Image

An image is needed for a container to be run. In other words, an image cannot do anything on its own as it needs a container to run. To provide a container to an imported image, go to the [Containers](https://192.168.0.100:9443/#!/2/docker/containers) tab on Portainer, where all current containers are listed.

![Containers](/img/port-3-build-container.png)

To add a container, select the option **Add container**. Provide a name for you container, and the image uploaded to Portainer in the last step.

![Setting Up a Container](/img/port-4-set-con.png)

One can also set the following parameters:

- Registry
- Always pull the image
- Create a container webhook
- Publish all exposed ports to random host ports
- Port Mapping
- Access control
- Auto Remove

You can also set **Advanced container settings**, which allows the user to set additional environment parameters like:

- Commands and logging
- Volumes
- Network
- Env
- Labels
- Restart policy
- Runtime and resources
- Capabilities

Once all options are set, select **Deploy the container**. The container will now start running.
