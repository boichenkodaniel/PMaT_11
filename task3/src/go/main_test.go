package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/status", statusHandler)
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
}

func TestHealthHandlerPort(t *testing.T) {
	os.Setenv("PORT", "9999")
	defer os.Unsetenv("PORT")

	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Port != "9999" {
		t.Errorf("expected port '9999', got '%s'", resp.Port)
	}
}

func TestHealthHandlerDefaultPort(t *testing.T) {
	os.Unsetenv("PORT")

	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Port != "8080" {
		t.Errorf("expected port '8080', got '%s'", resp.Port)
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
	req := httptest.NewRequest(http.MethodGet, "/hello?name=Docker", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Hello, Docker!" {
		t.Errorf("expected 'Hello, Docker!', got '%s'", resp.Message)
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

func TestStatusHandlerRouteExists(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}

	var resp StatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Go != "ok" {
		t.Errorf("expected Go status 'ok', got '%s'", resp.Go)
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

func TestHelloResponseContentType(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}
}
