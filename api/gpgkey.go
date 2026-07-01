package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type GPGKey interface {
	List() (*GPGKeyList, error)
	Get(keyID string) (*GPGKeyModel, error)
	Create(key *GPGKeyModel) (*GPGKeyModel, error)
	Delete(keyID string) error
}

type GPGKeyModel struct {
	KeyID      string `json:"keyID,omitempty"`
	KeyData    string `json:"keyData"`
	SubType    string `json:"subType,omitempty"`
}

type GPGKeyList struct {
	Items []*GPGKeyModel `json:"items"`
}

type GPGKeyStandard struct {
	client *resty.Client
}

func NewGPGKey(client *resty.Client) GPGKey {
	return &GPGKeyStandard{client: client}
}

func (g *GPGKeyStandard) List() (*GPGKeyList, error) {
	var result GPGKeyList
	resp, err := g.client.R().
		SetResult(&result).
		Get("/api/v1/gpgkeys")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (g *GPGKeyStandard) Get(keyID string) (*GPGKeyModel, error) {
	var result GPGKeyModel
	resp, err := g.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/gpgkeys/%s", keyID))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (g *GPGKeyStandard) Create(key *GPGKeyModel) (*GPGKeyModel, error) {
	var result GPGKeyModel
	resp, err := g.client.R().
		SetBody(key).
		SetResult(&result).
		Post("/api/v1/gpgkeys")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (g *GPGKeyStandard) Delete(keyID string) error {
	resp, err := g.client.R().
		Delete(fmt.Sprintf("/api/v1/gpgkeys/%s", keyID))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

var _ GPGKey = (*GPGKeyStandard)(nil)
