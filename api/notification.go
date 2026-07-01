package api

import (
	"github.com/go-resty/resty/v2"
)

type Notification interface {
	ListTriggers() (*NotificationTriggerList, error)
	ListServices() (*NotificationServiceList, error)
	ListTemplates() (*NotificationTemplateList, error)
}

type NotificationTrigger struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type NotificationTriggerList struct {
	Items []NotificationTrigger `json:"items"`
}

type NotificationService struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type NotificationServiceList struct {
	Items []NotificationService `json:"items"`
}

type NotificationTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type NotificationTemplateList struct {
	Items []NotificationTemplate `json:"items"`
}

type NotificationStandard struct {
	client *resty.Client
}

func NewNotification(client *resty.Client) Notification {
	return &NotificationStandard{client: client}
}

func (n *NotificationStandard) ListTriggers() (*NotificationTriggerList, error) {
	var result NotificationTriggerList
	resp, err := n.client.R().
		SetResult(&result).
		Get("/api/v1/notifications/triggers")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (n *NotificationStandard) ListServices() (*NotificationServiceList, error) {
	var result NotificationServiceList
	resp, err := n.client.R().
		SetResult(&result).
		Get("/api/v1/notifications/services")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (n *NotificationStandard) ListTemplates() (*NotificationTemplateList, error) {
	var result NotificationTemplateList
	resp, err := n.client.R().
		SetResult(&result).
		Get("/api/v1/notifications/templates")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

var _ Notification = (*NotificationStandard)(nil)
