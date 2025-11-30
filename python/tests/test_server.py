#!/usr/bin/env python3
"""
Unit tests for Python Web Server
"""

import unittest
import json
import logging
import os
import sys
from unittest.mock import patch, MagicMock
from io import BytesIO
from server import WebServerHandler, PythonWebServer, main

# Add the server module to the path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


class TestWebServerHandler(unittest.TestCase):
    """Test cases for WebServerHandler class"""

    def setUp(self):
        """Set up test fixtures"""
        # Create a mock handler without calling the parent constructor
        self.handler = WebServerHandler.__new__(WebServerHandler)

        # Set up required attributes manually
        self.handler.client_address = ("127.0.0.1", 12345)
        self.handler.server = MagicMock()
        self.handler.server.server_name = "localhost"
        self.handler.server.server_port = 8000
        self.handler.wfile = BytesIO()
        self.handler.rfile = BytesIO()
        self.handler.path = "/"
        self.handler.requestline = "GET / HTTP/1.1"
        self.handler.request_version = "HTTP/1.1"

    def test_do_get_root_endpoint(self):
        """Test GET request to root endpoint"""
        self.handler.path = "/"
        with patch.object(self.handler, "_send_hello_response") as mock_hello:
            self.handler.do_GET()
            mock_hello.assert_called_once()

    def test_do_get_health_endpoint(self):
        """Test GET request to health endpoint"""
        self.handler.path = "/health"
        with patch.object(
            self.handler, "_send_health_response"
        ) as mock_health:
            self.handler.do_GET()
            mock_health.assert_called_once()

    def test_do_get_info_endpoint(self):
        """Test GET request to info endpoint"""
        self.handler.path = "/info"
        with patch.object(self.handler, "_send_info_response") as mock_info:
            self.handler.do_GET()
            mock_info.assert_called_once()

    def test_do_get_image_endpoint(self):
        """Test GET request to image endpoint"""
        self.handler.path = "/image"
        with patch.object(self.handler, "_send_image_response") as mock_image:
            self.handler.do_GET()
            mock_image.assert_called_once()

    def test_do_get_unknown_endpoint(self):
        """Test GET request to unknown endpoint returns 404"""
        self.handler.path = "/unknown"
        with patch.object(self.handler, "_send_404_response") as mock_404:
            self.handler.do_GET()
            mock_404.assert_called_once()

    def test_send_hello_response(self):
        """Test hello world response"""
        with patch.object(self.handler, "_send_response") as mock_send:
            self.handler._send_hello_response()
            mock_send.assert_called_once_with(
                200, "Hello World!\n", "text/plain"
            )

    def test_send_health_response(self):
        """Test health check response"""
        with patch.object(self.handler, "_send_response") as mock_send:
            self.handler._send_health_response()
            mock_send.assert_called_once_with(200, "OK", "text/plain")

    def test_send_info_response(self):
        """Test server info response"""
        with patch.dict(os.environ, {"GITSHA": "abc123"}):
            with patch.object(
                self.handler, "_send_json_response"
            ) as mock_json:
                self.handler._send_info_response()
                mock_json.assert_called_once()
                args = mock_json.call_args[0]
                self.assertEqual(args[0], 200)
                info_data = args[1]

                # Check required fields
                self.assertEqual(info_data["Programming Language"], "Python")
                self.assertEqual(info_data["Repository"], "RepoOfWebServers")
                self.assertEqual(info_data["gitSha"], "abc123")
                self.assertIn("endpoints", info_data)
                self.assertEqual(len(info_data["endpoints"]), 4)

    def test_send_info_response_no_gitsha(self):
        """Test server info response without GITSHA environment variable"""
        with patch.dict(os.environ, {}, clear=True):
            if "GITSHA" in os.environ:
                del os.environ["GITSHA"]
            with patch.object(
                self.handler, "_send_json_response"
            ) as mock_json:
                self.handler._send_info_response()
                args = mock_json.call_args[0]
                info_data = args[1]
                self.assertEqual(info_data["gitSha"], "N/A")

    def test_send_image_response(self):
        """Test image handler response"""
        with patch.object(self.handler, "_send_response") as mock_send:
            self.handler._send_image_response()
            mock_send.assert_called_once_with(200, "", "text/plain")

    def test_send_404_response(self):
        """Test 404 not found response"""
        self.handler.path = "/nonexistent"
        with patch.object(self.handler, "_send_json_response") as mock_json:
            self.handler._send_404_response()
            mock_json.assert_called_once()
            args = mock_json.call_args[0]
            self.assertEqual(args[0], 404)
            error_data = args[1]

            self.assertEqual(error_data["error"], "Not Found")
            self.assertEqual(error_data["status_code"], 404)
            self.assertIn("not found", error_data["message"].lower())

    def test_send_response_headers(self):
        """Test HTTP response headers are set correctly"""
        with patch.object(self.handler, "send_response") as mock_status:
            with patch.object(self.handler, "send_header") as mock_header:
                with patch.object(self.handler, "end_headers") as mock_end:
                    with patch.object(
                        self.handler.wfile, "write"
                    ) as mock_write:
                        self.handler._send_response(
                            200, "Test content", "text/plain"
                        )
                        mock_status.assert_called_once_with(200)
                        mock_header.assert_any_call(
                            "Content-Type", "text/plain"
                        )
                        mock_header.assert_any_call(
                            "Content-Length", "12"
                        )  # len("Test content")
                        mock_end.assert_called_once()
                        mock_write.assert_called_once_with(b"Test content")

    def test_send_json_response(self):
        """Test JSON response formatting"""
        test_data = {"key": "value", "number": 42}
        with patch.object(self.handler, "_send_response") as mock_send:
            self.handler._send_json_response(200, test_data)
            args = mock_send.call_args[0]
            self.assertEqual(args[0], 200)
            self.assertEqual(args[2], "application/json")
            # Verify JSON is properly formatted
            json_data = json.loads(args[1])
            self.assertEqual(json_data["key"], "value")
            self.assertEqual(json_data["number"], 42)


