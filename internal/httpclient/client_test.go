package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alphabatem/flux_cli/dto"
)

func TestNew(t *testing.T) {
	c := New("https://example.com/", "mykey", "Authorization")

	if c.BaseURL != "https://example.com" {
		t.Errorf("expected trailing slash stripped, got %s", c.BaseURL)
	}
	if c.APIKey != "mykey" {
		t.Errorf("expected APIKey mykey, got %s", c.APIKey)
	}
	if c.HeaderName != "Authorization" {
		t.Errorf("expected HeaderName Authorization, got %s", c.HeaderName)
	}
	if c.HTTPClient.Timeout.Seconds() != 30 {
		t.Errorf("expected 30s timeout, got %v", c.HTTPClient.Timeout)
	}
}

func TestNew_NoTrailingSlash(t *testing.T) {
	c := New("https://example.com", "key", "X-API-KEY")
	if c.BaseURL != "https://example.com" {
		t.Errorf("expected no change, got %s", c.BaseURL)
	}
}

func TestGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/test/path" {
			t.Errorf("expected /test/path, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "test-key" {
			t.Errorf("expected Authorization header, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != "flux-cli/0.1.0" {
			t.Errorf("expected flux-cli User-Agent, got %s", r.Header.Get("User-Agent"))
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c := New(srv.URL, "test-key", "Authorization")
	var result map[string]string
	err := c.Get("/test/path", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %v", result)
	}
}

func TestGet_CustomHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != "rugcheck-key" {
			t.Errorf("expected X-API-KEY header, got %s", r.Header.Get("X-API-KEY"))
		}
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "rugcheck-key", "X-API-KEY")
	var result map[string]interface{}
	err := c.Get("/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_NoAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("expected no Authorization header when key is empty")
		}
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "", "Authorization")
	var result map[string]interface{}
	c.Get("/test", &result)
}

func TestGet_APIError(t *testing.T) {
	tests := []struct {
		status int
		body   string
	}{
		{400, `{"error": "bad request"}`},
		{401, `{"error": "unauthorized"}`},
		{403, `{"error": "forbidden"}`},
		{404, `{"error": "not found"}`},
		{500, `{"error": "internal"}`},
		{503, `{"error": "unavailable"}`},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("HTTP_%d", tt.status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := New(srv.URL, "key", "Authorization")
			var result interface{}
			err := c.Get("/test", &result)
			if err == nil {
				t.Fatal("expected error for non-2xx status")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if apiErr.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, apiErr.StatusCode)
			}
			if apiErr.Body != tt.body {
				t.Errorf("expected body %s, got %s", tt.body, apiErr.Body)
			}
		})
	}
}

func TestPost_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"test":true}` {
			t.Errorf("unexpected body: %s", string(body))
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]int{"id": 42})
	}))
	defer srv.Close()

	c := New(srv.URL, "key", "Authorization")
	var result map[string]int
	err := c.Post("/create", strings.NewReader(`{"test":true}`), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != 42 {
		t.Errorf("expected id 42, got %v", result)
	}
}

func TestPostRaw(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte(`raw response`))
	}))
	defer srv.Close()

	c := New(srv.URL, "key", "Authorization")
	body, status, err := c.PostRaw("/raw", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 201 {
		t.Errorf("expected 201, got %d", status)
	}
	if string(body) != "raw response" {
		t.Errorf("expected 'raw response', got %s", string(body))
	}
}

func TestDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "key", "Authorization")
	var result map[string]bool
	err := c.Delete("/resource/1", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["deleted"] {
		t.Error("expected deleted=true")
	}
}

func TestPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"updated":true}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "key", "Authorization")
	var result map[string]bool
	err := c.Put("/resource/1", strings.NewReader(`{"name":"new"}`), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["updated"] {
		t.Error("expected updated=true")
	}
}

func TestGet_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	c := New(srv.URL, "key", "Authorization")
	var result map[string]string
	err := c.Get("/test", &result)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGet_NilTarget(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"data":"ignored"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "key", "Authorization")
	err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("expected no error with nil target, got: %v", err)
	}
}

func TestAPIError_Error(t *testing.T) {
	e := &APIError{StatusCode: 401, Body: "unauthorized"}
	expected := "API error (HTTP 401): unauthorized"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestExitCodeForError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"auth 401", &APIError{StatusCode: 401}, dto.ExitAuthError},
		{"auth 403", &APIError{StatusCode: 403}, dto.ExitAuthError},
		{"unavailable 503", &APIError{StatusCode: 503}, dto.ExitServiceDown},
		{"api error 400", &APIError{StatusCode: 400}, dto.ExitAPIError},
		{"api error 404", &APIError{StatusCode: 404}, dto.ExitAPIError},
		{"api error 500", &APIError{StatusCode: 500}, dto.ExitAPIError},
		{"non-api error", fmt.Errorf("connection refused"), dto.ExitGeneralError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := ExitCodeForError(tt.err)
			if code != tt.expected {
				t.Errorf("expected exit code %d, got %d", tt.expected, code)
			}
		})
	}
}
