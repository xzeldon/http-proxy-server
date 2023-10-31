# Simple HTTP Tunneling Proxy Server in Go

This is a lightweight proxy server written in Go that supports HTTP requests and HTTPS tunneling.

## Features

- Proxying HTTP requests
- Tunneling HTTPS requests
- Optional Basic Proxy Authentication
- Customizable port, username, and password using command-line arguments
- Logging

## Requirements

- Go (version 1.x+)

## Usage

### Running the Proxy from source

1. Clone the repository:
```bash
git clone https://github.com/xzeldon/http-proxy-server.git
cd http-proxy-server
```

2. Run the proxy:
```bash
go run main.go --port <PORT> [--username <USERNAME> --password <PASSWORD>]
```

By default, the proxy will run on port `3000`. If both `--username` and `--password` are omitted, authentication will be bypassed.

For example:

- To run the proxy on port `1489` without authentication:
```bash
go run main.go --port 1489
```

- To run the proxy on port `1489` with authentication:
```bash
go run main.go --port 1489 --username admin --password admin123
```

### Authentication

If you specify both `--username` and `--password` when starting the proxy, it will enforce Basic Proxy Authentication with the given credentials. If these parameters are omitted, the proxy will not require authentication.