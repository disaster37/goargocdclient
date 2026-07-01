package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRepoCredsList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepoCredsList{Items: []*RepoCredsModel{
			{URL: "https://github.com", Username: "user"},
		}})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	list, err := r.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Username != "user" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestRepoCredsCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var creds RepoCredsModel
		json.NewDecoder(r.Body).Decode(&creds)
		jsonResponse(w, 201, creds)
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	creds, err := r.Create(&RepoCredsModel{URL: "https://github.com", Username: "user"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if creds.Username != "user" {
		t.Errorf("unexpected creds: %+v", creds)
	}
}

func TestRepoCredsUpdate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepoCredsModel{URL: "https://github.com", Username: "updated"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	creds, err := r.Update(&RepoCredsModel{URL: "https://github.com", Username: "updated"})
	if err != nil {
		t.Fatal(err)
	}
	if creds.Username != "updated" {
		t.Errorf("unexpected creds: %+v", creds)
	}
}

func TestRepoCredsDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	if err := r.Delete("https://github.com"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoCredsListWrite_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepoCredsList{Items: []*RepoCredsModel{
			{URL: "https://write.github.com"},
		}})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	list, err := r.ListWrite()
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestRepoCredsCreateWrite_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 201, RepoCredsModel{URL: "https://write.github.com", Username: "writer"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	creds, err := r.CreateWrite(&RepoCredsModel{URL: "https://write.github.com", Username: "writer"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if creds.Username != "writer" {
		t.Errorf("unexpected creds: %+v", creds)
	}
}

func TestRepoCredsUpdateWrite_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepoCredsModel{URL: "https://write.github.com", Username: "updated-writer"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	creds, err := r.UpdateWrite(&RepoCredsModel{URL: "https://write.github.com", Username: "updated-writer"})
	if err != nil {
		t.Fatal(err)
	}
	if creds.Username != "updated-writer" {
		t.Errorf("unexpected creds: %+v", creds)
	}
}

func TestRepoCredsDeleteWrite_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	if err := r.DeleteWrite("https://write.github.com"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoCredsList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	_, err := r.List()
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	_, err := r.Create(&RepoCredsModel{URL: "https://github.com", Username: "user"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsUpdate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	_, err := r.Update(&RepoCredsModel{URL: "https://github.com", Username: "updated"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	err := r.Delete("https://github.com")
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsListWrite_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	_, err := r.ListWrite()
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsCreateWrite_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	_, err := r.CreateWrite(&RepoCredsModel{URL: "https://write.github.com", Username: "writer"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsUpdateWrite_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	_, err := r.UpdateWrite(&RepoCredsModel{URL: "https://write.github.com", Username: "updated-writer"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCredsDeleteWrite_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepoCreds(client)
	err := r.DeleteWrite("https://write.github.com")
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepoCreds_NetworkError(t *testing.T) {
	client := newFailingClient()
	r := NewRepoCreds(client)
	_, err := r.List()
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.Create(&RepoCredsModel{URL: "https://github.com", Username: "user"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.Update(&RepoCredsModel{URL: "https://github.com", Username: "updated"})
	if err == nil {
		t.Error("expected error")
	}
	if err = r.Delete("https://github.com"); err == nil {
		t.Error("expected error")
	}
	_, err = r.ListWrite()
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.CreateWrite(&RepoCredsModel{URL: "https://write.github.com", Username: "writer"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.UpdateWrite(&RepoCredsModel{URL: "https://write.github.com", Username: "updated-writer"})
	if err == nil {
		t.Error("expected error")
	}
	if err = r.DeleteWrite("https://write.github.com"); err == nil {
		t.Error("expected error")
	}
}
