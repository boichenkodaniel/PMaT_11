package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Arch    string `json:"arch"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Message: "Server is running",
		Status:  "ok",
		Arch:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
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
		Arch:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/hello", helloHandler)

	http.ListenAndServe(":"+port, nil)
}
