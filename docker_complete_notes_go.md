# Docker — Complete Notes & Practice (with Go Applications)

---

## 1. Introduction to Docker

**Docker** is a platform for building, shipping, and running applications inside lightweight, isolated environments called **containers**.

### Why Docker?
- **"Works on my machine" problem** — Docker packages your app + all its dependencies (runtime, libraries, config) into one portable unit.
- **Lightweight** — containers share the host OS kernel, unlike VMs which virtualize an entire OS. This makes them start in milliseconds and use far less memory/disk.
- **Consistency** — the same image runs identically on your laptop, CI server, and production.
- **Isolation** — each container has its own filesystem, process space, and network interface, isolated from others.

### Containers vs Virtual Machines
| | Container | Virtual Machine |
|---|---|---|
| Virtualizes | OS-level (shares host kernel) | Hardware-level (full OS per VM) |
| Size | MBs | GBs |
| Startup time | Milliseconds–seconds | Minutes |
| Isolation | Process-level (namespaces, cgroups) | Full OS-level |
| Density | Hundreds per host | Few per host |

### Key Docker Concepts (terms you'll see everywhere)
- **Image**: A read-only template/blueprint (your app + dependencies + OS layers) used to create containers.
- **Container**: A running (or stopped) instance of an image — the actual live process.
- **Dockerfile**: A text file with instructions to build an image.
- **Registry**: A place to store/share images (Docker Hub is the default public one).
- **Volume**: Persistent storage that lives outside the container's writable layer.
- **Network**: Virtual networking that lets containers communicate with each other and the outside world.

---

## 2. Docker Architecture and Components

Docker uses a **client-server architecture**:

```
 ┌─────────────┐         REST API / socket        ┌────────────────────┐
 │ Docker CLI   │ ───────────────────────────────► │   Docker Daemon     │
 │ (docker ...) │ ◄─────────────────────────────── │   (dockerd)         │
 └─────────────┘                                   └─────────┬──────────┘
                                                               │ manages
                                    ┌──────────────┬───────────┼───────────────┐
                                    ▼              ▼           ▼               ▼
                               Images        Containers    Networks        Volumes
```

### Core Components

1. **Docker Client (`docker`)** — the CLI tool you type commands into. Sends commands to the daemon via REST API (usually over a Unix socket `/var/run/docker.sock`).

2. **Docker Daemon (`dockerd`)** — the background service that actually builds images, runs containers, manages networks/volumes. Does the real work.

3. **containerd** — a lower-level container runtime (a CNCF project) that Docker delegates actual container lifecycle management to (start/stop/pause containers). Sits between `dockerd` and the OS kernel.

4. **runc** — the low-level tool that actually creates containers using Linux kernel primitives (namespaces + cgroups). `containerd` calls `runc`.

5. **Docker Registry** — stores and distributes images.
   - **Docker Hub**: default public registry (`docker.io`).
   - Private registries: AWS ECR, GCP Artifact Registry, GitHub Container Registry (ghcr.io), self-hosted (`registry:2` image).

6. **Docker Objects**:
   - **Images** — layered, read-only filesystems.
   - **Containers** — writable layer on top of an image + running process.
   - **Networks** — bridge, host, overlay, none.
   - **Volumes** — persistent data storage.

### What makes containers "lightweight" (Linux primitives used under the hood)
- **Namespaces** — isolate what a process can *see* (PID namespace, network namespace, mount namespace, etc.) — each container thinks it's the only thing running.
- **Cgroups (control groups)** — limit what a process can *use* (CPU, memory, disk I/O).
- **Union filesystems (e.g., OverlayFS)** — let images be built in layers, and containers add a thin writable layer on top without duplicating the whole filesystem.

---

## 3. Installing Docker Environment

### Linux (Ubuntu/Debian)
```bash
# Remove old versions
sudo apt-get remove docker docker-engine docker.io containerd runc

# Set up Docker's official repo
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Run docker without sudo
sudo usermod -aG docker $USER
newgrp docker
```

### macOS / Windows
- Install **Docker Desktop** (bundles daemon, CLI, Compose, and a lightweight Linux VM since containers need a Linux kernel).

### Verify Installation
```bash
docker --version
docker compose version
docker run hello-world     # pulls & runs a test image; confirms daemon works
docker info                # shows daemon, storage driver, containers, images summary
```

### Useful setup checks
```bash
docker context ls          # see which docker context (local/remote) is active
docker system df           # disk usage by images/containers/volumes
```

