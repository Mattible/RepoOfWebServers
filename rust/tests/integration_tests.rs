use rust_webserver::{WebServer, ServerInfo};
use std::env;
use serial_test::serial;

// Tests moved from main.rs

#[test]
fn test_main_function_creates_server() {
    let server = WebServer::new();
    assert!(!server.get_port().is_empty());
}

#[test]
#[serial]
fn test_main_function_respects_environment() {
    // Store original PORT value
    let original_port = env::var("WEBSERVER_PORT").ok();
    
    // Test with PORT set
    unsafe { env::set_var("WEBSERVER_PORT", "4000"); }
    let server = WebServer::new();
    assert_eq!(server.get_port(), "4000");
    
    // Test with PORT unset
    unsafe { env::remove_var("WEBSERVER_PORT"); }
    let server = WebServer::new();
    assert_eq!(server.get_port(), "8000");
    
    // Restore original environment
    match original_port {
        Some(port) => unsafe { env::set_var("WEBSERVER_PORT", port) },
        None => unsafe { env::remove_var("WEBSERVER_PORT") },
    }
}

#[test]
fn test_main_server_functionality() {
    // Test that a server created like in main() works properly
    let server = WebServer::new();
    
    // Test that it can handle requests
    let response = server.handle_request("GET", "/");
    assert!(response.contains("Hello, world!"));
    
    let health_response = server.handle_request("GET", "/health");
    assert!(response.contains("200 OK"));
    assert!(health_response.contains("OK"));
    
    // Test server info
    let info = server.create_server_info();
    assert_eq!(info.repository, "RepoOfWebServers");
    assert_eq!(info.programming_language, "Rust");
}

#[test]
#[serial]
fn test_full_server_integration() {
    // Set up test environment with a specific test port
    unsafe { env::set_var("WEBSERVER_PORT", "8888"); }
    let server = WebServer::new();

    // Test server creation
    assert_eq!(server.get_port(), "8888");

    // Test server info creation
    let info = server.create_server_info();
    assert_eq!(info.repository, "RepoOfWebServers");
    assert_eq!(info.programming_language, "Rust");
    assert_eq!(info.routes.len(), 4);

    // Clean up
    unsafe { env::remove_var("WEBSERVER_PORT"); }
}

#[test]
#[serial]
fn test_environment_variable_handling() {
    // Test with WEBSERVER_PORT set
    unsafe { env::set_var("WEBSERVER_PORT", "9999"); }
    let server = WebServer::new();
    assert_eq!(server.get_port(), "9999");

    // Test with PORT unset
    unsafe { env::remove_var("WEBSERVER_PORT"); }
    let server = WebServer::new();
    assert_eq!(server.get_port(), "8000");
}

#[test]
fn test_all_endpoints_responses() {
    let server = WebServer::with_port("8000");

    // Test all endpoints
    let endpoints = vec![
        ("/", "Hello, world!"),
        ("/health", "OK"),
        ("/image", ""),
    ];

    for (path, expected_content) in endpoints {
        let response = server.handle_request("GET", path);
        assert!(response.contains("200 OK"));
        assert!(response.contains(expected_content));
    }

    // Test 404
    let response = server.handle_request("GET", "/nonexistent");
    assert!(response.contains("404 Not Found"));
}

#[test]
fn test_json_info_endpoint() {
    let server = WebServer::with_port("9090");
    let response = server.handle_request("GET", "/info");

    assert!(response.contains("200 OK"));
    assert!(response.contains("application/json"));

    // Extract JSON from response
    let json_start = response.find("\r\n\r\n").unwrap() + 4;
    let json = &response[json_start..];

    // Parse and validate JSON
    let info: ServerInfo = serde_json::from_str(json).unwrap();
    assert_eq!(info.routes.len(), 4);
    assert_eq!(info.routes[0].path, "/");
    assert_eq!(info.routes[0].description, "Hello World");
}

#[test]
fn test_http_method_handling() {
    let server = WebServer::with_port("8080");

    // Test POST request (should return 404)
    let response = server.handle_request("POST", "/health");
    assert!(response.contains("404 Not Found"));

    // Test PUT request (should return 404)
    let response = server.handle_request("PUT", "/");
    assert!(response.contains("404 Not Found"));

    // Test valid GET request
    let response = server.handle_request("GET", "/health");
    assert!(response.contains("200 OK"));
}

