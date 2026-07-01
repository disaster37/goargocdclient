package goargocdclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew_InvalidURL(t *testing.T) {
	_, err := New("://invalid-url")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestNew_MissingScheme(t *testing.T) {
	_, err := New("/relative/path/no/host")
	if err == nil {
		t.Error("expected error for missing scheme")
	}
}

func TestNew_WithToken(t *testing.T) {
	authChecked := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer test-token" {
			authChecked = true
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	api, err := New(server.URL, WithToken("test-token"))
	if err != nil {
		t.Fatal(err)
	}
	if api == nil {
		t.Fatal("expected non-nil API")
	}
	_ = authChecked
}

func TestNew_WithUsernamePassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/session" && r.Method == "POST" {
			var body map[string]string
			json.NewDecoder(r.Body).Decode(&body)
			if body["username"] != "admin" || body["password"] != "password" {
				w.WriteHeader(401)
				json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"token": "session-token"})
			return
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	api, err := New(server.URL, WithUsernamePassword("admin", "password"))
	if err != nil {
		t.Fatal(err)
	}
	if api == nil {
		t.Fatal("expected non-nil API")
	}
}

func TestNew_WithUsernamePasswordLoginFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]any{"code": float64(401), "message": "unauthorized"})
	}))
	defer server.Close()

	_, err := New(server.URL, WithUsernamePassword("admin", "wrong"))
	if err == nil {
		t.Error("expected error for failed login")
	}
}

func TestNew_WithUsernamePasswordInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("not-json"))
	}))
	defer server.Close()

	_, err := New(server.URL, WithUsernamePassword("admin", "password"))
	if err == nil {
		t.Error("expected error for invalid JSON response")
	}
}

func TestNew_WithUsernamePasswordNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	_, err := New(server.URL, WithUsernamePassword("admin", "password"))
	if err == nil {
		t.Error("expected error for network error during login")
	}
}

func TestWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	api, err := New(server.URL, WithTimeout(5*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if api == nil {
		t.Fatal("expected non-nil API")
	}
}

func TestWithInsecure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	api, err := New(server.URL, WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	if api == nil {
		t.Fatal("expected non-nil API")
	}
}
