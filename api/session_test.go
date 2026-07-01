package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestSessionCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/session" {
			w.WriteHeader(404)
			return
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["username"] != "admin" || body["password"] != "pass" {
			w.WriteHeader(403)
			return
		}
		jsonResponse(w, 200, SessionResponse{Token: "abc123"})
	}))
	defer server.Close()

	s := NewSession(client)
	resp, err := s.Create(&SessionCreateOptions{Username: "admin", Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token != "abc123" {
		t.Errorf("expected token abc123, got %s", resp.Token)
	}
}

func TestSessionCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 401, APIError{Code: 401, Message: "bad credentials"})
	}))
	defer server.Close()

	s := NewSession(client)
	_, err := s.Create(&SessionCreateOptions{Username: "admin", Password: "wrong"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestSessionDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" || r.URL.Path != "/api/v1/session" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	s := NewSession(client)
	if err := s.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestSessionDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	s := NewSession(client)
	if err := s.Delete(); err == nil {
		t.Error("expected error")
	}
}

func TestSessionGetUserInfo_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, UserInfo{LoggedIn: true, Username: "admin"})
	}))
	defer server.Close()

	s := NewSession(client)
	info, err := s.GetUserInfo()
	if err != nil {
		t.Fatal(err)
	}
	if !info.LoggedIn || info.Username != "admin" {
		t.Errorf("unexpected user info: %+v", info)
	}
}

func TestSessionGetUserInfo_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 403, APIError{Code: 403, Message: "forbidden"})
	}))
	defer server.Close()

	s := NewSession(client)
	_, err := s.GetUserInfo()
	if err == nil {
		t.Error("expected error")
	}
}

func TestSession_NetworkError(t *testing.T) {
	client := newFailingClient()
	s := NewSession(client)
	_, err := s.Create(&SessionCreateOptions{Username: "admin", Password: "pass"})
	if err == nil {
		t.Error("expected error")
	}
	if err = s.Delete(); err == nil {
		t.Error("expected error")
	}
	_, err = s.GetUserInfo()
	if err == nil {
		t.Error("expected error")
	}
}
