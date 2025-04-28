# Goflect
Goflect is a super lightweight and fast URL shortener written in Go.

It's built to be simple to run, easy to scale and powerful enough to handle serious traffic.

# Installation
## Prerequisites
You have two options to run Goflect:
1. Using Docker (recommended)
2. Running locally

## Clone the Repository
```sh
git clone https://github.com/ardytstrn/goflect.git
cd goflect
```

## Environment Variables
Configuration is managed through a `.env` file. You can find an example `.env.example` file to guide you.
```sh
cp .env.example .env
```
Modify it if needed.

# Running Goflect
## 1. Running with Docker (Recommended)
Build and start the application:
```sh
docker-compose up -d --build
```
This will:
- Build the Go application inside a Docker container.
- Start the server at `http://localhost:8000` (or your configured port).

To stop the server:
```sh
docker-compose down
```

## 2. Running Locally
If you prefer to run it directly with Go:
```sh
go build ./cmd/goflect
./goflect
```
Note that you should manually configure PostgreSQL and Redis servers.

# License
MIT License - free to use, free to modify.
