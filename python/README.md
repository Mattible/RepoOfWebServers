# Python Web Server

A slightly overkill Python webserver, slithering in with graceful shutdown, logging, and multiple endpoints.

## Features

- ✅ **HTTP/1.1 Support** - Built on Python's http.server
- ✅ **Multiple Endpoints** - Health checks, info and custom handlers  
- ✅ **Graceful Shutdown** - Clean shutdown on SIGINT/SIGTERM signals
- ✅ **Request Logging** - Comprehensive logging with client IP and request details
- ✅ **Environment Configuration** - Configurable host and port via environment variables
- ✅ **JSON API Responses** - Structured JSON responses for most endpoints
- ✅ **Threading Support** - Non-blocking server operation


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
python server.py

# Run with custom port
WEBSERVER_PORT=3000 python server.py

# Run with custom host and port
HOST=127.0.0.1 WEBSERVER_PORT=3000 python server.py

# Make executable and run directly
chmod +x server.py
./server.py
```

## Docker Usage

### Build the Docker Image

```bash
docker build -t repo-ws-py.
```

### Run the Container

```bash
# Run with default port mapping
docker run -p 8000:8000 repo-ws-py

# Run with custom port
docker run -p 3000:3000 -e WEBSERVER_PORT=3000 repo-ws-py

# Run with custom host and port
docker run -p 3000:3000 -e HOST=0.0.0.0 -e WEBSERVER_PORT=3000 repo-ws-py
```

## Configuration

The server can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `0.0.0.0` | Host address to bind the server |
| `WEBSERVER_PORT` | `8000` | Port number for the HTTP server |
| `GITSHA` | `N/A` | Git SHA for version info (optional) |

## Examples

### Hello World

```bash
curl http://localhost:8000/
# Response: Hello  World!
```

### Health Check

```bash
curl http://localhost:8000/health
# Response: JSON with status, timestamp, and server info
```

### Server Information

```bash
curl http://localhost:8000/info
# Response: 
# {
#     {
#   "Programming Language": "Python",
#   "Repository": "RepoOfWebServers",
#   "URL": "https://github.com/Mattible/RepoOfWebServers",
#   "version": "0.1.0",
#   "git Sha": "xxxxxx",
#   "endpoints": []
# }
```

## Development

### Project Structure

```
python/
├── server.py         # Main server implementation
├── Dockerfile        # Container configuration (to be created)
├── requirements.txt  # Python dependencies (optional)
└── README.md        # This file
```

## Dependencies

This server uses Python 3 standard library:
- `http.server` - HTTP server implementation
- `threading` - Multi-threading support
- `logging` - Comprehensive logging
- `json` - JSON response handling
- `signal` - Graceful shutdown handling

### Recommended Settings

```bash
# Set production environment
export WEBSERVER_PORT=8080
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
