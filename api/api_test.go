package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestAPIInterface(t *testing.T) {
	client := resty.New()
	api := New(client)

	if api.Account() == nil {
		t.Error("Account() returned nil")
	}
	if api.Application() == nil {
		t.Error("Application() returned nil")
	}
	if api.ApplicationSet() == nil {
		t.Error("ApplicationSet() returned nil")
	}
	if api.Certificate() == nil {
		t.Error("Certificate() returned nil")
	}
	if api.Cluster() == nil {
		t.Error("Cluster() returned nil")
	}
	if api.GPGKey() == nil {
		t.Error("GPGKey() returned nil")
	}
	if api.Notification() == nil {
		t.Error("Notification() returned nil")
	}
	if api.Project() == nil {
		t.Error("Project() returned nil")
	}
	if api.RepoCreds() == nil {
		t.Error("RepoCreds() returned nil")
	}
	if api.Repository() == nil {
		t.Error("Repository() returned nil")
	}
	if api.Session() == nil {
		t.Error("Session() returned nil")
	}
	if api.Settings() == nil {
		t.Error("Settings() returned nil")
	}
	if api.Version() == nil {
		t.Error("Version() returned nil")
	}
}

func TestAPI_ServiceAccessors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	client := resty.New().SetBaseURL(srv.URL)
	api := New(client)

	_ = api.Account()
	_ = api.Application()
	_ = api.ApplicationSet()
	_ = api.Certificate()
	_ = api.Cluster()
	_ = api.GPGKey()
	_ = api.Notification()
	_ = api.Project()
	_ = api.RepoCreds()
	_ = api.Repository()
	_ = api.Session()
	_ = api.Settings()
	_ = api.Version()
}
