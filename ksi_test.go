package ksi

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) (Response, error) {
	log.Printf("method: %s, path: %s", r.Method, r.URL.Path)
	log.Printf("user-agent: %s", r.Header.Get("User-Agent"))
	w.Header().Set("X-Test", "ksi-reflection")
	return Ok(map[string]string{"status": "ok", "path": r.URL.Path}), nil
}

func TestEndpointCreation(t *testing.T) {
	k := NewKsi(":8080")
	k.Get("/test", testHandler)

	go func() {
		if err := k.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // wait for server to start

	resp, err := http.Get("http://localhost:8080/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("couldn't decode body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}
	if body["path"] != "/test" {
		t.Errorf("expected path /test, got %s", body["path"])
	}

	if resp.Header.Get("X-Test") != "ksi-reflection" {
		t.Errorf("expected X-Test header, got %s", resp.Header.Get("X-Test"))
	}
}
