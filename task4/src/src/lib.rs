use serde::Serialize;

#[derive(Serialize)]
pub struct Response {
    pub message: String,
    pub status: String,
    #[serde(skip_serializing_if = "String::is_empty")]
    pub port: String,
}

pub fn build_health_response() -> Response {
    let port = std::env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    Response {
        message: "Server is running".to_string(),
        status: "ok".to_string(),
        port,
    }
}

pub fn build_hello_response(path: &str) -> Response {
    let name = if let Some(pos) = path.find("name=") {
        let start = pos + 5;
        let end = path[start..].find('&').map(|i| start + i).unwrap_or(path.len());
        &path[start..end]
    } else {
        "World"
    };
    let name = if name.is_empty() { "World" } else { name };
    Response {
        message: format!("Hello, {}!", name),
        status: "ok".to_string(),
        port: String::new(),
    }
}

pub fn serialize_response(data: &impl Serialize) -> String {
    serde_json::to_string(data).unwrap()
}

pub fn build_response(status: u16, body: &str) -> String {
    let status_line = match status {
        200 => "HTTP/1.1 200 OK",
        404 => "HTTP/1.1 404 Not Found",
        _ => "HTTP/1.1 500 Internal Server Error",
    };
    format!(
        "{}\r\nContent-Type: application/json\r\nContent-Length: {}\r\nConnection: close\r\n\r\n{}",
        status_line,
        body.len(),
        body
    )
}

pub fn route(path: &str) -> (u16, String) {
    if path == "/health" {
        let resp = build_health_response();
        (200, serialize_response(&resp))
    } else if path.starts_with("/hello") {
        let resp = build_hello_response(path);
        (200, serialize_response(&resp))
    } else {
        (404, "404 Not Found".to_string())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_health_response() {
        std::env::remove_var("PORT");
        let resp = build_health_response();
        assert_eq!(resp.status, "ok");
        assert_eq!(resp.message, "Server is running");
        assert!(resp.port.contains("8080"));
        let json = serialize_response(&resp);
        assert!(json.contains("message"));
        assert!(json.contains("port"));
    }

    #[test]
    fn test_hello_default_name() {
        let resp = build_hello_response("/hello");
        assert_eq!(resp.message, "Hello, World!");
        assert_eq!(resp.status, "ok");
    }

    #[test]
    fn test_hello_with_name() {
        let resp = build_hello_response("/hello?name=Musl");
        assert_eq!(resp.message, "Hello, Musl!");
    }

    #[test]
    fn test_hello_empty_name() {
        let resp = build_hello_response("/hello?name=");
        assert_eq!(resp.message, "Hello, World!");
    }

    #[test]
    fn test_hello_multiple_params() {
        let resp = build_hello_response("/hello?name=Rust&extra=1");
        assert_eq!(resp.message, "Hello, Rust!");
    }

    #[test]
    fn test_route_health() {
        let (status, body) = route("/health");
        assert_eq!(status, 200);
        assert!(body.contains("Server is running"));
    }

    #[test]
    fn test_route_hello() {
        let (status, body) = route("/hello?name=Test");
        assert_eq!(status, 200);
        assert!(body.contains("Hello, Test!"));
    }

    #[test]
    fn test_route_unknown() {
        let (status, body) = route("/unknown");
        assert_eq!(status, 404);
        assert_eq!(body, "404 Not Found");
    }

    #[test]
    fn test_route_root() {
        let (status, body) = route("/");
        assert_eq!(status, 404);
        assert_eq!(body, "404 Not Found");
    }

    #[test]
    fn test_response_format_200() {
        let body = r#"{"ok":true}"#;
        let response = build_response(200, body);
        assert!(response.starts_with("HTTP/1.1 200 OK\r\n"));
        assert!(response.contains("Content-Type: application/json\r\n"));
        assert!(response.contains(body));
    }

    #[test]
    fn test_response_format_404() {
        let response = build_response(404, "404 Not Found");
        assert!(response.starts_with("HTTP/1.1 404 Not Found\r\n"));
    }

    #[test]
    fn test_json_serialization() {
        let resp = Response {
            message: "test".to_string(),
            status: "ok".to_string(),
            port: String::new(),
        };
        let json = serialize_response(&resp);
        assert!(json.contains(r#""message":"test""#));
        assert!(!json.contains("port"));
    }
}
