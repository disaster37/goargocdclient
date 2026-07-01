package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCertificateList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, CertificateList{
			Items: []CertificateModel{{ServerName: "example.com", CertType: "https"}},
		})
	}))
	defer server.Close()

	c := NewCertificate(client)
	list, err := c.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].ServerName != "example.com" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestCertificateList_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, CertificateList{Items: []CertificateModel{{ServerName: "example.com"}}})
	}))
	defer server.Close()

	c := NewCertificate(client)
	list, err := c.List(&CertificateQuery{HostNamePattern: "*.com", CertType: "https"})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(list.Items))
	}
}

func TestCertificateCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 201, CertificateList{
			Items: []CertificateModel{{ServerName: "new.com", CertType: "https"}},
		})
	}))
	defer server.Close()

	c := NewCertificate(client)
	list, err := 	c.Create(&CertificateCreateRequest{
		HTTPS: &CertificateModel{ServerName: "new.com", CertType: "https", CertData: "data", CertInfo: "info"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].ServerName != "new.com" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestCertificateDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCertificate(client)
	if err := c.Delete(&CertificateQuery{HostNamePattern: "*.com", CertType: "https"}); err != nil {
		t.Fatal(err)
	}
}

func TestCertificateDelete_NilOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCertificate(client)
	if err := c.Delete(nil); err != nil {
		t.Fatal(err)
	}
}

func TestCertificateList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var query CertificateQuery
		json.NewDecoder(r.Body).Decode(&query)
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCertificate(client)
	_, err := c.List(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCertificateCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCertificate(client)
	_, err := 	c.Create(&CertificateCreateRequest{
		HTTPS: &CertificateModel{ServerName: "new.com", CertType: "https", CertData: "data", CertInfo: "info"},
	}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCertificateDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCertificate(client)
	err := c.Delete(&CertificateQuery{HostNamePattern: "*.com", CertType: "https"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestCertificateDelete_ErrorNilOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCertificate(client)
	err := c.Delete(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCertificate_NetworkError(t *testing.T) {
	client := newFailingClient()
	c := NewCertificate(client)
	_, err := c.List(nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = 	c.Create(&CertificateCreateRequest{
		HTTPS: &CertificateModel{ServerName: "new.com", CertType: "https", CertData: "data", CertInfo: "info"},
	}, nil)
	if err == nil {
		t.Error("expected error")
	}
	if err = c.Delete(&CertificateQuery{HostNamePattern: "*.com", CertType: "https"}); err == nil {
		t.Error("expected error")
	}
}
