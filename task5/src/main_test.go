package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/hello", helloHandler)
	return mux
}

func TestHealthHandler(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
	if resp.Message != "Server is running" {
		t.Errorf("expected message 'Server is running', got '%s'", resp.Message)
	}
	if resp.Arch == "" {
		t.Error("expected arch to be set")
	}
}

func TestHelloHandlerDefault(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got '%s'", resp.Message)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestHelloHandlerWithName(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/hello?name=Buildx", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Hello, Buildx!" {
		t.Errorf("expected 'Hello, Buildx!', got '%s'", resp.Message)
	}
}

func TestHelloHandlerEmptyName(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/hello?name=", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Hello, World!" {
		t.Errorf("expected 'Hello, World!' for empty name, got '%s'", resp.Message)
	}
}

func TestRouting404(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status %d for unknown route, got %d", http.StatusNotFound, status)
	}
}

func TestRoutingRoot(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status %d for root route, got %d", http.StatusNotFound, status)
	}
}

func TestResponseContentType(t *testing.T) {
	mux := newTestMux()
	for _, path := range []string{"/health", "/hello"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
			}
		})
	}
}
