package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Account interface {
	List() (*AccountsList, error)
	Get(name string) (*AccountDetail, error)
	CanI(resource, action, subresource string) (*CanIResponse, error)
	UpdatePassword(currentPassword, newPassword, name string) error
	CreateToken(name string, expiresIn int64, id string) (*TokenResponse, error)
	DeleteToken(name, id string) error
}

type AccountDetail struct {
	Name      string   `json:"name"`
	Capabilities []string `json:"capabilities"`
	Tokens    []TokenInfo `json:"tokens,omitempty"`
}

type TokenInfo struct {
	ID        string `json:"id"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type AccountsList struct {
	Items []*AccountDetail `json:"items"`
}

type CanIResponse struct {
	Value string `json:"value"`
}

type AccountStandard struct {
	client *resty.Client
}

func NewAccount(client *resty.Client) Account {
	return &AccountStandard{client: client}
}

func (a *AccountStandard) List() (*AccountsList, error) {
	var result AccountsList
	resp, err := a.client.R().
		SetResult(&result).
		Get("/api/v1/account")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *AccountStandard) Get(name string) (*AccountDetail, error) {
	var result AccountDetail
	resp, err := a.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/account/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *AccountStandard) CanI(resource, action, subresource string) (*CanIResponse, error) {
	var result CanIResponse
	resp, err := a.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/account/can-i/%s/%s/%s", resource, action, subresource))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *AccountStandard) UpdatePassword(currentPassword, newPassword, name string) error {
	resp, err := a.client.R().
		SetBody(map[string]string{
			"currentPassword": currentPassword,
			"newPassword":     newPassword,
			"name":            name,
		}).
		Put("/api/v1/account/password")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *AccountStandard) CreateToken(name string, expiresIn int64, id string) (*TokenResponse, error) {
	var result TokenResponse
	resp, err := a.client.R().
		SetBody(map[string]any{
			"expiresIn": expiresIn,
			"id":        id,
		}).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/account/%s/token", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *AccountStandard) DeleteToken(name, id string) error {
	resp, err := a.client.R().
		Delete(fmt.Sprintf("/api/v1/account/%s/token/%s", name, id))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

var _ Account = (*AccountStandard)(nil)
