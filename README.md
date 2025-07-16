# image-app
## Image Upload & Gallery App
Technical test for Software Engineer in DoIT

A Dockerized application for uploading JPEG images and viewing them in a browser gallery. It includes:
- Nginx as reverse proxy & static file server
- Go API for handling uploads and metadata
- PostgreSQL as image metadata repository
- Persistent volume for storing uploaded images

## Overview
### Nginx :
* Serves static files (upload.html, gallery.html, and uploaded images)
* Proxies /api/* to Go backend

### Go App :
* POST /upload to accept and validate JPEG uploads (<=10MB)
* Stores image metadata in PostgreSQL
* Writes image file to specified directory with `filename-uuid.jpeg` format to avoid collision 
* GET /images to list uploaded image metadata (filename, URL, size, timestamp)

### PostgreSQL :
* Stores metadata in project.images table
* Initializes schema on first launch from migrations/database.sql

### Network separation :
* api-network: Nginx and Go containers
* repository-network: Go and PostgreSQL containers

## Getting Started
### 1. Clone The Repo

```bash
git clone https://github.com/rahadianir/image-app.git
cd image-app
```
### 2. Build & Run Services
Make sure you have docker compose installed and run this command to build and run the services.

```bash
docker compose up --build -d
```
### 3. Browse the app

Upload page: http://localhost/upload.html

Gallery page: http://localhost/gallery.html

Nginx serves on port 80

Backend API on port 8080

Postgresql on port 5432

### 4. Clean up
```bash
docker compose down --volumes
```

## Repository Structure
```
.
├── README.md
├── app
│   ├── Dockerfile                  # Dockerfile to build api service
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── core
│   │   │   └── dependency.go
│   │   ├── image                   # Main package to handle upload process
│   │   │   ├── handler.go
│   │   │   ├── logic.go
│   │   │   ├── port.go
│   │   │   └── repository.go
│   │   ├── model                   
│   │   │   └── image.go
│   │   ├── pkg
│   │   │   ├── logger
│   │   │   │   └── logger.go
│   │   │   ├── pagination
│   │   │   │   └── pagination.go
│   │   │   └── xhttp
│   │   │       └── xhttp.go
│   │   └── server
│   │       └── http.go
│   ├── main.go
│   └── static                      # Static HTML page for Web UI
│       ├── gallery.html
│       └── upload.html
├── docker-compose.yml              # Docker compose as IaC to setup/provision all needed components
├── migrations                      # Migration directory to store database schema
│   └── database.sql
└── proxy                           # Proxy directory to store nginx (reverse proxy) configuration
    └── nginx.conf
```

