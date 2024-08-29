# Simple Host Documentation
This simple Go program creates an HTTP server that listens on a specified port and responds with a message indicating the port number. It's useful for testing purposes, such as verifying the behavior of a reverse proxy or checking network configurations.

# How It Works
- Port Specification: The port on which the server listens is specified via a command-line flag (-port). If no port is specified, the server defaults to port 8080.
- HTTP Response: When accessed, the server responds with a message that includes the port number on which it is running.

# Code Explanation
- flag.Int: Used to parse the port number from the command line. It defaults to 8080 if not provided.
- http.NewServeMux: Sets up a new HTTP request multiplexer (router).
- router.HandleFunc: Registers a simple handler that responds to all incoming requests with a message including the port number.
- http.Server: Configures and starts the HTTP server on the specified port.

# Usage
1. Compile and Run:
- Compile and run the program, specifying the desired port number:
```bash
go run simple_host.go -port=8001
```
If you don't specify a port, the server will default to port 8080:

```bash
go run simple_host.go
```
2. Test the Server:

- Open your browser or use curl to make a request to the server:

```bash
curl http://localhost:8001/
```
The server will respond with:

```text
Hello from port 8001
```
3. Shutdown:
- To stop the server, press Ctrl+C in the terminal where the server is running.

Conclusion
This simple HTTP server is a quick and easy tool for testing. It can be used to verify reverse proxies, load balancers, or any network configuration that requires a lightweight HTTP service.