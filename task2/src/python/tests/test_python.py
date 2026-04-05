import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

import io
import json
import unittest
from main import AppHandler


class MockSocket:
    def __init__(self, request_data):
        self._data = request_data
        self._wbuf = io.BytesIO()

    def makefile(self, mode, bufsize=-1):
        if "r" in mode:
            return io.BytesIO(self._data)
        return self._wbuf

    def getpeername(self):
        return ("127.0.0.1", 8080)

    def sendall(self, data):
        self._wbuf.write(data)

    def get_response(self):
        self._wbuf.seek(0)
        return self._wbuf.read().decode()


def parse_response(raw):
    header_body = raw.split("\r\n\r\n", 1)
    return json.loads(header_body[1])


class TestHealthEndpoint(unittest.TestCase):
    def _handle(self, request_data):
        sock = MockSocket(request_data)
        AppHandler(sock, ("127.0.0.1", 8080), None)
        return sock.get_response()

    def test_health_status(self):
        response = self._handle(b"GET /health HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["status"], "ok")

    def test_health_message(self):
        response = self._handle(b"GET /health HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["message"], "Server is running")

    def test_health_default_port(self):
        if "PORT" in os.environ:
            del os.environ["PORT"]
        response = self._handle(b"GET /health HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["port"], "8080")

    def test_health_custom_port(self):
        os.environ["PORT"] = "9999"
        try:
            response = self._handle(b"GET /health HTTP/1.1\r\nHost: localhost\r\n\r\n")
            data = parse_response(response)
            self.assertEqual(data["port"], "9999")
        finally:
            del os.environ["PORT"]

    def test_health_content_type(self):
        response = self._handle(b"GET /health HTTP/1.1\r\nHost: localhost\r\n\r\n")
        headers = response.split("\r\n\r\n", 1)[0]
        self.assertIn("application/json", headers)

    def test_health_status_code(self):
        response = self._handle(b"GET /health HTTP/1.1\r\nHost: localhost\r\n\r\n")
        self.assertIn("200 OK", response)


class TestHelloEndpoint(unittest.TestCase):
    def _handle(self, request_data):
        sock = MockSocket(request_data)
        AppHandler(sock, ("127.0.0.1", 8080), None)
        return sock.get_response()

    def test_hello_default_name(self):
        response = self._handle(b"GET /hello HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["message"], "Hello, World!")

    def test_hello_with_name(self):
        response = self._handle(b"GET /hello?name=Rust HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["message"], "Hello, Rust!")

    def test_hello_with_special_name(self):
        response = self._handle(b"GET /hello?name=Go HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["message"], "Hello, Go!")

    def test_hello_empty_name(self):
        response = self._handle(b"GET /hello?name= HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["message"], "Hello, World!")

    def test_hello_status(self):
        response = self._handle(b"GET /hello HTTP/1.1\r\nHost: localhost\r\n\r\n")
        data = parse_response(response)
        self.assertEqual(data["status"], "ok")

    def test_hello_content_type(self):
        response = self._handle(b"GET /hello HTTP/1.1\r\nHost: localhost\r\n\r\n")
        headers = response.split("\r\n\r\n", 1)[0]
        self.assertIn("application/json", headers)

    def test_hello_status_code(self):
        response = self._handle(b"GET /hello HTTP/1.1\r\nHost: localhost\r\n\r\n")
        self.assertIn("200 OK", response)


class TestRouting(unittest.TestCase):
    def _handle(self, request_data):
        sock = MockSocket(request_data)
        AppHandler(sock, ("127.0.0.1", 8080), None)
        return sock.get_response()

    def test_root_404(self):
        response = self._handle(b"GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")
        self.assertIn("404", response)

    def test_unknown_404(self):
        response = self._handle(b"GET /unknown HTTP/1.1\r\nHost: localhost\r\n\r\n")
        self.assertIn("404", response)


if __name__ == "__main__":
    unittest.main()
