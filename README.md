# kafNodeX 🚀
**Real-Time Data Streaming API with Golang and Redpanda (Kafka)**

This project implements a robust and high-performance backend API in Go for real-time data streaming using Redpanda (Kafka) as the message broker, WebSockets for live updates, and Prometheus for performance monitoring. It includes secure API key middleware, global rate limiting, and full test coverage.

---

## 📦 Features

- Real-time streaming of payloads using Kafka topics
- WebSocket connections for pushing processed data
- Prometheus metrics integration (`/metrics` endpoint)
- Global rate limiting via middleware
- API key-based authentication
- Load testing with `wrk`
- Unit + integration tests using Go's testing framework

---

## ✅ Prerequisites

Ensure you have the following installed:

- [Homebrew](https://brew.sh/)  
- Go (v1.17 or higher)
- Redpanda (via `rpk`)
- Docker (for containerized Redpanda)
- Prometheus
- wrk (for load testing)
- Websocat (WebSocket CLI)

> ⚠️ **macOS users**: You must manually allow Prometheus under  
System Settings → Security & Privacy → Advanced → “Allow Prometheus”.

---

## 🛠 Installation

### 1. Install Go and Tools

```bash
brew update
brew install go
brew install vectorizedio/tap/redpanda
brew install --cask docker
brew install prometheus wrk websocat kafkacat
```

### 2. Start Docker & Redpanda

```bash
open /Applications/Docker.app
rpk container start -n 1
```

### 3. Create Kafka Topic

```bash
rpk topic create stream-topic
rpk topic list
```

---

## ⚙️ Setup Go Project

```bash
go get -u github.com/gorilla/mux
go get -u github.com/segmentio/kafka-go
go get -u github.com/gorilla/websocket
go mod tidy
```

---

## 🔐 Environment Variable

```bash
export API_KEY="my_secret_api_key_12345"
```

---

## 🚀 Run the Server

```bash
cd cmd/api
go run main.go middleware.go
```

Server runs at: [http://localhost:8080](http://localhost:8080)

---

## 🧪 Running Tests

```bash
go test ./tests -v
```

---

## 🔁 Rate Limiter Testing

Edit `internal/api/ratelimiter.go`:
```go
const globalRequestsPerSecond = 1
const globalBurstLimit = 2
```

**Terminal 1:**
```bash
go run main.go middleware.go
```

**Terminal 2:**
```bash
for i in {1..20}; do curl -X POST http://localhost:8080/stream/start -H "X-API-Key: my_secret_api_key_12345"; done
```

---

## 🔑 API Key Authentication Testing

**Valid Key:**
```bash
curl -X POST http://localhost:8080/stream/{stream_id}/send \
  -H "X-API-Key: my_secret_api_key_12345" \
  -H "Content-Type: application/json" \
  -d '{"data": "test_data"}'
```

**Invalid Key:**
```bash
curl -X POST http://localhost:8080/stream/{stream_id}/send \
  -H "X-API-Key: invalid_key" \
  -H "Content-Type: application/json" \
  -d '{"data": "test_data"}'
```

---

## 📊 Prometheus Setup

### 1. Create Prometheus Config

Open your system Prometheus config (or use the one in `monitoring/prometheus/prometheus.yml`):

```yaml
global:
  scrape_interval: 15s
scrape_configs:
  - job_name: 'golang_api'
    static_configs:
      - targets: ['localhost:8080']
```

### 2. Start Prometheus

```bash
brew services start prometheus
brew services list
```

Access Prometheus at [http://localhost:9090](http://localhost:9090)

---

## 🔥 Load Testing with wrk

### `/stream/start` endpoint

```bash
wrk -t12 -c1000 -d30s -s benchmark/start_stream.lua http://localhost:8080
```

### `/stream/{stream_id}/send` endpoint

Replace `{stream_id}` with a real ID:

```bash
wrk -t12 -c1000 -d30s -s benchmark/send_data.lua http://localhost:8080
```

---

## 📈 Monitoring Endpoint

You can view all collected metrics here:

```bash
http://localhost:8080/metrics
```

Metrics include:
- HTTP request counts and status codes
- Kafka messages produced/consumed
- Request duration histograms

---

## 🧾 Summary Commands

```bash
# Set API key
export API_KEY="my_secret_api_key_12345"

# Start Redpanda
rpk container start -n 1

# Create Kafka topic
rpk topic create stream-topic

# Run API server
go run ./cmd/api/main.go

# Run tests
go test ./tests -v

# Start Prometheus
brew services start prometheus

# Load testing
wrk -t12 -c1000 -d30s -s benchmark/start_stream.lua http://localhost:8080
```

---

## 📄 License

This project was developed by **Ayush Bhardwaj** 
---
