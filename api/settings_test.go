package api

import (
	"net/http"
	"testing"
)

func TestRevisionList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, SettingsModel{URL: "https://argocd.local"})
	}))
	defer server.Close()

	s := NewSettings(client)
	settings, err := s.Get()
	if err != nil {
		t.Fatal(err)
	}
	if settings.URL != "https://argocd.local" {
		t.Errorf("unexpected settings: %+v", settings)
	}
}

func TestSettingsGetPlugins_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, PluginsModel{
			Plugins: []*PluginInfo{{Name: "plugin1"}},
		})
	}))
	defer server.Close()

	s := NewSettings(client)
	plugins, err := s.GetPlugins()
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins.Plugins) != 1 || plugins.Plugins[0].Name != "plugin1" {
		t.Errorf("unexpected plugins: %+v", plugins)
	}
}

func TestSettingsGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	s := NewSettings(client)
	_, err := s.Get()
	if err == nil {
		t.Error("expected error")
	}
}

func TestSettingsGetPlugins_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	s := NewSettings(client)
	_, err := s.GetPlugins()
	if err == nil {
		t.Error("expected error")
	}
}

func TestSettings_NetworkError(t *testing.T) {
	client := newFailingClient()
	s := NewSettings(client)
	_, err := s.Get()
	if err == nil {
		t.Error("expected error")
	}
	_, err = s.GetPlugins()
	if err == nil {
		t.Error("expected error")
	}
}
