# reverse-proxy

## Description
A dead simple reverse-proxy written with Go. Welcome to the Simple Reverse Proxy project! This project demonstrates a basic reverse proxy written in Go, showcasing its ability to work across different frameworks.

## Project Overview
This project includes:

- A simple reverse proxy written in Go using plain `net/http` (no other web framework).
- Two sample services located in the `services` folder:
    - **service1**: A service written in Go.
    - **service2**: A service written in Node.js.

The main file contains a simple web server with a single route that re-routes requests to the corresponding service based on the request's country. The country is determined from a `Country` header in the original request.

## What is a Reverse Proxy?
A reverse proxy is a server that sits between client devices and backend servers, forwarding client requests to the appropriate backend server. It can be used for load balancing, improving security, and caching content.

### Use Cases for a Reverse Proxy:
- **Load Balancing**: Distributing incoming traffic across multiple servers to ensure no single server becomes overwhelmed.
- **Security**: Hiding the identity and characteristics of backend servers.
- **Caching**: Storing copies of frequently accessed content to reduce load on backend servers.

## 🚀 Running the Services

### Service 1 (Go)
Navigate to the `services/service1` directory:
```bash
cd services/service1
```

Run the service:
```bash
go run main.go
```

### Service 2 (Node.js)
Navigate to the `services/service2` directory:
```bash
cd services/service2
```

Run the service:
```bash
node server.js
```

## 🌐 Running the Reverse Proxy
Navigate to the root directory:
```bash
cd reverse-proxy
```

Run the reverse proxy:
```bash
go run main.go
```

Now you can send a test request to the reverse proxy server:
```bash
curl -H "Country: mx" http://localhost:8080/resources
```

If you change the Country header to "ar", the request will be routed to the second service.

```bash
curl -H "Country: ar" http://localhost:8080/resources
```

## How It Works
The reverse proxy server re-routes incoming requests to the appropriate service based on the Country header in the request. For example:

    Requests with Country: MX are routed to service1.
    Requests with Country: AR are routed to service2.

This setup can be easily modified to implement load balancing or other routing logic.

Enjoy exploring the Simple Reverse Proxy project!

### New section added to the readme.