---

## 4. Creating Dockerfiles for Go Applications

### 4.1 A Simple Go App
```go
// main.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Dockerized Go app!")
	})
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
```
```
// go.mod
module myapp

go 1.22
```

### 4.2 Basic Dockerfile (single-stage — works but not optimal)
```dockerfile
FROM golang:1.22

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o server .

EXPOSE 8080
CMD ["./server"]
```
**Problem:** the final image includes the entire Go toolchain (~800MB+), even though you only need the compiled binary to run.

### 4.3 Multi-Stage Build (recommended — small, production-ready image)
```dockerfile
# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependency downloads separately from source code changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 for a fully static binary (needed for scratch/alpine base)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ---- Final stage ----
FROM alpine:3.19
# alpine needs ca-certificates if your app makes HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
```
Result: final image is often **~15-20MB** instead of ~800MB.

### 4.4 Even smaller: `scratch` base (no shell, no OS at all)
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM scratch
COPY --from=builder /app/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/server"]
```
- `scratch` is an empty image — smallest possible (just your binary). No shell means you can't `docker exec` into it easily for debugging.

### 4.5 Dockerfile Instruction Reference
| Instruction | Purpose |
|---|---|
| `FROM` | Base image to build from |
| `WORKDIR` | Sets working directory for subsequent instructions |
| `COPY` | Copies files from host into image |
| `ADD` | Like COPY, but also supports URL fetch & tar auto-extraction (prefer COPY unless you need this) |
| `RUN` | Executes a command at build time (creates a new layer) |
| `ENV` | Sets environment variable |
| `ARG` | Build-time variable (not present in final container) |
| `EXPOSE` | Documents which port the container listens on (doesn't actually publish it) |
| `CMD` | Default command to run when container starts (overridable) |
| `ENTRYPOINT` | Fixed main command (args from CMD/docker run get appended) |
| `USER` | Run as non-root user (security best practice) |
| `HEALTHCHECK` | Defines how Docker checks if container is healthy |

### 4.6 .dockerignore (important — reduces build context size)
```
.git
*.md
Dockerfile
.dockerignore
tmp/
*.log
node_modules
```

### 4.7 Best Practices for Go Dockerfiles
1. Use multi-stage builds — always.
2. Copy `go.mod`/`go.sum` **before** the rest of the source, so `go mod download` is cached and only re-runs when dependencies change.
3. Use `CGO_ENABLED=0` for static binaries compatible with `alpine`/`scratch`.
4. Run as a non-root user:
   ```dockerfile
   RUN adduser -D -g '' appuser
   USER appuser
   ```
5. Pin base image versions (`golang:1.22-alpine`, not `golang:latest`) for reproducible builds.
6. Add a `HEALTHCHECK`:
   ```dockerfile
   HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/health || exit 1
   ```

---

## 5. Building and Running Containers

### 5.1 Build an image
```bash
docker build -t myapp:1.0 .
docker build -t myapp:1.0 -f Dockerfile.prod .    # custom Dockerfile name
docker build --no-cache -t myapp:1.0 .            # ignore layer cache
docker build --build-arg VERSION=1.2 -t myapp .   # pass ARG values
```
- `-t` tags the image (`name:tag`).
- `.` is the **build context** (directory sent to the daemon).

### 5.2 Run a container
```bash
docker run myapp:1.0                     # run in foreground
docker run -d myapp:1.0                  # detached (background)
docker run -d -p 8080:8080 myapp:1.0     # map host:container port
docker run -d --name myserver myapp:1.0  # give it a name
docker run -it myapp:1.0 sh              # interactive shell
docker run --rm myapp:1.0                # auto-remove container when it stops
docker run -e ENV=production myapp:1.0   # pass environment variable
docker run --env-file .env myapp:1.0     # load env vars from file
```

### 5.3 Inspect running containers
```bash
docker ps                  # running containers
docker ps -a                # all containers (incl. stopped)
docker logs <container>     # view stdout/stderr
docker logs -f <container>  # follow logs live
docker inspect <container>  # full JSON metadata
docker stats                # live CPU/memory usage
docker top <container>      # processes inside container
```

### 5.4 Execute commands inside a running container
```bash
docker exec -it <container> sh          # open a shell (or bash if available)
docker exec <container> ls /app
```

### 5.5 Stop / start / remove
```bash
docker stop <container>
docker start <container>
docker restart <container>
docker rm <container>
docker rm -f <container>       # force remove (even if running)
docker kill <container>        # SIGKILL immediately (vs graceful SIGTERM from stop)
```

---

## 6. Docker Images and Container Lifecycle

### 6.1 Image commands
```bash
docker images                      # list local images
docker pull golang:1.22-alpine     # download an image
docker push myrepo/myapp:1.0       # upload to a registry
docker rmi myapp:1.0               # remove an image
docker tag myapp:1.0 myrepo/myapp:1.0  # add a tag/alias
docker history myapp:1.0           # show image layers & sizes
docker image prune                 # remove dangling (untagged) images
docker image prune -a              # remove ALL unused images
```

### 6.2 Image Layers
Every Dockerfile instruction (`RUN`, `COPY`, `ADD`) creates a new **layer**. Layers are cached and reused across builds if nothing above them changed — this is why instruction **order matters** (put things that change least often at the top).

```
Layer 4: COPY . .              <- changes often (source code)
Layer 3: RUN go mod download   <- changes rarely (deps)
Layer 2: COPY go.mod go.sum .  <- changes rarely
Layer 1: FROM golang:1.22-alpine  <- base, changes almost never
```

### 6.3 Container Lifecycle (states)

```
        docker create
             │
             ▼
         [Created] ──docker start──► [Running] ──docker stop──► [Stopped/Exited]
                                          │                              │
                                     docker pause                   docker start
                                          ▼                              │
                                      [Paused] ──docker unpause──► [Running]
                                                                          
         [Running/Stopped] ──docker rm──► [Removed]
