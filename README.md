<h1 align="center" style="border-bottom: none">
    <img alt="logo" src="./docs/govpp-logo.png"><br>RapidVPP
</h1>

RapidVPP controller for managing FD.io VPP using Go and govpp. Provides a modular REST API for interfaces, IP configuration, and ACLs via VPP binapi.

## Useful Links

- **GoVPP Documentation**: [GitHub Wiki](https://github.com/FDio/govpp/wiki)
- **VPP Documentation**: [Official Docs](https://s3-docs.fd.io/vpp/26.02/)
- **GoVPP Repository**: [FDio/govpp](https://github.com/FDio/govpp)

## Features & Modules

RapidVPP provides four core modules for managing FD.io VPP through a REST API:

### 1. VPP Information Module
**Purpose**: Retrieve basic VPP system information  
**Key Function**: Get VPP version and system details  
**Use Case**: Health checks, compatibility verification, system monitoring

### 2. Interface Management Module
**Purpose**: Create, configure, and manage network interfaces  
**Key Functions**:
- Create/delete loopback interfaces
- Set interface administrative states (up/down)
- Configure IP addresses on interfaces
- Attach/detach ACLs to interfaces  

**Use Case**: Network interface provisioning, interface state management

### 3. IP Configuration Module
**Purpose**: Manage routing tables and VRF instances  
**Key Functions**:
- Add/delete routes
- Create/delete VRF tables
- List routes by VRF
- VRF cache initialization

**Use Case**: Dynamic routing, multi-tenant network isolation, route management

### 4. ACL Management Module
**Purpose**: Create and manage Access Control Lists  
**Key Functions**:
- Create/update/delete ACL rules
- List configured ACLs
- ACL rule management (permit/deny)

**Use Case**: Network security, traffic filtering, policy enforcement

## API Reference

### VPP Information
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/vpp/version` | Get VPP version information |

### Interface Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/interfaces/` | List all interfaces |
| `POST` | `/interfaces/loopback` | Create a new loopback interface |
| `POST` | `/interfaces/{id}/state` | Set interface state (up/down) |
| `POST` | `/interfaces/{id}/ip` | Add IP address to interface |
| `DELETE` | `/interfaces/{id}` | Delete loopback interface |
| `GET` | `/interfaces/{id}/acl` | List ACLs attached to interface |
| `POST` | `/interfaces/{id}/acl` | Attach ACL to interface |
| `DELETE` | `/interfaces/{id}/acl` | Detach ACL from interface |


### IP Configuration & Routing
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/routes` | List all routes |
| `GET` | `/routes/{vrf}` | Get routes for specific VRF |
| `POST` | `/routes` | Add a new route |
| `DELETE` | `/routes` | Delete a route |
| `GET` | `/vrf` | List all VRF tables |
| `POST` | `/vrf` | Create VRF table |
| `DELETE` | `/vrf/{id}` | Delete VRF table |

### ACL Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/acl` | List all ACLs |
| `POST` | `/acl` | Create new ACL |
| `PUT` | `/acl/{id}` | Update existing ACL |
| `DELETE` | `/acl/{id}` | Delete ACL |


## 🚀 Build & Run

This section describes how to build and run the **RapidVPP** controller.
It requires a working **FD.io VPP** environment and **Go** tooling.

---

## Prerequisites

### System Dependencies (Ubuntu/Debian)

* Install **VPP** from the official **FD.io** repository
* Go (recommended version according to the project)

---

## Build the Project

You can build the project using the provided **Makefile** or directly with `go build`.

### Using the Makefile (Recommended)

```bash
# Build the server binary into the ./bin/ directory
make build
```

### Manually with Go

```bash
go build -o bin/rapidvpp-server ./cmd/app
```

---

## Run the Server

Ensure the **VPP** service is running before starting the controller:

```bash
sudo systemctl start vpp
```

### Using the Makefile

```bash
# This command will build and immediately run the project
make run
```

### Manually with Go

```bash
go run ./cmd/app/main.go
```

---

##  API Access

Upon successful startup, the REST API server will be available at:

```
http://localhost:8080
```

---

## Quick Command Reference (Makefile)

* `make build` — Build the binary
* `make run` — Build and run the server (development mode)
