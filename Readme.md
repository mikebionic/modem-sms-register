# Modem SMS Register

A robust Go CLI application that interfaces with GSM modems to receive SMS messages and forward them to HTTP endpoints for processing. This service acts as a bridge between physical GSM modems and web applications, enabling SMS-based integrations and workflows.

## 🎯 Project Aim

This project was created to solve the challenge of integrating SMS functionality into modern web applications without relying on expensive SMS gateway services. By using readily available GSM modems and SIM cards, organizations can:

* **Reduce SMS costs** - Use local SIM cards instead of premium SMS APIs
* **Gain independence** - No reliance on third-party SMS services
* **Ensure reliability** - Direct hardware control over SMS reception
* **Enable custom workflows** - Process SMS messages according to specific business logic

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────────┐    ┌─────────────────────┐
│   GSM Modem     │    │  modem-sms-register  │    │   Your Web API      │
│   (Hardware)    │───▶│     (This App)       │───▶│   (HTTP Endpoint)   │
│                 │    │                      │    │                     │
│ • Receives SMS  │    │ • Reads from modem   │    │ • Processes SMS     │
│ • USB/Serial    │    │ • Forwards to API    │    │ • Business logic    │
│ • SIM card      │    │ • Error handling     │    │ • Database storage  │
└─────────────────┘    └──────────────────────┘    └─────────────────────┘
```

---

## 🏃‍♂️ Running the Application

You can run the application with:

```bash
# Specify a custom config location
export CONFIG_PATH=/path/to/config.json

# Run the app directly
go run cmd/main.go

# Or build and run
make build
./bin/modem-sms-register
```

> By default, the application will look for `config.json` in the working directory if `CONFIG_PATH` is not set, don't forget to copy the `config.sample.json`.

---

## 📁 Project Structure

```
modem-sms-register/
├── cmd/
│   └── main.go            # Entry point
├── pkg/
│   ├── config/
│   │   └── config.go      # Reads configuration
│   ├── modem/
│   │   └── modem.go       # Modem and SMS receiver logic
│   └── sms/
│       └── request.go     # HTTP request handling (Send / makeRequest)
├── go.mod
├── go.sum
├── config.json             # Default config
├── config.example.json
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## 🚀 Use Cases

### Business Applications

* **Two-Factor Authentication** - Receive verification codes for secure login systems
* **Customer Support** - Process incoming support requests via SMS
* **Order Notifications** - Handle order confirmations and status updates
* **Emergency Alerts** - Receive critical system notifications

### IoT and Automation

* **Remote Monitoring** - Collect sensor data transmitted via SMS
* **Device Management** - Receive status updates from remote devices
* **Alert Systems** - Process emergency notifications from field equipment

### Personal Projects

* **Home Automation** - Control smart home devices via SMS commands
* **Weather Monitoring** - Receive weather station updates
* **Security Systems** - Get alerts from security cameras or sensors

---

## 📋 Requirements

### Hardware

* **GSM Modem** - USB or serial-connected GSM modem (e.g., Huawei E173, SIM800, etc.)
* **SIM Card** - Active SIM card with SMS capability
* **USB Port** - For USB modems, or serial port for RS232 modems

### Software

* **Go 1.21+** - For building from source
* **Linux/Windows/macOS** - Cross-platform support
* **Docker** (optional) - For containerized deployment

### Network

* **HTTP Endpoint** - Target API server to receive SMS data
* **Internet Connection** - For HTTP requests (cellular or WiFi)

---

## ⚙️ Docker & Makefile Notes

* Docker Compose automatically sets `CONFIG_PATH=/app/config.json`
* Build the app in Docker using:

```bash
docker build -t modem-sms-register .
docker-compose up -d
```

* Local development with Makefile:

```bash
make run          # Build & run
make run-verbose  # Build & run with debug logs
```
