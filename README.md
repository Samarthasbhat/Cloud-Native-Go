
# Cloud Native Go ☁️

This repository is part of my learning journey in building **Cloud-Native applications using Go**.  
It includes experiments, services, and examples inspired by the book *Cloud Native Go* by Matthew A. Titmus.


## 🚀 What’s Inside

- **concurrency-patterns/**
  - Examples like debounce, retry, throttle, and timeout written in Go.
- **cloud-native-service/**
  - Simple key-value store and microservices built with `net/http` and Go routines.
- **Dockerfile**
  - Containerization setup for Go services.
- **README.md**
  - Project overview and usage guide.



## 📦 Getting Started

### Prerequisites

- Go (v1.20+)
- Docker & Docker Desktop (for container builds)
- OpenSSL (optional — for generating TLS certs)

### Clone the Repository

```bash
git clone https://github.com/Samarthasbhat/Cloud-Native-Go.git
cd Cloud-Native-Go
````



## 🧠 Run a Service Locally

To run an example from **concurrency-patterns**:

```bash
cd concurrency-patterns/retry
go run main.go
```

To run the key-value store service:

```bash
cd cloud-native-service/key-value-store
go run main.go
```

If the service uses TLS, generate your own local certificate:

```bash
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout key.pem -out cert.pem -days 365 \
  -subj "/CN=localhost"
```

Then run:

```bash
go run main.go
```

Test with `curl`:

```bash
curl -k https://localhost:8080/v1/somekey
```



## 🐳 Docker Setup

### Option A — Auto-generate certs inside container

Add this to your Dockerfile:

```dockerfile
FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN apk add --no-cache openssl \
    && openssl req -x509 -newkey rsa:4096 -nodes \
         -keyout key.pem -out cert.pem -days 365 \
         -subj "/CN=localhost"

RUN go build -o server .

CMD ["./server"]
```

Build and run:

```bash
docker build -t go-kv-store .
docker run -p 8080:8080 go-kv-store
```

### Option B — Use local certs at runtime

Keep `cert.pem` and `key.pem` locally (ignored by Git):

```bash
docker build -t go-kv-store .
docker run -v $(pwd)/cert.pem:/app/cert.pem \
           -v $(pwd)/key.pem:/app/key.pem \
           -p 8080:8080 go-kv-store
```



## 🧭 Project Structure

```
Cloud-Native-Go/
├── concurrency-patterns/
│   ├── debounce/
│   ├── retry/
│   ├── throttle/
│   └── timeout/
├── cloud-native-service/
│   └── key-value-store/
├── Dockerfile
└── README.md
```



## 🤝 Contributing

This repo is primarily for learning, but contributions are welcome!

1. Fork the repo
2. Create a feature branch
3. Add or improve a service/pattern
4. Open a pull request


## 📚 References

* *Cloud Native Go* — Matthew A. Titmus
* Go standard library (`net/http`, `context`, `sync`)
* Docker and TLS basics
* Concurrency patterns in Go



## ⚠️ Notes

* Don’t commit your `cert.pem` or `key.pem` — add them to `.gitignore`.
* For self-signed certs, use `curl -k` to skip verification.
* Each new service can have its own Dockerfile and Go module.



## 👨‍💻 Maintainer

Maintained by [@Samarthasbhat](https://github.com/Samarthasbhat)
Learning Cloud-Native Go, one service at a time ☁️💻

````



You can now copy that whole block into your local `README.md` file and commit:

```bash
git add README.md
git commit -m "Update README with project overview and Docker setup"
git push origin main
````
