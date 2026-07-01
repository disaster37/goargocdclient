package api

import (
	"net/http"
	"testing"
)

func TestNotificationListTriggers_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, NotificationTriggerList{
			Items: []NotificationTrigger{{Name: "on-sync-status-unknown"}},
		})
	}))
	defer server.Close()

	n := NewNotification(client)
	triggers, err := n.ListTriggers()
	if err != nil {
		t.Fatal(err)
	}
	if len(triggers.Items) != 1 || triggers.Items[0].Name != "on-sync-status-unknown" {
		t.Errorf("unexpected triggers: %+v", triggers)
	}
}

func TestNotificationListServices_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, NotificationServiceList{
			Items: []NotificationService{{Name: "slack"}},
		})
	}))
	defer server.Close()

	n := NewNotification(client)
	services, err := n.ListServices()
	if err != nil {
		t.Fatal(err)
	}
	if len(services.Items) != 1 || services.Items[0].Name != "slack" {
		t.Errorf("unexpected services: %+v", services)
	}
}

func TestNotificationListTemplates_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, NotificationTemplateList{
			Items: []NotificationTemplate{{Name: "app-sync-succeeded"}},
		})
	}))
	defer server.Close()

	n := NewNotification(client)
	templates, err := n.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates.Items) != 1 || templates.Items[0].Name != "app-sync-succeeded" {
		t.Errorf("unexpected templates: %+v", templates)
	}
}

func TestNotificationListTriggers_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	n := NewNotification(client)
	_, err := n.ListTriggers()
	if err == nil {
		t.Error("expected error")
	}
}

func TestNotificationListServices_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	n := NewNotification(client)
	_, err := n.ListServices()
	if err == nil {
		t.Error("expected error")
	}
}

func TestNotificationListTemplates_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	n := NewNotification(client)
	_, err := n.ListTemplates()
	if err == nil {
		t.Error("expected error")
	}
}

func TestNotification_NetworkError(t *testing.T) {
	client := newFailingClient()
	n := NewNotification(client)
	_, err := n.ListTriggers()
	if err == nil {
		t.Error("expected error")
	}
	_, err = n.ListServices()
	if err == nil {
		t.Error("expected error")
	}
	_, err = n.ListTemplates()
	if err == nil {
		t.Error("expected error")
	}
}
