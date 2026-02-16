---
title: Docker Authorization Policy Scope
sidebar_position: 101
---

import Admonition from '@theme/Admonition';


# Docker Authorization Policy Scope

In this section, the Docker Authz authorization plugin is introduced, and policy restrictions are listed.

## 1. Introduction

To secure the TDC-E device's Docker engine, an authorization plugin is developed according to a security and compliance policy. This policy **restricts** certain operations like elevating privileges, accessing sensitive devices or writing risky configurations. The plugin is implemented as an **OPA (Open Policy Agent)** policy for Docker.

The policy is implemented as a set of rules Docker commands need to adhere to, and blocks any commands that would violate the device's security.

If any of the commands match the deny rules of the plugin, the request is denied and the plugin returns an `error` message. Otherwise, the Docker API processes the command.

In the following sections, Docker policy restrictions are listed. 

> **📝 Note**
>
> When developing a custom Docker container and image, make sure to adhere to the policy.


## 2. Policy Segments

### 2.1. Privileged Containers

_Running containers in privileged mode is **denied**._

Example:
```json
"HostConfig": { "Privileged": false }   // Allowed
"HostConfig": { "Privileged": true }    // Denied
```

### 2.2. Bind Mounts

_Mounting sensitive paths is **denied**. Make sure to use only **allowed bind mounts**!_

**Allowed Bind Paths:**

- `/datafs/operator`
- `/var/run/iolink`
- `/var/run/cc`, `/run/cc`
- `/var/run/hal`, `/run/hal`
- `/var/run/docker.sock`
- `/var/volatile/tdcx`
- `/etc/tdcx`
- `/etc/os-release`
- `/media`
- `/datafs/tdc-engine`
- `/etc/machine-id`
- `/run/log/journal`
- `/usr/lib/os-release`
- `/dev/mmcblk1`, `/dev/mmcblk1p*`
- `/datafs/appruntime`

**Allowed Read-Only Bind Paths:**

- `/etc/tdcx`
- `/etc/os-release`
- `/etc/machine-id`
- `/run/log/journal`
- `/usr/lib/os-release`

> **📝 Note**
>
> All other bind mounts are **denied**.


Example:
```json
"BindMounts": [
    { "Source": "/etc/tdcx", "ReadOnly": true }    // Allowed
    { "Source": "/etc/tdcx", "ReadOnly": false }   // Denied
    { "Source": "/root", "ReadOnly": true }        // Denied
]
```

### 2.3. Volume Creation

_Volume creation is subjected to the same **bind mount restrictions as containers**._

**Denied volumes:**

- Bind mounts to **non-whitelisted** or **non-read-only** sensitive paths

### 2.4. Dangerous Capabilities

_Adding dangerous Linux capabilities is **denied**._

**Denied capabilities:**

- `all`
- `audit_control`
- `sys_admin`
- `sys_module`
- `sys_ptrace`
- `syslog`
- `dac_read_search`

Other capabilities are considered non-dangerous and **allowed**.

Example:
```json
"HostConfig": { "CapAdd": ["NET_ADMIN"] }      // Allowed
"HostConfig": { "CapAdd": ["SYS_ADMIN"] }      // Denied
```

### 2.5. Host Namespaces

_Using host namespaces is **denied**._

Therefore, the following are **denied**:

- `UsernsMode: "host"`
- `PidMode: "host"`
- `UTSMode: "host"`
- `CgroupnsMode: "host"`
- `IPCMode: "host"`

Example:
```json
"HostConfig": { "UsernsMode": "host" }   // Denied
"HostConfig": { "PidMode": "host" }      // Denied
```

### 2.6. Publishing All Ports

_Publishing all ports is **denied**._

Example:
```json
"HostConfig": { "PublishAllPorts": true }   // Denied
```

### 2.7. Security Options

_Unconfined Seccomp, AppArmor, or disabled SELinux labeling is **denied**._

**Denied options:**

- `seccomp=unconfined`
- `apparmor=unconfined`
- `label:disable`

Example:
```json
"HostConfig": { "SecurityOpt": ["seccomp=unconfined"] }   // Denied
```

### 2.8. Sensitive Sysctls

_Sensitive sysctls are **denied**._

The following are **denied**:
- `kernel.core_pattern`
- any sysctl starting with `net.`

Example:
```json
"HostConfig": { "Sysctls": { "kernel.core_pattern": "..." } }   // Denied
"HostConfig": { "Sysctls": { "net.ipv4.ip_forward": "1" } }     // Denied
```

### 2.9. Device Access

_Access to host devices is **restricted**._

Only the following devices are **allowed:**

- `/dev/serial/by-id/`
- `/dev/mmcblk1`, `/dev/mmcblk1p*`

> **📝 Note**
>
> All other devices are **denied**.
> `DeviceCgroupRules` and `DeviceRequests` usage is **denied**.


Examples:
```json
"HostConfig": {
    "Devices": [{ "PathOnHost": "/dev/serial/by-id/tdcx-serial" }]   // Allowed
    "Devices": [{ "PathOnHost": "/dev/sda" }]                        // Denied
    "DeviceCgroupRules": ["c 42:* rmw"]                              // Denied
    "DeviceRequests": [{ ... }]                                      // Denied
}
```

### 2.10. Copying Content into Containers

_Copying content from the host into containers using the archive API is **denied**._

**Denied PUT requests:**
- PUT requests to `/containers/{id}/archive`

If any deny rule listed above matches, the Docker command will fail.

## 3. FAQ

**Q | In some cases my docker container does not start. Looking at logs following error messages can be seen: 'Failed to allocate and map portaddress already in useError starting userland proxy'.**

A | Ensure that no other Docker containers are using the same port. 

If no conflicts are found, the issue may be related to a known Docker Engine issue. Restarting the device should resolve the problem and allow the container to start normally.

---

**Q | I've made changes to TDC-E SSH work environment and now I have troubles working with it.**

A | To reset SSH workspace issue the following operation: `container-reset-userworkspace` via REST API `/api/v1/system/administration/operation`.

Make sure to confirm the operation with `/api/v1/system/administration/operation/confirm`. Access the Control Center API UI via the page Resources → Documentation → Control Center.

---

**Q | I've made changes to AppEngine either by deploying apps or manipulating its variables and now I have troubles working with it.**

A | To reset AppEngine issue an operation `container-reset-appengine` via REST API `/api/v1/system/administration/operation`.

Make sure to confirm the operation with `/api/v1/system/administration/operation/confirm`. Access the Control Center API UI via the page Resources → Documentation → Control Center.

---

**Q | After updating to TDC-E FW 1.4.0 it's no longer possible to reach IO-Link REST API on host port 9000 from inside of container.**

A | From version 1.4.0 IO-Link API is available on localhost endpoint `127.0.0.1:19005`. If your container is using `network=host` then this API is reachable via `127.0.0.1:19005`.

If your container is using a `bridge network`, then add a parameter `--add-host` ("extra_host" in docker compose) `"host.docker.internal:host-gateway"`. 
This way, IO-Link API will be reachable on endpoint `host.docker.internal:19005`.

---

**Q | After updating to TDC-E FW 1.4.0 I cannot deploy my docker applications anymore. Error mentions 'authorization denied by plugin'.**

A | From version 1.4.0 Docker API authorization is activated which enforces a policy for every Docker API operation.
