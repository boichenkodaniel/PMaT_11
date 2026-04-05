import json
import os
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs


class AppHandler(BaseHTTPRequestHandler):
    def _send_json(self, data):
        body = json.dumps(data).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        parsed = urlparse(self.path)
        if parsed.path == "/health":
            self._handle_health()
        elif parsed.path == "/hello":
            self._handle_hello(parse_qs(parsed.query))
        else:
            self.send_response(404)
            self.end_headers()

    def _handle_health(self):
        port = os.environ.get("PORT", "8080")
        self._send_json({
            "message": "Server is running",
            "status": "ok",
            "port": port,
        })

    def _handle_hello(self, query):
        name = query.get("name", ["World"])[0]
        if not name:
            name = "World"
        self._send_json({
            "message": f"Hello, {name}!",
            "status": "ok",
        })

    def log_message(self, format, *args):
        pass


def main():
    port = int(os.environ.get("PORT", "8080"))
    server = HTTPServer(("0.0.0.0", port), AppHandler)
    server.serve_forever()


if __name__ == "__main__":
    main()