class TestPythonWebServer(unittest.TestCase):
    """Test cases for PythonWebServer class"""

    def test_server_initialization_default_values(self):
        """Test server initialization with default values"""
        server = PythonWebServer()
        self.assertEqual(server.host, "0.0.0.0")
        self.assertEqual(server.port, 8000)
        self.assertIsNone(server.server)

    def test_server_initialization_custom_values(self):
        """Test server initialization with custom values"""
        server = PythonWebServer(host="127.0.0.1", port=9090)
        self.assertEqual(server.host, "127.0.0.1")
        self.assertEqual(server.port, 9090)

    def test_signal_handler_setup(self):
        """Test signal handlers are properly set up"""
        with patch("signal.signal") as mock_signal:
            _ = PythonWebServer()
            # Signal should have been called twice (SIGINT and SIGTERM)
            self.assertEqual(mock_signal.call_count, 2)

    def test_stop_server(self):
        """Test server stop functionality"""
        server = PythonWebServer()
        mock_server = MagicMock()
        server.server = mock_server
        server.stop()
        mock_server.shutdown.assert_called_once()
        mock_server.server_close.assert_called_once()

    def test_stop_server_no_server_instance(self):
        """Test stop when no server instance exists"""
        server = PythonWebServer()
        # Should not raise an exception
        server.stop()

    def test_signal_handler_calls_stop(self):
        """Test signal handler sets shutdown event"""
        server = PythonWebServer()
        # Verify shutdown event is not set initially
        self.assertFalse(server._shutdown_event.is_set())
        # Call signal handler
        server._signal_handler(2, None)  # SIGINT = 2
        # Verify shutdown event is now set
        self.assertTrue(server._shutdown_event.is_set())