```

- `docker run` = `docker create` + `docker start` combined.
- **Created** — container exists (filesystem allocated) but process hasn't started.
- **Running** — process is executing.
- **Paused** — process frozen (via cgroups freezer), no CPU scheduling, but state kept in memory.
- **Stopped/Exited** — process has ended (either normally or via `stop`/crash); filesystem still exists so you can inspect logs or restart it.
- **Removed** — container and its writable layer are deleted permanently.

### 6.4 Container restart policies
```bash
docker run --restart=no myapp          # default: never restart automatically
docker run --restart=on-failure myapp  # restart only if exit code != 0
docker run --restart=always myapp      # always restart, even after daemon restart
docker run --restart=unless-stopped myapp  # like always, but not if manually stopped
```

---

## 7. Docker Volumes and Networking

### 7.1 Why volumes?
Containers are **ephemeral** — data written inside a container's writable layer is lost when the container is removed. **Volumes** provide persistent storage that survives container removal and can be shared between containers.

### 7.2 Volume types
| Type | Managed by | Use case |
|---|---|---|
| **Named volume** | Docker | Databases, persistent app data — Docker manages the storage location |
| **Bind mount** | You (host path) | Local development — mount your source code directly into the container |
| **tmpfs mount** | Memory only | Sensitive/temp data that shouldn't touch disk |

```bash
# Named volume
docker volume create mydata
docker run -d -v mydata:/var/lib/data myapp

# Bind mount (absolute host path required)
docker run -d -v /home/user/app:/app myapp
docker run -d -v $(pwd):/app myapp        # common in local dev

# tmpfs (Linux only)
docker run -d --tmpfs /app/cache myapp
```

### 7.3 Volume management commands
```bash
docker volume ls
docker volume inspect mydata
docker volume rm mydata
docker volume prune          # remove all unused volumes
```

### 7.4 Example: Go app + Postgres with persistent volume
```bash
docker volume create pgdata

docker run -d \
  --name pg \
  -e POSTGRES_PASSWORD=secret \
  -v pgdata:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:16
```
Even if you `docker rm pg`, the data in `pgdata` survives and can be reattached to a new container.

### 7.5 Docker Networking — network drivers
| Driver | Description |
|---|---|
| `bridge` (default) | Private internal network on the host; containers can talk to each other via container name (DNS) if on the same custom bridge network |
| `host` | Container shares the host's network stack directly (no isolation, best performance, Linux only) |
| `none` | No networking at all |
| `overlay` | Multi-host networking for Docker Swarm/clusters |
| `macvlan` | Assigns a MAC address so container appears as a physical device on the network |

### 7.6 Networking commands
```bash
docker network ls
docker network create mynet
docker network inspect mynet
docker network connect mynet <container>
docker network rm mynet
```

### 7.7 Why you need a **custom** bridge network (not the default one)
On the **default bridge network**, containers can only reach each other by IP, not by name. On a **custom bridge network**, Docker provides automatic DNS resolution by container name.

```bash
docker network create app-net

