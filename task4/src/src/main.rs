use server::{build_response, route};
use std::io::{BufRead, BufReader, Write};
use std::net::{TcpListener, TcpStream};

fn parse_request(stream: &TcpStream) -> String {
    let reader = BufReader::new(stream.try_clone().unwrap());
    let request_line = reader.lines().next().unwrap().unwrap();
    let parts: Vec<&str> = request_line.split_whitespace().collect();
    if parts.len() > 1 {
        parts[1].to_string()
    } else {
        "/".to_string()
    }
}

fn handle_client(mut stream: TcpStream) {
    let path = parse_request(&stream);
    let (status, body) = route(&path);
    let response = build_response(status, &body);
    let _ = stream.write_all(response.as_bytes());
}

fn main() {
    let port = std::env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    let addr = format!("0.0.0.0:{}", port);
    let listener = TcpListener::bind(&addr).unwrap();

    for stream in listener.incoming() {
        if let Ok(stream) = stream {
            handle_client(stream);
        }
    }
}
