# Rust Web Server

A high-performance, memory-safe HTTP web server written in Rust with comprehensive testing, graceful shutdown, and JSON API responses.

## Features

- ✅ **Async/Await Runtime** - Built with async-std for high-performance concurrent I/O
- ✅ **Multiple Endpoints** - Health checks, info, and custom handlers
- ✅ **Graceful Shutdown** - Clean shutdown on SIGINT/SIGTERM signals
- ✅ **Request Logging** - Comprehensive logging with client IP and request details
- ✅ **Environment Configuration** - Configurable port via environment variables
- ✅ **JSON API Responses** - Structured responses for info endpoint
- ✅ **Comprehensive Testing** - 23 integration tests with 40+ test scenarios
- ✅ **Docker Support** - Multi-stage builds for optimized containers

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
cargo run

# Run with custom port
WEBSERVER_PORT=3000 cargo run

# Build optimized release version
cargo build --release
./target/release/main
```

### Running Tests

```bash
# Run all tests
cargo test

# Run with test coverage
cargo install cargo-tarpaulin
cargo tarpaulin --out Html

## Docker Usage

### Build the Docker Image

```bash
docker build -t repo-ws-rust .
```

### Run the Container

```bash
# Run with default port mapping
docker run -p 8000:8000 repo-ws-rust

# Run with custom port
docker run -p 3000:3000 -e WEBSERVER_PORT=3000 repo-ws-rust

# Run in detached mode
docker run -d -p 8000:8000 repo-ws-rust
```

## Configuration

The server can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBSERVER_PORT` | `8000` | Port number for the HTTP server |
| `GITSHA` | `""` | Git SHA for version tracking (optional) |

## Examples

### Health Check

```bash
curl http://localhost:8000/health
# Response: OK
```

### Server Information

```bash
curl http://localhost:8000/info
# Response:
# {
#   "Programming Language": "Rust",
#   "Repository": "RepoOfWebServers",
#   "URL": "https://github.com/Mattible/RepoOfWebServers",
#   "version": "0.1.0",
#   "gitSha": "xxxxxx",
#   "routes": []
# }
```

### Basic Request

```bash
curl http://localhost:8000/
# Response: Hello, world!
```

## Development

### Project Structure

```
rust/
├── src/
│   ├── main.rs              # Application entry point
│   └── lib.rs               # Core server implementation
├── tests/
│   └── integration_tests.rs # Comprehensive test suite
├── Cargo.toml              # Rust project configuration
├── Dockerfile              # Multi-stage container build
└── README.md               # This file
```

### Dependencies

```toml
[dependencies]
serde = { version = "1.0", features = ["derive"] }  # Serialization/deserialization framework
serde_json = "1.0"                                  # JSON serialization support
async-std = { version = "1.12", features = ["attributes"] } # Async runtime for efficient I/O
ctrlc = "3.4"                                        # Cross-platform signal handling for graceful shutdown

[dev-dependencies]
serial_test = "3.0"                                  # Test isolation for integration tests
```

## License

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.