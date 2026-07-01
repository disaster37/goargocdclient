package api

import (
	"net/http"
	"testing"
)

func TestVersionGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, VersionInfo{Version: "v2.8.0", GoVersion: "go1.21"})
	}))
	defer server.Close()

	v := NewVersion(client)
	info, err := v.Get()
	if err != nil {
		t.Fatal(err)
	}
	if info.Version != "v2.8.0" {
		t.Errorf("expected v2.8.0, got %s", info.Version)
	}
	if info.GoVersion != "go1.21" {
		t.Errorf("expected go1.21, got %s", info.GoVersion)
	}
}

func TestVersionGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	v := NewVersion(client)
	_, err := v.Get()
	if err == nil {
		t.Error("expected error")
	}
}

func TestVersion_NetworkError(t *testing.T) {
	client := newFailingClient()
	v := NewVersion(client)
	_, err := v.Get()
	if err == nil {
		t.Error("expected error")
	}
}
