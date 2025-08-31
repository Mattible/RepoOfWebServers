# Go Web Server

A simple, production-ready HTTP web server written in Go with graceful shutdown, logging, and multiple endpoints.

## Features

- ✅ **HTTP/1.1 and HTTP/2 Support** - Modern protocol support
- ✅ **Multiple Endpoints** - Health checks, info, and custom handlers
- ✅ **Graceful Shutdown** - Clean shutdown on SIGINT/SIGTERM signals
- ✅ **Request Logging** - Comprehensive logging with client IP and request details
- ✅ **Environment Configuration** - Configurable port via environment variables
- ✅ **JSON API Responses** - Structured responses for info endpoint
- ✅ **Comprehensive Tests** - Full test coverage for all endpoints

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/`      | Hello World response |
| GET    | `/health` | Health check (returns "OK") |
| GET    | `/info`  | Server information (JSON) |
| GET    | `/image` | Image handler (placeholder) |
| *      | `*`      | 404 handler for undefined routes |

## Quick Start

### Local Development

```bash
# Run with default port (8000)
go run main.go

# Run with custom port
WEBSERVER_PORT=3000 go run main.go

# Build and run
go build -o webserver main.go
./webserver
```

### Running Tests

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run tests with coverage
go test -cover

# Generate detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Docker Usage

### Build the Docker Image

```bash
docker build -t go-webserver .
```

### Run the Container

```bash
# Run with default port mapping
docker run -p 8000:8000 go-webserver

# Run with custom port
docker run -p 3000:3000 -e WEBSERVER_PORT=3000 go-webserver

```

## Configuration

The server can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBSERVER_PORT` | `8000` | Port number for the HTTP server |

## Examples

### Health Check

```bash
curl http://localhost:8000/health
# Response: OK
```

### Server Information

```bash
curl http://localhost:8000/info
# Response: JSON with server details
```

### WebSocket Test

```bash
# Using websocat (install with: cargo install websocat)
echo "Hello WebSocket" | websocat ws://localhost:8000/ws
```

## Development

### Project Structure

```
go/
├── main.go           # Main server implementation
├── main_test.go      # Comprehensive test suite
├── Dockerfile        # Container configuration
├── go.mod           # Go module dependencies
├── go.sum           # Dependency checksums
└── README.md        # This file
```

## Dependencies

This server uses Go standard library plus:
- `golang.org/x/net/http2` - HTTP/2 support

Install dependencies:
```bash
go mod tidy
```

## Performance

- **Concurrent Handling**: Built-in Go concurrency for request handling
- **HTTP/2 Support**: Multiplexed connections and server push
- **Efficient Routing**: Standard library HTTP mux
- **Memory Efficient**: Minimal memory footprint

## Production Deployment

### Recommended Settings

```bash
# Set production environment
export WEBSERVER_PORT=8080
export GOMAXPROCS=$(nproc)

# Run with systemd or similar process manager
./webserver
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.