#[test]
fn test_server_info_serialization_full() {
    let server = WebServer::with_port("8080");
    let info = server.create_server_info();

    let json = serde_json::to_string(&info).unwrap();

    // Verify all expected fields are present with correct serde renames
    assert!(json.contains("\"Repository\":\"RepoOfWebServers\""));
    assert!(json.contains("\"URL\":\"https://github.com/Mattible/RepoOfWebServers\""));
    assert!(json.contains("\"Programming Language\":\"Rust\""));
    assert!(json.contains("\"version\":\"0.1.0\""));
    assert!(json.contains("\"gitSha\""));
    assert!(json.contains("\"routes\":["));

    // Verify we can deserialize it back
    let deserialized: ServerInfo = serde_json::from_str(&json).unwrap();
    assert_eq!(deserialized.routes.len(), 4);
}

#[test]
fn test_port_binding_validation() {
    // Test that we can create servers with different ports
    let server1 = WebServer::with_port("8000");
    let server2 = WebServer::with_port("9090");
    let server3 = WebServer::with_port("3000");

    assert_eq!(server1.get_port(), "8000");
    assert_eq!(server2.get_port(), "9090");
    assert_eq!(server3.get_port(), "3000");
}

#[test]
fn test_request_parsing_edge_cases() {
    // Test various malformed requests
    assert_eq!(WebServer::parse_http_request(""), None);
    assert_eq!(WebServer::parse_http_request("GET"), None);
    assert_eq!(WebServer::parse_http_request("GET /"), None);
    assert_eq!(WebServer::parse_http_request("GET / HTTP/1.1"), Some(("GET", "/")));
    assert_eq!(WebServer::parse_http_request("POST /api HTTP/1.1"), Some(("POST", "/api")));
}

#[test]
fn test_response_format_compliance() {
    let server = WebServer::with_port("8080");

    // All responses should follow HTTP format
    let responses = vec![
        server.handle_request("GET", "/"),
        server.handle_request("GET", "/health"),
        server.handle_request("GET", "/info"),
        server.handle_request("GET", "/nonexistent"),
    ];

    for response in responses {
        // Should start with HTTP version
        assert!(response.starts_with("HTTP/1.1"));

        // Should have status line
        assert!(response.contains(" "));

        // Should have headers
        assert!(response.contains("\r\n"));

        // Should have body separator
        assert!(response.contains("\r\n\r\n"));
    }
}

#[test]
fn test_concurrent_server_creation() {
    use std::sync::mpsc;
    use std::thread;

    let (tx, rx) = mpsc::channel();

    // Create multiple servers in parallel
    for i in 0..5 {
        let tx_clone = tx.clone();
        thread::spawn(move || {
            let port = format!("800{}", i);
            let server = WebServer::with_port(&port);
            tx_clone.send(server.get_port().to_string()).unwrap();
        });
    }

    // Collect results
    for _ in 0..5 {
        let port = rx.recv().unwrap();
        assert!(port.starts_with("800"));
    }
}

#[test]
fn test_stress_request_handling() {
    let server = WebServer::with_port("8000");
    
    // Test handling many requests in sequence
    for i in 0..100 {
        let path = if i % 4 == 0 { "/" } 
                  else if i % 4 == 1 { "/health" }
                  else if i % 4 == 2 { "/info" }
                  else { "/image" };
        
        let response = server.handle_request("GET", path);
        assert!(response.contains("200 OK"));
    }
}

#[test]
fn test_malformed_http_requests() {
    // Test various malformed HTTP requests
    let malformed_requests = vec![
        "",
        "GET",
        "GET /",
        "INVALID REQUEST FORMAT",
        "GET\n/\nHTTP/1.1",
        "   ",
        "\t\t\t",
        "GET / HTTP", // Missing version
        "/ HTTP/1.1", // Missing method
        "GET HTTP/1.1", // Missing path
    ];
    
    for request in malformed_requests {
        let result = WebServer::parse_http_request(request);
        if request.split_whitespace().count() >= 3 {
            assert!(result.is_some());
        } else {
            assert!(result.is_none());
        }
    }
}

#[test]
fn test_various_http_methods() {
    let server = WebServer::with_port("8080");
    let methods = vec!["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE"];
    
    for method in methods {
        let response = server.handle_request(method, "/");
        if method == "GET" {
            assert!(response.contains("200 OK"));
            assert!(response.contains("Hello, world!"));
        } else {
            assert!(response.contains("404 Not Found"));
        }
    }
}

#[test]
fn test_response_content_types() {
    let server = WebServer::with_port("8000");
    
    // Test text/plain responses
    let text_endpoints = vec!["/", "/health", "/image"];
    for endpoint in text_endpoints {
        let response = server.handle_request("GET", endpoint);
        assert!(response.contains("Content-Type: text/plain"));
    }
    
    // Test JSON response
    let json_response = server.handle_request("GET", "/info");
    assert!(json_response.contains("Content-Type: application/json"));
    
    // Test 404 response
    let not_found_response = server.handle_request("GET", "/nonexistent");
    assert!(not_found_response.contains("Content-Type: text/plain"));
}

