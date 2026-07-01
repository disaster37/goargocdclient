package api

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type RepoCreds interface {
	List() (*RepoCredsList, error)
	Create(creds *RepoCredsModel, opts *RepoCredsCreateOptions) (*RepoCredsModel, error)
	Update(creds *RepoCredsModel) (*RepoCredsModel, error)
	Delete(url string) error
	ListWrite() (*RepoCredsList, error)
	CreateWrite(creds *RepoCredsModel, opts *RepoCredsCreateOptions) (*RepoCredsModel, error)
	UpdateWrite(creds *RepoCredsModel) (*RepoCredsModel, error)
	DeleteWrite(url string) error
}

type RepoCredsModel struct {
	URL                        string `json:"url"`
	Username                   string `json:"username,omitempty"`
	Password                   string `json:"password,omitempty"`
	SSHPrivateKey              string `json:"sshPrivateKey,omitempty"`
	TLSClientCertData          string `json:"tlsClientCertData,omitempty"`
	TLSClientCertKey           string `json:"tlsClientCertKey,omitempty"`
	GCPServiceAccountKey       string `json:"gcpServiceAccountKey,omitempty"`
	GitHubAppPrivateKey        string `json:"githubAppPrivateKey,omitempty"`
	GitHubAppID                int64  `json:"githubAppID,omitempty"`
	GitHubAppInstallationID    int64  `json:"githubAppInstallationID,omitempty"`
	GitHubAppEnterpriseBaseURL string `json:"githubAppEnterpriseBaseUrl,omitempty"`
	Type                       string `json:"type,omitempty"`
	Name                       string `json:"name,omitempty"`
	EnableOCI                  bool   `json:"enableOCI,omitempty"`
	Project                    string `json:"project,omitempty"`
	Proxy                      string `json:"proxy,omitempty"`
}

type RepoCredsList struct {
	Items []*RepoCredsModel `json:"items"`
}

type RepoCredsCreateOptions struct {
	Upsert bool `json:"upsert,omitempty"`
}

type RepoCredsStandard struct {
	client *resty.Client
}

func NewRepoCreds(client *resty.Client) RepoCreds {
	return &RepoCredsStandard{client: client}
}

func (r *RepoCredsStandard) List() (*RepoCredsList, error) {
	var result RepoCredsList
	resp, err := r.client.R().
		SetResult(&result).
		Get("/api/v1/repocreds")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepoCredsStandard) Create(creds *RepoCredsModel, opts *RepoCredsCreateOptions) (*RepoCredsModel, error) {
	var result RepoCredsModel
	req := r.client.R().
		SetBody(creds).
		SetResult(&result)
	if opts != nil && opts.Upsert {
		req = req.SetQueryParam("upsert", "true")
	}
	resp, err := req.Post("/api/v1/repocreds")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepoCredsStandard) Update(creds *RepoCredsModel) (*RepoCredsModel, error) {
	var result RepoCredsModel
	resp, err := r.client.R().
		SetBody(creds).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/repocreds/%s", encodeRepoCredsURL(creds.URL)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepoCredsStandard) Delete(url string) error {
	resp, err := r.client.R().
		Delete(fmt.Sprintf("/api/v1/repocreds/%s", encodeRepoCredsURL(url)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (r *RepoCredsStandard) ListWrite() (*RepoCredsList, error) {
	var result RepoCredsList
	resp, err := r.client.R().
		SetResult(&result).
		Get("/api/v1/write-repocreds")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepoCredsStandard) CreateWrite(creds *RepoCredsModel, opts *RepoCredsCreateOptions) (*RepoCredsModel, error) {
	var result RepoCredsModel
	req := r.client.R().
		SetBody(creds).
		SetResult(&result)
	if opts != nil && opts.Upsert {
		req = req.SetQueryParam("upsert", "true")
	}
	resp, err := req.Post("/api/v1/write-repocreds")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepoCredsStandard) UpdateWrite(creds *RepoCredsModel) (*RepoCredsModel, error) {
	var result RepoCredsModel
	resp, err := r.client.R().
		SetBody(creds).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/write-repocreds/%s", encodeRepoCredsURL(creds.URL)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepoCredsStandard) DeleteWrite(url string) error {
	resp, err := r.client.R().
		Delete(fmt.Sprintf("/api/v1/write-repocreds/%s", encodeRepoCredsURL(url)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func encodeRepoCredsURL(target string) string {
	b := []byte(target)
	for i, ch := range b {
		switch ch {
		case '/':
			b[i] = '_'
		case ':':
			b[i] = '_'
		}
	}
	return url.QueryEscape(string(b))
}

var _ RepoCreds = (*RepoCredsStandard)(nil)
