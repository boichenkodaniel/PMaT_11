package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func waitForService(url string, timeout time.Duration) error {
	deadline := time.After(timeout)
	for {
		select {
		case <-deadline:
			return fmt.Errorf("service %s did not start in %s", url, timeout)
		default:
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func dockerCompose(args ...string) error {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Dir = "../"
	return cmd.Run()
}

func TestIntegrationNetwork(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	if !dockerAvailable() {
		t.Skip("docker not available — skipping integration test")
	}

	t.Log("Starting docker compose...")
	if err := dockerCompose("up", "-d", "--build"); err != nil {
		t.Fatalf("failed to start compose: %v", err)
	}
	defer func() {
		t.Log("Stopping docker compose...")
		dockerCompose("down")
	}()

	t.Log("Waiting for services to be ready...")
	time.Sleep(5 * time.Second)

	if err := waitForService("http://localhost:8080/health", 30*time.Second); err != nil {
		t.Fatalf("Go service not ready: %v", err)
	}

	t.Run("/status returns all services ok", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/status")
		if err != nil {
			t.Fatalf("failed to GET /status: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		t.Logf("/status response: %s", strings.TrimSpace(string(body)))

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var status StatusResponse
		if err := json.Unmarshal(body, &status); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		if status.Go != "ok" {
			t.Errorf("expected Go='ok', got '%s'", status.Go)
		}
		if status.Python != "ok" {
			t.Errorf("expected Python='ok', got '%s'", status.Python)
		}
		if status.Rust != "ok" {
			t.Errorf("expected Rust='ok', got '%s'", status.Rust)
		}
	})

	t.Run("/health endpoint works", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			t.Fatalf("failed to GET /health: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var r Response
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if r.Message != "Server is running" {
			t.Errorf("expected 'Server is running', got '%s'", r.Message)
		}
	})

	t.Run("/hello endpoint works", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/hello?name=Integration")
		if err != nil {
			t.Fatalf("failed to GET /hello: %v", err)
		}
		defer resp.Body.Close()

		var r Response
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if r.Message != "Hello, Integration!" {
			t.Errorf("expected 'Hello, Integration!', got '%s'", r.Message)
		}
	})
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
