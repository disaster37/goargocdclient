package api

import "github.com/go-resty/resty/v2"

type Version interface {
	Get() (*VersionInfo, error)
}

type VersionInfo struct {
	Version        string `json:"Version"`
	BuildDate      string `json:"BuildDate"`
	GitCommit      string `json:"GitCommit"`
	GitTreeState   string `json:"GitTreeState"`
	GoVersion      string `json:"GoVersion"`
	Compiler       string `json:"Compiler"`
	Platform       string `json:"Platform"`
	Kustomize      string `json:"KustomizeVersion"`
	Helm           string `json:"HelmVersion"`
	Kubectl        string `json:"KubectlVersion"`
	JSONnetVersion string `json:"JsonnetVersion"`
}

type VersionStandard struct {
	client *resty.Client
}

func NewVersion(client *resty.Client) Version {
	return &VersionStandard{client: client}
}

func (v *VersionStandard) Get() (*VersionInfo, error) {
	var result VersionInfo
	resp, err := v.client.R().
		SetResult(&result).
		Get("/api/version")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

var _ Version = (*VersionStandard)(nil)