class TestMainFunction(unittest.TestCase):
    """Test cases for main function and environment handling"""

    def test_main_function_default_port(self):
        """Test main function with default port"""
        with patch.dict(os.environ, {}, clear=True):
            # Clear WEBSERVER_PORT if it exists
            if "WEBSERVER_PORT" in os.environ:
                del os.environ["WEBSERVER_PORT"]

            with patch("server.PythonWebServer") as mock_server_class:
                mock_server = MagicMock()
                mock_server_class.return_value = mock_server

                with patch("server.PythonWebServer.start"):
                    main()

                mock_server_class.assert_called_once_with("0.0.0.0", 8000)

    def test_main_function_webserver_port_env(self):
        """Test main function with WEBSERVER_PORT environment variable"""
        with patch.dict(os.environ, {"WEBSERVER_PORT": "9999"}):
            with patch("server.PythonWebServer") as mock_server_class:
                mock_server = MagicMock()
                mock_server_class.return_value = mock_server

                with patch("server.PythonWebServer.start"):
                    main()

                mock_server_class.assert_called_once_with("0.0.0.0", 9999)

    def test_main_function_webserver_port_only(self):
        """Test main function only uses WEBSERVER_PORT environment variable"""
        with patch.dict(os.environ, {"WEBSERVER_PORT": "5555"}):
            with patch("server.PythonWebServer") as mock_server_class:
                mock_server = MagicMock()
                mock_server_class.return_value = mock_server

                with patch("server.PythonWebServer.start"):
                    main()

                mock_server_class.assert_called_once_with("0.0.0.0", 5555)

    def test_main_function_custom_host(self):
        """Test main function with custom HOST environment variable"""
        with patch.dict(
            os.environ, {"HOST": "192.168.1.100", "WEBSERVER_PORT": "8080"}
        ):
            with patch("server.PythonWebServer") as mock_server_class:
                mock_server = MagicMock()
                mock_server_class.return_value = mock_server

                with patch("server.PythonWebServer.start"):
                    main()

                mock_server_class.assert_called_once_with(
                    "192.168.1.100", 8080
                )


class TestIntegration(unittest.TestCase):
    """Integration tests for server components"""

    def _create_mock_handler(self, path="/"):
        """Helper method to create a mock handler"""
        handler = WebServerHandler.__new__(WebServerHandler)
        handler.client_address = ("127.0.0.1", 12345)
        handler.server = MagicMock()
        handler.server.server_name = "localhost"
        handler.server.server_port = 8000
        handler.wfile = BytesIO()
        handler.rfile = BytesIO()
        handler.path = path
        handler.requestline = f"GET {path} HTTP/1.1"
        handler.request_version = "HTTP/1.1"
        return handler

    def test_request_logging_integration(self):
        """Test request logging works end-to-end"""
        handler = self._create_mock_handler("/")

        # Capture log output
        with patch("server.logger") as mock_logger:
            with patch.object(handler, "_send_hello_response"):
                handler.do_GET()

            # Verify logging was called
            mock_logger.info.assert_called()
            log_call = mock_logger.info.call_args[0][0]
            self.assertIn("GET /", log_call)
            self.assertIn("127.0.0.1", log_call)

    def test_json_serialization_integration(self):
        """Test JSON serialization works correctly"""
        handler = self._create_mock_handler("/")

        test_data = {"test": "value", "nested": {"key": 123}}

        with patch.object(handler, "send_response"):
            with patch.object(handler, "send_header"):
                with patch.object(handler, "end_headers"):
                    with patch.object(handler.wfile, "write") as mock_write:
                        handler._send_json_response(200, test_data)

                        # Verify JSON was written correctly
                        written_data = mock_write.call_args[0][0]
                        parsed_json = json.loads(written_data.decode())
                        self.assertEqual(parsed_json["test"], "value")
                        self.assertEqual(parsed_json["nested"]["key"], 123)

    def test_path_parsing_integration(self):
        """Test URL path parsing with query parameters"""
        handler = self._create_mock_handler("/info?param=value&other=123")

        with patch.object(handler, "_send_info_response") as mock_info:
            handler.do_GET()
            mock_info.assert_called_once()

    def test_error_handling_integration(self):
        """Test error handling across components"""
        handler = self._create_mock_handler(
            "/nonexistent/path/with/params?test=123"
        )

        with patch.object(handler, "_send_json_response") as mock_json:
            handler.do_GET()

            # Verify 404 response
            args = mock_json.call_args[0]
            self.assertEqual(args[0], 404)
            error_data = args[1]
            self.assertEqual(error_data["status_code"], 404)
            self.assertIn(
                "/nonexistent/path/with/params?test=123", error_data["message"]
            )


if __name__ == "__main__":
    # Configure test logging to be less verbose
    logging.getLogger("server").setLevel(logging.WARNING)

    # Run the tests
    unittest.main(verbosity=2)
