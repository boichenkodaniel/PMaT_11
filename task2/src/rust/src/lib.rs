pub fn build_health_response() -> (u16, String) {
    let port = std::env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    (200, format!(r#"{{"message":"Server is running","status":"ok","port":"{}"}}"#, port))
}

pub fn build_hello_response(path: &str) -> (u16, String) {
    let name = if let Some(pos) = path.find("name=") {
        let start = pos + 5;
        let end = path[start..].find('&').map(|i| start + i).unwrap_or(path.len());
        &path[start..end]
    } else {
        "World"
    };
    let name = if name.is_empty() { "World" } else { name };
    (200, format!(r#"{{"message":"Hello, {}!","status":"ok"}}"#, name))
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
        build_health_response()
    } else if path.starts_with("/hello") {
        build_hello_response(path)
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
        let (status, body) = build_health_response();
        assert_eq!(status, 200);
        assert!(body.contains(r#""message":"Server is running""#));
        assert!(body.contains(r#""status":"ok""#));
        assert!(body.contains(r#""port":"#));
    }

    #[test]
    fn test_hello_default_name() {
        let (status, body) = build_hello_response("/hello");
        assert_eq!(status, 200);
        assert!(body.contains(r#""message":"Hello, World!""#));
    }

    #[test]
    fn test_hello_with_name() {
        let (status, body) = build_hello_response("/hello?name=Rust");
        assert_eq!(status, 200);
        assert!(body.contains(r#""message":"Hello, Rust!""#));
    }

    #[test]
    fn test_hello_empty_name() {
        let (status, body) = build_hello_response("/hello?name=");
        assert_eq!(status, 200);
        assert!(body.contains(r#""message":"Hello, World!""#));
    }

    #[test]
    fn test_hello_multiple_params() {
        let (status, body) = build_hello_response("/hello?name=Go&extra=1");
        assert_eq!(status, 200);
        assert!(body.contains(r#""message":"Hello, Go!""#));
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
        let response = build_response(200, r#"{"ok":true}"#);
        assert!(response.starts_with("HTTP/1.1 200 OK\r\n"));
        assert!(response.contains("Content-Type: application/json\r\n"));
        assert!(response.contains(r#"{"ok":true}"#));
    }

    #[test]
    fn test_response_format_404() {
        let response = build_response(404, "404 Not Found");
        assert!(response.starts_with("HTTP/1.1 404 Not Found\r\n"));
        assert!(response.contains("404 Not Found"));
    }
}
