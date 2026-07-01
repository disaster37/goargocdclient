package api

import (
	"github.com/go-resty/resty/v2"
)

type Session interface {
	Create(username, password string) (*SessionResponse, error)
	Delete() error
	GetUserInfo() (*UserInfo, error)
}

type SessionResponse struct {
	Token string `json:"token"`
}

type UserInfo struct {
	LoggedIn bool   `json:"loggedIn"`
	Username string `json:"username,omitempty"`
	Issuer   string `json:"iss,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

type SessionStandard struct {
	client *resty.Client
}

func NewSession(client *resty.Client) Session {
	return &SessionStandard{client: client}
}

func (s *SessionStandard) Create(username, password string) (*SessionResponse, error) {
	var result SessionResponse
	resp, err := s.client.R().
		SetBody(map[string]string{
			"username": username,
			"password": password,
		}).
		SetResult(&result).
		Post("/api/v1/session")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (s *SessionStandard) Delete() error {
	resp, err := s.client.R().
		Delete("/api/v1/session")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (s *SessionStandard) GetUserInfo() (*UserInfo, error) {
	var result UserInfo
	resp, err := s.client.R().
		SetResult(&result).
		Get("/api/v1/session/userinfo")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

var _ Session = (*SessionStandard)(nil)
