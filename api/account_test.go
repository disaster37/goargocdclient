package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAccountList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, AccountsList{Items: []*AccountDetail{
			{Name: "admin"},
		}})
	}))
	defer server.Close()

	a := NewAccount(client)
	list, err := a.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Name != "admin" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestAccountGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, AccountDetail{Name: "admin"})
	}))
	defer server.Close()

	a := NewAccount(client)
	detail, err := a.Get("admin")
	if err != nil {
		t.Fatal(err)
	}
	if detail.Name != "admin" {
		t.Errorf("unexpected detail: %+v", detail)
	}
}

func TestAccountCanI_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, CanIResponse{Value: "yes"})
	}))
	defer server.Close()

	a := NewAccount(client)
	resp, err := a.CanI("applications", "create", "default")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Value != "yes" {
		t.Errorf("unexpected response: %s", resp.Value)
	}
}

func TestAccountUpdatePassword_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["currentPassword"] != "old" || body["newPassword"] != "new" {
			w.WriteHeader(400)
			return
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewAccount(client)
	if err := a.UpdatePassword("old", "new", "admin"); err != nil {
		t.Fatal(err)
	}
}

func TestAccountCreateToken_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, TokenResponse{Token: "tok123"})
	}))
	defer server.Close()

	a := NewAccount(client)
	resp, err := a.CreateToken("admin", 3600, "id1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token != "tok123" {
		t.Errorf("unexpected token: %s", resp.Token)
	}
}

func TestAccountDeleteToken_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewAccount(client)
	if err := a.DeleteToken("admin", "id1"); err != nil {
		t.Fatal(err)
	}
}

func TestAccountList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewAccount(client)
	_, err := a.List()
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccountGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewAccount(client)
	_, err := a.Get("admin")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccountCanI_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewAccount(client)
	_, err := a.CanI("applications", "create", "default")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccountUpdatePassword_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewAccount(client)
	err := a.UpdatePassword("old", "new", "admin")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccountCreateToken_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewAccount(client)
	_, err := a.CreateToken("admin", 3600, "id1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccountDeleteToken_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewAccount(client)
	err := a.DeleteToken("admin", "id1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccount_NetworkError(t *testing.T) {
	client := newFailingClient()
	a := NewAccount(client)
	_, err := a.List()
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Get("admin")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.CanI("applications", "create", "default")
	if err == nil {
		t.Error("expected error")
	}
	if err = a.UpdatePassword("old", "new", "admin"); err == nil {
		t.Error("expected error")
	}
	_, err = a.CreateToken("admin", 3600, "id1")
	if err == nil {
		t.Error("expected error")
	}
	if err = a.DeleteToken("admin", "id1"); err == nil {
		t.Error("expected error")
	}
}
