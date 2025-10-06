#!/usr/bin/env python3
import os
import logging
import signal
import sys
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import json
import threading
import time

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)


class WebServerHandler(BaseHTTPRequestHandler):
    """Custom HTTP request handler"""
    
    def do_GET(self):
        """Handle GET requests"""
        parsed_url = urlparse(self.path)
        path = parsed_url.path
        query_params = parse_qs(parsed_url.query)
        
        logger.info(f"GET {self.path} from {self.client_address[0]}")
        
        if path == '/':
            self._send_hello_response()
        elif path == '/health':
            self._send_health_response()
        elif path == '/info':
            self._send_info_response()
        elif path == '/image':
            self._send_image_response()
        else:
            self._send_404_response()
    
    def _send_hello_response(self):
        """Send hello world response"""
        message = "Hello World!\n"
        self._send_response(200, message, 'text/plain')
    
    def _send_health_response(self):
        """Send health check response"""
        self._send_response(200, "OK", 'text/plain')
    
    def _send_info_response(self):
        """Send server info response"""
        info_data = {
            "Programming Language": "Python",
            "Repository": "RepoOfWebServers",
            "URL": "https://github.com/Mattible/RepoOfWebServers",
            "version": "0.1.0",
            "git Sha": os.getenv("GITSHA", "N/A"),
            # "git Tag": os.getenv("GIT_TAG", "N/A"),
            "endpoints": [
                {"path": "/", "method": "GET", "description": "Hello world from Python"},
                {"path": "/health", "method": "GET", "description": "Health check"},
                {"path": "/info", "method": "GET", "description": "Server information"},
                {"path": "/image", "method": "GET", "description": "Image handler"},
            ]
        }
        self._send_json_response(200, info_data)
    
    def _send_image_response(self):
        """Send image handler response"""
        # TODO: Add image path to a cloud provided CDN Cache or local Storage
        message = ""
        self._send_response(200, message, 'text/plain')
    
    def _send_404_response(self):
        """Send 404 not found response"""
        error_data = {
            "error": "Not Found",
            "message": f"The requested path '{self.path}' was not found",
            "status_code": 404
        }
        self._send_json_response(404, error_data)
    
    def _send_response(self, status_code, message, content_type):
        """Send HTTP response with given status, message and content type"""
        self.send_response(status_code)
        self.send_header('Content-Type', content_type)
        self.send_header('Content-Length', str(len(message.encode())))
        self.end_headers()
        self.wfile.write(message.encode())
    
    def _send_json_response(self, status_code, data):
        """Send JSON response"""
        json_data = json.dumps(data, indent=2)
        self._send_response(status_code, json_data, 'application/json')
    
    def log_message(self, format, *args):
        """Override default log message to use our logger"""
        logger.info(f"{self.client_address[0]} - {format % args}")


class PythonWebServer:
    """Python Web Server class"""

    def __init__(self, host='0.0.0.0', port=8000):
        self.host = host
        self.port = port
        self.server = None
        self._shutdown_event = threading.Event()
        self._setup_signal_handlers()

    def _setup_signal_handlers(self):
        """Setup signal handlers for graceful shutdown"""
        signal.signal(signal.SIGINT, self._signal_handler)
        signal.signal(signal.SIGTERM, self._signal_handler)

    def _signal_handler(self, signum, frame):
        """Handle shutdown signals by setting the shutdown event."""
        logger.info(f"Received signal {signum}, shutting down...")
        self._shutdown_event.set()

    def start(self):
        """Start the web server and wait for shutdown signal."""
        try:
            self.server = HTTPServer((self.host, self.port), WebServerHandler)
            server_thread = threading.Thread(target=self.server.serve_forever)
            server_thread.daemon = True

            logger.info(f"Starting Python web server on {self.host}:{self.port}")
            server_thread.start()
            logger.info("Press Ctrl+C to shutdown...")

            # Wait for shutdown event
            self._shutdown_event.wait()

        except OSError as e:
            logger.error(f"Failed to start server: {e}")
            sys.exit(1)
        finally:
            self.stop()

    def stop(self):
        """Stop the web server."""
        if self.server:
            logger.info("Server shutting down gracefully...")
            self.server.shutdown()
            self.server.server_close()
            logger.info("Server stopped.")


def main():
    """Main function"""
    # Get configuration from environment variables
    host = os.getenv('HOST', '0.0.0.0')
    # Use WEBSERVER_PORT environment variable or default to 8000
    port = int(os.getenv('WEBSERVER_PORT', 8000))

    # Create and start server
    server = PythonWebServer(host, port)
    server.start()


if __name__ == '__main__':
    main()