#[test]
fn test_large_path_handling() {
    let server = WebServer::with_port("8000");
    
    // Test very long path
    let long_path = format!("/{}", "a".repeat(1000));
    let response = server.handle_request("GET", &long_path);
    assert!(response.contains("404 Not Found"));
    
    // Test path with special characters
    let special_paths = vec![
        "/path%20with%20spaces",
        "/path?query=value",
        "/path#fragment",
        "/path/with/many/segments",
        "/path-with-dashes",
        "/path_with_underscores",
        "/path.with.dots",
    ];
    
    for path in special_paths {
        let response = server.handle_request("GET", path);
        assert!(response.contains("404 Not Found"));
    }
}

#[test]
fn test_server_info_json_structure() {
    let server = WebServer::with_port("7777");
    let response = server.handle_request("GET", "/info");
    
    // Extract JSON from response
    let json_start = response.find("\r\n\r\n").unwrap() + 4;
    let json = &response[json_start..];
    
    // Parse JSON and verify structure
    let info: ServerInfo = serde_json::from_str(json).unwrap();
    
    // Test all fields are present and valid
    assert!(!info.repository.is_empty());
    assert!(!info.url.is_empty());
    assert!(!info.programming_language.is_empty());
    assert!(!info.version.is_empty());
    assert!(!info.routes.is_empty());
    
    // Test specific values
    assert_eq!(info.routes.len(), 4);
    
    // Test that routes have required fields
    for route in info.routes {
        assert!(!route.path.is_empty());
        assert!(!route.description.is_empty());
        assert!(route.path.starts_with('/'));
    }
}

#[test]
#[serial]
fn test_environment_variable_persistence() {
    use std::env;

    // Test that each server reads the environment variable at creation time
    let original = env::var("WEBSERVER_PORT").ok();

    // Set WEBSERVER_PORT and create server
    unsafe { env::set_var("WEBSERVER_PORT", "5555"); }
    let server1 = WebServer::new();
    assert_eq!(server1.get_port(), "5555");

    // Change WEBSERVER_PORT and create another server
    unsafe { env::set_var("WEBSERVER_PORT", "6666"); }
    let server2 = WebServer::new();
    assert_eq!(server2.get_port(), "6666");

    // First server should still have its original port (read at creation time)
    assert_eq!(server1.get_port(), "5555");
    assert_eq!(server2.get_port(), "6666");

    // Remove WEBSERVER_PORT and create server (should default to 8000)
    unsafe { env::remove_var("WEBSERVER_PORT"); }
    let server3 = WebServer::new();
    assert_eq!(server3.get_port(), "8000");

    // Previous servers should retain their ports (they don't re-read the env var)
    assert_eq!(server1.get_port(), "5555");
    assert_eq!(server2.get_port(), "6666");

    // Restore environment
    match original {
        Some(port) => unsafe { env::set_var("WEBSERVER_PORT", port) },
        None => unsafe { env::remove_var("WEBSERVER_PORT") },
    }
}

#[test]
fn test_clone_independence() {
    let server1 = WebServer::with_port("8001");
    let server2 = server1.clone();
    
    // Cloned servers should behave identically
    assert_eq!(server1.get_port(), server2.get_port());
    
    let response1 = server1.handle_request("GET", "/health");
    let response2 = server2.handle_request("GET", "/health");
    assert_eq!(response1, response2);
    
    // Test info consistency
    let info1 = server1.create_server_info();
    let info2 = server2.create_server_info();
    assert_eq!(info1.repository, info2.repository);
}

#[test]
#[serial]
fn test_gitsha_environment_variable() {
    // Store original GITSHA value
    let original_gitsha = env::var("GITSHA").ok();
    
    // Test with GITSHA set
    unsafe { env::set_var("GITSHA", "abc123def"); }
    let server = WebServer::new();
    let info = server.create_server_info();
    assert_eq!(info.git_sha, "abc123def");
    
    // Test with GITSHA unset (should default to "N/A")
    unsafe { env::remove_var("GITSHA"); }
    let server = WebServer::new();
    let info = server.create_server_info();
    assert_eq!(info.git_sha, "N/A");
    
    // Restore original environment
    match original_gitsha {
        Some(gitsha) => unsafe { env::set_var("GITSHA", gitsha) },
        None => unsafe { env::remove_var("GITSHA") },
    }
}
