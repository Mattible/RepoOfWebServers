use std::env;
use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::thread;
use std::time::Duration;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ServerInfo {
    pub repository: String,
    pub url: String,
    #[serde(rename = "Programming Language")]
    pub programming_language: String,
    pub version: String,
    pub routes: Vec<RouteInfo>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct RouteInfo {
    pub path: String,
    pub description: String,
}

pub struct WebServer {
    port: String,
}

impl WebServer {
    pub fn new() -> Self {
        let port_env = env::var("WEBSERVER_PORT").unwrap_or_else(|_| "8080".to_string());
        let port = port_env.trim();
        let port = if port.is_empty() { "8080" } else { port };
        Self { port: port.to_string() }
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
            routes: vec![
                RouteInfo { path: "/".to_string(), description: "Hello World".to_string() },
                RouteInfo { path: "/health".to_string(), description: "Health check".to_string() },
                RouteInfo { path: "/info".to_string(), description: "Server info".to_string() },
                RouteInfo { path: "/image".to_string(), description: "Image handler".to_string() },
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

    pub fn run_server(&self) -> std::io::Result<()> {
        let listener = TcpListener::bind(format!("127.0.0.1:{}", self.port))?;
        println!("Server is listening on http://127.0.0.1:{}", self.port);
        println!("Press Ctrl+C to shutdown server...");

        let running = Arc::new(AtomicBool::new(true));
        let _running_clone = Arc::clone(&running);

        // Handle graceful shutdown in a separate thread
        thread::spawn(move || {
            // Simple signal handling - in a real implementation you'd use signal hooks
            thread::sleep(Duration::from_secs(1)); // Placeholder for signal handling
            // For now, we'll just wait - in production you'd listen for SIGINT/SIGTERM
        });

        // Accept connections while running
        for stream in listener.incoming() {
            if !running.load(Ordering::Relaxed) {
                println!("Shutting down server...");
                break;
            }

            match stream {
                Ok(stream) => {
                    let running_clone = Arc::clone(&running);
                    let server = self.clone();
                    thread::spawn(move || {
                        Self::handle_client(stream, server, running_clone);
                    });
                }
                Err(e) => {
                    eprintln!("Failed to accept connection: {}", e);
                    if !running.load(Ordering::Relaxed) {
                        break;
                    }
                }
            }
        }

        println!("Server shutdown gracefully");
        Ok(())
    }

    fn handle_client(mut stream: TcpStream, server: WebServer, _running: Arc<AtomicBool>) {
        let mut buffer = [0; 1024];
        match stream.read(&mut buffer) {
            Ok(n) => {
                if n == 0 {
                    return;
                }

                let request = String::from_utf8_lossy(&buffer[..n]);
                println!("Received request: {}", request.lines().next().unwrap_or(""));

                if let Some((method, path)) = Self::parse_http_request(&request) {
                    let response = server.handle_request(method, path);
                    if let Err(e) = stream.write_all(response.as_bytes()) {
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

impl Clone for WebServer {
    fn clone(&self) -> Self {
        Self {
            port: self.port.clone(),
        }
    }
}
