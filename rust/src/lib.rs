use async_std::net::{TcpListener, TcpStream};
use async_std::prelude::*;
use async_std::task;
use serde::{Deserialize, Serialize};
use std::env;
use std::sync::atomic::{AtomicBool, Ordering};

static SHUTDOWN_FLAG: AtomicBool = AtomicBool::new(false);

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ServerInfo {
    #[serde(rename = "Repository")]
    pub repository: String,
    #[serde(rename = "URL")]
    pub url: String,
    #[serde(rename = "Programming Language")]
    pub programming_language: String,
    pub version: String,
    #[serde(rename = "gitSha")]
    pub git_sha: String,
    pub routes: Vec<RouteInfo>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct RouteInfo {
    pub path: String,
    pub description: String,
}

#[derive(Clone)]
pub struct WebServer {
    port: String,
}

impl Default for WebServer {
    fn default() -> Self {
        Self::new()
    }
}

impl WebServer {
    pub fn new() -> Self {
        let port_env = env::var("WEBSERVER_PORT").unwrap_or_else(|_| "8000".to_string());
        let port = port_env.trim();
        let port = if port.is_empty() { "8000" } else { port };
        Self {
            port: port.to_string(),
        }
    }

    pub fn with_port(port: &str) -> Self {
        Self {
            port: port.to_string(),
        }
    }

    pub fn get_port(&self) -> &str {
        &self.port
    }

    pub fn create_server_info(&self) -> ServerInfo {
        ServerInfo {
            repository: "RepoOfWebServers".to_string(),
            url: "https://github.com/Mattible/RepoOfWebServers".to_string(),
            programming_language: "Rust".to_string(),
            version: "0.1.0".to_string(),
            git_sha: env::var("GITSHA").unwrap_or_else(|_| "N/A".to_string()),
            routes: vec![
                RouteInfo {
                    path: "/".to_string(),
                    description: "Hello World".to_string(),
                },
                RouteInfo {
                    path: "/health".to_string(),
                    description: "Health check".to_string(),
                },
                RouteInfo {
                    path: "/info".to_string(),
                    description: "Server info".to_string(),
                },
                RouteInfo {
                    path: "/image".to_string(),
                    description: "Image handler".to_string(),
                },
            ],
        }
    }

    pub fn handle_request(&self, method: &str, path: &str) -> String {
        match (method, path) {
            ("GET", "/") => "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, world!".to_string(),
            ("GET", "/health") => "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nOK".to_string(),
            ("GET", "/info") => {
                match serde_json::to_string(&self.create_server_info()) {
                    Ok(json) => format!("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{}", json),
                    Err(_) => "HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\n\r\nFailed to encode info".to_string(),
                }
            },
            ("GET", "/image") => "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n".to_string(),
            _ => "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\nNot Found".to_string(),
        }
    }

    pub fn parse_http_request(request: &str) -> Option<(&str, &str)> {
        let parts: Vec<&str> = request.split_whitespace().collect();
        if parts.len() >= 3 {
            Some((parts[0], parts[1]))
        } else {
            None
        }
    }

    pub async fn run_server(&self) -> std::io::Result<()> {
        let listener = TcpListener::bind(format!("0.0.0.0:{}", self.port)).await?;
        println!("Server is listening on http://0.0.0.0:{}", self.port);
        println!("Press Ctrl+C to shutdown server...");

        // Reset shutdown flag
        SHUTDOWN_FLAG.store(false, Ordering::Relaxed);

        // Handle graceful shutdown with proper signal handling
        ctrlc::set_handler(move || {
            if !SHUTDOWN_FLAG.load(Ordering::Relaxed) {
                SHUTDOWN_FLAG.store(true, Ordering::Relaxed);
                println!("\nReceived Ctrl+C, shutting down gracefully...");
            }
        })
        .expect("Error setting Ctrl+C handler");

        // Accept connections asynchronously
        let mut incoming = listener.incoming();

        loop {
            if SHUTDOWN_FLAG.load(Ordering::Relaxed) {
                break;
            }

            // Use timeout to periodically check shutdown flag
            match async_std::future::timeout(std::time::Duration::from_millis(100), incoming.next())
                .await
            {
                Ok(Some(stream_result)) => match stream_result {
                    Ok(stream) => {
                        let server = self.clone();
                        task::spawn(async move {
                            Self::handle_client(stream, server).await;
                        });
                    }
                    Err(e) => {
                        eprintln!("Failed to accept connection: {}", e);
                    }
                },
                Ok(None) => break,
                Err(_) => continue, // Timeout occurred, check shutdown flag again
            }
        }

        println!("Server shutdown gracefully");
        Ok(())
    }

    async fn handle_client(mut stream: TcpStream, server: WebServer) {
        let mut buffer = [0; 1024];
        match stream.read(&mut buffer).await {
            Ok(n) => {
                if n == 0 {
                    return;
                }

                let request = String::from_utf8_lossy(&buffer[..n]);
                println!("Received request: {}", request.lines().next().unwrap_or(""));

                if let Some((method, path)) = Self::parse_http_request(&request) {
                    let response = server.handle_request(method, path);
                    if let Err(e) = stream.write_all(response.as_bytes()).await {
                        eprintln!("Failed to write response: {}", e);
                    }
                }
            }
            Err(e) => {
                eprintln!("Failed to read from socket: {}", e);
            }
        }
    }
}
