package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Port    string `json:"port"`
}

type StatusResponse struct {
	Go     string `json:"go"`
	Python string `json:"python"`
	Rust   string `json:"rust"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	resp := Response{
		Message: "Server is running",
		Status:  "ok",
		Port:    port,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	resp := Response{
		Message: fmt.Sprintf("Hello, %s!", name),
		Status:  "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 3 * time.Second}

	statuses := StatusResponse{
		Go:     "ok",
		Python: "unreachable",
		Rust:   "unreachable",
	}

	resp, err := client.Get("http://python:8080/health")
	if err == nil {
		body, e := io.ReadAll(resp.Body)
		resp.Body.Close()
		if e == nil {
			var h Response
			if json.Unmarshal(body, &h) == nil {
				statuses.Python = h.Status
			} else {
				statuses.Python = "error"
			}
		}
	}

	resp, err = client.Get("http://rust:8080/health")
	if err == nil {
		body, e := io.ReadAll(resp.Body)
		resp.Body.Close()
		if e == nil {
			var h Response
			if json.Unmarshal(body, &h) == nil {
				statuses.Rust = h.Status
			} else {
				statuses.Rust = "error"
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/status", statusHandler)

	http.ListenAndServe(":"+port, nil)
}
