package goargocdclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"disaster37/goargocdclient/api"

	"github.com/go-resty/resty/v2"
)

type Option func(*config)

type config struct {
	token    string
	username string
	password string
	insecure bool
	timeout  time.Duration
}

func WithToken(token string) Option {
	return func(c *config) {
		c.token = token
	}
}

func WithUsernamePassword(username, password string) Option {
	return func(c *config) {
		c.username = username
		c.password = password
	}
}

func WithInsecure() Option {
	return func(c *config) {
		c.insecure = true
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}

func New(serverURL string, opts ...Option) (api.API, error) {
	cfg := &config{
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}
	if parsedURL.Scheme == "" {
		return nil, fmt.Errorf("server URL must include scheme (http or https)")
	}

	client := resty.New().
		SetBaseURL(serverURL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTimeout(cfg.timeout)

	if cfg.insecure {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	if cfg.token != "" {
		client = client.SetAuthToken(cfg.token)
	} else if cfg.username != "" && cfg.password != "" {
		sessionResp, err := performLogin(client, cfg.username, cfg.password)
		if err != nil {
			return nil, fmt.Errorf("failed to login: %w", err)
		}
		client = client.SetAuthToken(sessionResp.Token)
	}

	return api.New(client), nil
}

func performLogin(client *resty.Client, username, password string) (*api.SessionResponse, error) {
	resp, err := client.R().
		SetBody(map[string]string{
			"username": username,
			"password": password,
		}).
		Post("/api/v1/session")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, api.ParseErrorFromBody(resp.StatusCode(), resp.Body())
	}
	var sessionResp api.SessionResponse
	if err := json.Unmarshal(resp.Body(), &sessionResp); err != nil {
		return nil, fmt.Errorf("failed to parse session response: %w", err)
	}
	return &sessionResp, nil
}