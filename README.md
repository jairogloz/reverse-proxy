![golang](https://camo.githubusercontent.com/29f331ff0b9cd5621d1233c541c575511c7ebb7cd6c09cb18c175c8bc729d14b/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f676f2d2532333030414444382e7376673f7374796c653d666f722d7468652d6261646765266c6f676f3d676f266c6f676f436f6c6f723d7768697465)
# 🚌 Reverse Proxy 
This documentation outlines the implementation of a basic reverse proxy in Go that utilizes a JSON configuration file to dynamically set up routing based on HTTP headers.

# 💻 Overview
The reverse proxy routes incoming HTTP requests to different backend services based on the value of a specified HTTP header. This allows for flexible and dynamic routing without hardcoding routes in the proxy itself.

# 🔧 JSON Configuration
The proxy's behavior is controlled by a JSON configuration file. This file defines the port on which the proxy listens and a list of endpoints with routing rules.

## Sample Configuration (config.json):

```json
{
    "port":"8000",
    "endpoints":[
        {
            "prefix":"/apiv1",
            "header_identifier":"country",
            "backend_urls":{"col":"http://localhost:8001","mex":"http://localhost:8002"}
        },
        {
            "prefix":"/apiv2",
            "header_identifier":"color",
            "backend_urls":{"blue":"http://localhost:8003","green":"http://localhost:8004"}
        }
    ]
}
```

## Explanation:

- port: Defines the port on which the reverse proxy listens.
- endpoints: A list of endpoint configurations.
- prefix: The URL prefix that the proxy listens for.
- header_identifier: The HTTP header used to determine the backend URL.
- backend_urls: A map where the key is the expected value of the header_identifier, and the value is the backend URL to which the request should be proxied.

# ✈ Implementation
The Go application reads the configuration file, sets up the routes, and proxies requests based on the provided configuration.

## Key Components:

1. Configuration Structs:
    - Config: Represents individual endpoint configuration.
    - ConfigFile: Represents the overall configuration file structure.

2. generateHandler Method:
    - Creates an HTTP handler for each endpoint that checks the specified header and routes to the appropriate backend URL.

3. reverseProxy Function:
    - Handles the actual proxying of requests to the target backend URL.

# Usage
1. Create the Configuration File:
    - Define the routing rules in a config.json file.
2. Run the Proxy:
    - Execute the Go program. The proxy will start on the port specified in the configuration file.
3. Send Requests:
    - Make HTTP requests to the proxy with the correct header and prefix. For example, if you send a request to /apiv1 with the header country: col, the proxy will route the request to http://localhost:8001.

## Example Request
```bash
curl -H "country: col" http://localhost:8000/apiv1
```
This request will be proxied to http://localhost:8001.

## Error Handling
- If the header value does not match any key in backend_urls, the proxy will return a 404 Not Found response.

# 🏢 Running the Project Locally
To run the project locally, you'll set up four simple servers using the simple_host application, each listening on a different port. Then, you'll start the reverse proxy that routes requests to these servers based on HTTP headers. Below are the steps to set up and run everything locally.

## Step 1: Run the Simple Servers
- Open four terminal windows or tabs.
- In each terminal, run the simple_host application with a different port:
```bash
go run simple_host.go -port=8001
go run simple_host.go -port=8002
go run simple_host.go -port=8003
go run simple_host.go -port=8004
```
Each server will respond with a message indicating its port when accessed.

## Step 2: Configure and Run the Reverse Proxy
1. Prepare the Configuration File:
    - Create a config.json file with the following content:
```json
{
    "port":"8000",
    "endpoints":[
        {
            "prefix":"/apiv1",
            "header_identifier":"country",
            "backend_urls":{"col":"http://localhost:8001","mex":"http://localhost:8002"}
        },
        {
            "prefix":"/apiv2",
            "header_identifier":"color",
            "backend_urls":{"blue":"http://localhost:8003","green":"http://localhost:8004"}
        }
    ]
}
```
2. Run the Reverse Proxy:
    - Save the reverse proxy code to a file (e.g., reverse_proxy.go).
    - Run the proxy server with the command:
```bash
go run reverse_proxy.go
```
The proxy will start on port 8000.

## Step 3: Test the Setup
1. Send Requests to the Proxy:
    - Use curl or a browser to send requests to the proxy. The proxy will route the requests to the appropriate backend server based on the headers.

### Example requests:

```bash
curl -H "country: col" http://localhost:8000/apiv1
```
This should return: `Hello from port 8001`

```bash
curl -H "color: blue" http://localhost:8000/apiv2
```
This should return: `Hello from port 8003`

