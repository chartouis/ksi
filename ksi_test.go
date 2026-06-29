package ksi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) (response, error) {
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
	defer func() {
		t.Log(resp.Body.Close())
	}()

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

type createRequest struct {
	Name string `json:"name"`
}

func testHandlerWithBody(body createRequest) (response, error) {
	if body.Name == "" {
		return response{}, HTTPError{Status: 400, Message: "name is required"}
	}
	return Ok(map[string]string{"name": body.Name}), nil
}

func TestEndpointWithBody(t *testing.T) {
	time.Sleep(200 * time.Millisecond) // let TestEndpointCreation's server settle

	k := NewKsi(":8081")
	k.Post("/create", testHandlerWithBody)
	go func() {
		if err := k.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// valid request
	resp, err := http.Post("http://localhost:8081/create", "application/json",
		strings.NewReader(`{"name":"yzarr"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		t.Log(resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("couldn't decode body: %v", err)
	}
	if body["name"] != "yzarr" {
		t.Errorf("expected name yzarr, got %s", body["name"])
	}

	// empty body should still decode, name will be empty string -> 400
	resp2, err := http.Post("http://localhost:8081/create", "application/json",
		strings.NewReader(`{"name":""}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		t.Log(resp2.Body.Close())
	}()

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp2.StatusCode)
	}
}
