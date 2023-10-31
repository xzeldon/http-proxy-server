# Simple HTTP Tunneling Proxy Server in Go

This is a lightweight proxy server written in Go that supports HTTP requests and HTTPS tunneling.

## Features

- Proxying HTTP requests
- Tunneling HTTPS requests
- Optional Basic Proxy Authentication
- Customizable port, username, and password using command-line arguments
- Logging

## Prerequisites

- Golang (if building from source)

## Usage

### Running the Proxy

#### For Darwin/Linux:

1. Download the latest release for your platform from GitHub releases
2. Extract the downloaded archive:

   ```bash
   tar -xvzf http-proxy-server-v{version}-{platform}-{arch}.tar.gz
   ```

3. Navigate to the extracted folder:

   ```bash
   cd http-proxy-server-v{version}-{platform}-{arch}
   ```

4. Run the proxy:

   ```bash
   ./http-proxy-server --port <PORT> [--username <USERNAME> --password <PASSWORD>]
   ```

#### For Windows:

1. Download the latest release for Windows from GitHub releases
2. Extract the downloaded ZIP archive
3. Navigate to the extracted folder.
4. Open a command prompt in this directory.
5. Run the proxy:

   ```bash
   proxy_name.exe --port <PORT> [--username <USERNAME> --password <PASSWORD>]
   ```

#### From Source

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

### License

This project is licensed under the MIT License.