docker run -d --name db --network app-net postgres:16
docker run -d --name api --network app-net -e DB_HOST=db myapp
# inside "api" container, connecting to host "db" just works (DNS resolves to db's IP)
```

### 7.8 Port publishing vs container-to-container communication
- `-p 8080:80` (host:container) is only needed to expose a port to **outside** the Docker network (your laptop's browser, external clients).
- Containers on the same custom network can talk to each other on their internal ports **without** publishing anything.

---

## 8. Practice: Complete Go + Docker Example

### Project structure
```
myapp/
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── .dockerignore
└── docker-compose.yml
```

### `main.go`
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Dockerized Go app!")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Printf("Server listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

### `Dockerfile` (production-grade multi-stage)
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates && \
    adduser -D -g '' appuser
WORKDIR /home/appuser
COPY --from=builder /app/server .
USER appuser

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/health || exit 1
ENTRYPOINT ["./server"]
```

### `docker-compose.yml` (app + database, custom network, volume)
```yaml
version: "3.9"

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DB_HOST=db
    depends_on:
      - db
    networks:
      - app-net

  db:
    image: postgres:16
    environment:
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=myapp
    volumes:
      - pgdata:/var/lib/postgresql/data
    networks:
      - app-net

networks:
  app-net:
    driver: bridge

volumes:
  pgdata:
```

### Run it
```bash
# Build & run manually
docker build -t myapp:1.0 .
docker run -d -p 8080:8080 --name myapp-container myapp:1.0
curl http://localhost:8080

# Or with Compose (builds + wires networking/volumes automatically)
docker compose up -d --build
docker compose logs -f api
docker compose down            # stop and remove containers/network
docker compose down -v         # also remove volumes
```

---

## 9. Common Docker Commands Cheat-Sheet

| Task | Command |
|---|---|
| Build image | `docker build -t name:tag .` |
| Run container | `docker run -d -p 8080:8080 name:tag` |
| List running containers | `docker ps` |
| List all containers | `docker ps -a` |
| View logs | `docker logs -f <container>` |
| Shell into container | `docker exec -it <container> sh` |
| Stop container | `docker stop <container>` |
| Remove container | `docker rm <container>` |
| List images | `docker images` |
| Remove image | `docker rmi <image>` |
| Create volume | `docker volume create <name>` |
| Create network | `docker network create <name>` |
| Clean up everything unused | `docker system prune -a --volumes` |
| Compose up | `docker compose up -d --build` |
| Compose down | `docker compose down -v` |

---

## 10. Practice Questions & Answers

**Q1. What's the difference between an image and a container?**
**A:** An image is a static, read-only template (layers of filesystem + metadata). A container is a running instance of that image, with its own writable layer added on top.

**Q2. Why use multi-stage builds for Go apps?**
**A:** The Go compiler and build toolchain aren't needed at runtime — only the compiled binary is. Multi-stage builds compile in a `golang` image, then copy just the binary into a minimal final image (`alpine`/`scratch`), drastically reducing image size and attack surface.

**Q3. What happens to data written inside a container if you don't use a volume?**
**A:** It's stored in the container's writable layer and is lost permanently when the container is removed (`docker rm`).

**Q4. Why does container-name DNS resolution fail on the default bridge network but work on a custom one?**
**A:** Docker only provides automatic embedded DNS for containers on user-defined (custom) networks. The default bridge network requires manually linking containers or using IP addresses.

**Q5. What's the difference between `CMD` and `ENTRYPOINT`?**
**A:** `ENTRYPOINT` sets the fixed executable that always runs; `CMD` supplies default arguments (which can be overridden by `docker run image <args>`). If both are set, `CMD`'s values become arguments to `ENTRYPOINT`.

**Q6. How do you keep Docker layer caching effective in a Go Dockerfile?**
**A:** Copy `go.mod`/`go.sum` and run `go mod download` **before** copying the rest of the source code. Since dependencies change far less often than source, this layer stays cached across most builds, speeding up rebuilds significantly.

**Q7. What Linux kernel features make containers possible?**
**A:** Namespaces (isolate what a process can see: PID, network, mount, etc.) and cgroups (limit what resources a process can use: CPU, memory). Union filesystems like OverlayFS enable the layered image model.

---

*End of notes. Try building the example Go app above from scratch, running it with both `docker run` and `docker compose`, then intentionally break the network/volume config to see the failure modes — that's the fastest way to really understand this.*
