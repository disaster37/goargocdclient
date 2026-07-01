package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Session interface {
	Create(opts *SessionCreateOptions) (*SessionResponse, error)
	Delete() error
	GetUserInfo() (*UserInfo, error)
}

type SessionResponse struct {
	Token string `json:"token"`
}

type UserInfo struct {
	LoggedIn bool     `json:"loggedIn"`
	Username string   `json:"username,omitempty"`
	Issuer   string   `json:"iss,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

type SessionCreateOptions struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

func (o *SessionCreateOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	if o.Token != "" {
		return nil
	}
	if o.Username == "" || o.Password == "" {
		return fmt.Errorf("either token or username+password must be provided")
	}
	return nil
}

type SessionStandard struct {
	client *resty.Client
}

func NewSession(client *resty.Client) Session {
	return &SessionStandard{client: client}
}

func (s *SessionStandard) Create(opts *SessionCreateOptions) (*SessionResponse, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	var result SessionResponse
	req := s.client.R().SetResult(&result)
	if opts.Token != "" {
		req = req.SetHeader("Authorization", "Bearer "+opts.Token)
	} else {
		req = req.SetBody(map[string]string{
			"username": opts.Username,
			"password": opts.Password,
		})
	}
	resp, err := req.Post("/api/v1/session")
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
