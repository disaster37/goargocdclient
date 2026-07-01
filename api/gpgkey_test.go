package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGPGKeyList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, GPGKeyList{Items: []*GPGKeyModel{
			{KeyID: "abc123", KeyData: "key-data"},
		}})
	}))
	defer server.Close()

	g := NewGPGKey(client)
	list, err := g.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].KeyID != "abc123" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestGPGKeyGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, GPGKeyModel{KeyID: "abc123", KeyData: "key-data"})
	}))
	defer server.Close()

	g := NewGPGKey(client)
	key, err := g.Get("abc123")
	if err != nil {
		t.Fatal(err)
	}
	if key.KeyID != "abc123" {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestGPGKeyCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key GPGKeyModel
		json.NewDecoder(r.Body).Decode(&key)
		jsonResponse(w, 201, key)
	}))
	defer server.Close()

	g := NewGPGKey(client)
	key, err := g.Create(&GPGKeyModel{KeyData: "new-key-data"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if key.KeyData != "new-key-data" {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestGPGKeyDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	g := NewGPGKey(client)
	if err := g.Delete("abc123"); err != nil {
		t.Fatal(err)
	}
}

func TestGPGKeyGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	g := NewGPGKey(client)
	_, err := g.Get("nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestGPGKeyList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	g := NewGPGKey(client)
	_, err := g.List()
	if err == nil {
		t.Error("expected error")
	}
}

func TestGPGKeyCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	g := NewGPGKey(client)
	_, err := g.Create(&GPGKeyModel{KeyData: "new-key-data"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestGPGKeyDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	g := NewGPGKey(client)
	err := g.Delete("abc123")
	if err == nil {
		t.Error("expected error")
	}
}

func TestGPGKey_NetworkError(t *testing.T) {
	client := newFailingClient()
	g := NewGPGKey(client)
	_, err := g.List()
	if err == nil {
		t.Error("expected error")
	}
	_, err = g.Get("abc123")
	if err == nil {
		t.Error("expected error")
	}
	_, err = g.Create(&GPGKeyModel{KeyData: "new-key-data"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	if err = g.Delete("abc123"); err == nil {
		t.Error("expected error")
	}
}
