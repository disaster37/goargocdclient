package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Certificate interface {
	List(opts *CertificateQuery) (*CertificateList, error)
	Create(certs *CertificateCreateRequest, opts *CertificateCreateOptions) (*CertificateList, error)
	Delete(opts *CertificateQuery) error
}

type CertificateModel struct {
	ServerName  string `json:"serverName"`
	CertType    string `json:"certType"`
	CertSubType string `json:"certSubType"`
	CertData    string `json:"certData"`
	CertInfo    string `json:"certInfo"`
}

type CertificateList struct {
	Items []CertificateModel `json:"items"`
}

type CertificateCreateRequest struct {
	HTTPS   *CertificateModel `json:"https,omitempty"`
	SSH     *CertificateModel `json:"ssh,omitempty"`
	TLSCert *CertificateModel `json:"tls,omitempty"`
}

type CertificateQuery struct {
	HostNamePattern string `json:"hostNamePattern,omitempty"`
	CertType        string `json:"certType,omitempty"`
	CertSubType     string `json:"certSubType,omitempty"`
}

type CertificateCreateOptions struct {
	Upsert bool `json:"upsert,omitempty"`
}

type CertificateStandard struct {
	client *resty.Client
}

func NewCertificate(client *resty.Client) Certificate {
	return &CertificateStandard{client: client}
}

func (c *CertificateStandard) List(opts *CertificateQuery) (*CertificateList, error) {
	var result CertificateList
	req := c.client.R().SetResult(&result)
	if opts != nil {
		req = req.SetQueryParams(map[string]string{
			"hostNamePattern": opts.HostNamePattern,
			"certType":        opts.CertType,
			"certSubType":     opts.CertSubType,
		})
	}
	resp, err := req.Get("/api/v1/certificates")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (c *CertificateStandard) Create(certs *CertificateCreateRequest, opts *CertificateCreateOptions) (*CertificateList, error) {
	var result CertificateList
	req := c.client.R().
		SetBody(certs).
		SetResult(&result)
	if opts != nil && opts.Upsert {
		req = req.SetQueryParam("upsert", "true")
	}
	resp, err := req.Post("/api/v1/certificates")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (c *CertificateStandard) Delete(opts *CertificateQuery) error {
	req := c.client.R()
	if opts != nil {
		req = req.SetQueryParams(map[string]string{
			"hostNamePattern": opts.HostNamePattern,
			"certType":        opts.CertType,
			"certSubType":     opts.CertSubType,
		})
	}
	resp, err := req.Delete(fmt.Sprintf("/api/v1/certificates"))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

var _ Certificate = (*CertificateStandard)(nil)
