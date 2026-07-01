package api

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type Repository interface {
	List() (*RepositoryList, error)
	Get(repo string) (*RepositoryModel, error)
	Create(repo *RepositoryModel) (*RepositoryModel, error)
	Update(repo *RepositoryModel) (*RepositoryModel, error)
	Delete(repo string) error
	ListApps(repo string) (*RepositoryAppsList, error)
	GetAppDetails(opts *RepoAppDetailsQuery) (*RepoAppDetails, error)
	GetHelmCharts(repo string) (*HelmChartsResponse, error)
	ListRefs(repo string) (*RefsResponse, error)
	ListOCITags(repo string) (*OCITagsResponse, error)
	ValidateAccess(opts *RepoAccessQuery) error
	ListWriteRepositories() (*RepositoryList, error)
	GetWrite(repo string) (*RepositoryModel, error)
	CreateWriteRepository(repo *RepositoryModel) (*RepositoryModel, error)
	UpdateWriteRepository(repo *RepositoryModel) (*RepositoryModel, error)
	DeleteWriteRepository(repo string) error
	ValidateWriteAccess(opts *RepoAccessQuery) error
}

type RepositoryModel struct {
	TypeMeta
	ObjectMeta
	Repo          string `json:"repo"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	SSHPrivateKey string `json:"sshPrivateKey,omitempty"`
	Insecure      bool   `json:"insecure,omitempty"`
	EnableLFS     bool   `json:"enableLfs,omitempty"`
	EnableOCI     bool   `json:"enableOCI,omitempty"`
	Type          string `json:"type,omitempty"`
	Name          string `json:"name,omitempty"`
	Project       string `json:"project,omitempty"`
	InheritedCreds bool  `json:"inheritedCreds,omitempty"`
	TLSClientCertData string `json:"tlsClientCertData,omitempty"`
	TLSClientCertKey  string `json:"tlsClientCertKey,omitempty"`
	GCPServiceAccountKey string `json:"gcpServiceAccountKey,omitempty"`
	Proxy          string `json:"proxy,omitempty"`
	ForceHTTPBasicAuth bool `json:"forceHttpBasicAuth,omitempty"`
	ConnectionState ConnectionState `json:"connectionState,omitempty"`
}

type RepositoryList struct {
	Items []*RepositoryModel `json:"items"`
}

type RepositoryAppsList struct {
	Items []*RepoApp `json:"items"`
}

type RepoApp struct {
	RepoURL       string `json:"repoUrl"`
	Path          string `json:"path"`
	Type          string `json:"type"`
	AppName       string `json:"appName,omitempty"`
	ClusterName   string `json:"clusterName,omitempty"`
	NameSpace     string `json:"namespace,omitempty"`
}

type RepoAppDetailsQuery struct {
	Source    ApplicationSource       `json:"source"`
	AppName   string                  `json:"appName,omitempty"`
	AppProject   string               `json:"appProject,omitempty"`
}

type RepoAppDetails struct {
	Type        string            `json:"type"`
	Path        string            `json:"path"`
	Kustomize   *KustomizeAppSpec `json:"kustomize,omitempty"`
	Directory   *DirectoryAppSpec `json:"directory,omitempty"`
	Helm        *HelmAppSpec      `json:"helm,omitempty"`
	Plugin      *PluginAppSpec    `json:"plugin,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

type KustomizeAppSpec struct {
	Path  string   `json:"path,omitempty"`
	Images []string `json:"images,omitempty"`
}

type DirectoryAppSpec struct {
}

type HelmAppSpec struct {
	Name          string       `json:"name,omitempty"`
	ValueFiles    []string     `json:"valueFiles,omitempty"`
	Parameters    []HelmParameter `json:"parameters,omitempty"`
	FileParameters []HelmFileParameter `json:"fileParameters,omitempty"`
	Values        string       `json:"values,omitempty"`
}

type PluginAppSpec struct {
	Name       string           `json:"name,omitempty"`
	Env        []EnvEntry       `json:"env,omitempty"`
	Parameters []ParameterEntry `json:"parameters,omitempty"`
}

type HelmChartsResponse struct {
	Items []*HelmChart `json:"items"`
}

type HelmChart struct {
	Name        string   `json:"name"`
	Versions    []string `json:"versions"`
	Deprecated  bool     `json:"deprecated,omitempty"`
	Description string   `json:"description,omitempty"`
}

type RefsResponse struct {
	Branches []string `json:"branches"`
	Tags     []string `json:"tags"`
}

type OCITagsResponse struct {
	Tags []string `json:"tags"`
}

type RepoAccessQuery struct {
	Repo              string `json:"repo"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	SSHPrivateKey     string `json:"sshPrivateKey,omitempty"`
	TLSClientCertData string `json:"tlsClientCertData,omitempty"`
	TLSClientCertKey  string `json:"tlsClientCertKey,omitempty"`
	Type              string `json:"type,omitempty"`
	Name              string `json:"name,omitempty"`
	Insecure          bool   `json:"insecure,omitempty"`
	EnableOCI         bool   `json:"enableOCI,omitempty"`
	Proxy             string `json:"proxy,omitempty"`
	ForceHTTPBasicAuth bool  `json:"forceHttpBasicAuth,omitempty"`
}

type RepositoryStandard struct {
	client *resty.Client
}

func NewRepository(client *resty.Client) Repository {
	return &RepositoryStandard{client: client}
}

func (r *RepositoryStandard) List() (*RepositoryList, error) {
	var result RepositoryList
	resp, err := r.client.R().
		SetResult(&result).
		Get("/api/v1/repositories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) Get(repo string) (*RepositoryModel, error) {
	var result RepositoryModel
	resp, err := r.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) Create(repo *RepositoryModel) (*RepositoryModel, error) {
	var result RepositoryModel
	resp, err := r.client.R().
		SetBody(repo).
		SetResult(&result).
		Post("/api/v1/repositories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) Update(repo *RepositoryModel) (*RepositoryModel, error) {
	var result RepositoryModel
	resp, err := r.client.R().
		SetBody(repo).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/repositories/%s", url.PathEscape(repo.Repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) Delete(repo string) error {
	resp, err := r.client.R().
		Delete(fmt.Sprintf("/api/v1/repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (r *RepositoryStandard) ListApps(repo string) (*RepositoryAppsList, error) {
	var result RepositoryAppsList
	resp, err := r.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/repositories/%s/apps", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) GetAppDetails(opts *RepoAppDetailsQuery) (*RepoAppDetails, error) {
	var result RepoAppDetails
	resp, err := r.client.R().
		SetBody(opts).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/repositories/%s/appdetails", url.PathEscape(opts.Source.RepoURL)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) GetHelmCharts(repo string) (*HelmChartsResponse, error) {
	var result HelmChartsResponse
	resp, err := r.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/repositories/%s/helmcharts", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) ListRefs(repo string) (*RefsResponse, error) {
	var result RefsResponse
	resp, err := r.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/repositories/%s/refs", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) ListOCITags(repo string) (*OCITagsResponse, error) {
	var result OCITagsResponse
	resp, err := r.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/repositories/%s/oci/tags", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) ValidateAccess(opts *RepoAccessQuery) error {
	resp, err := r.client.R().
		SetBody(opts).
		Post(fmt.Sprintf("/api/v1/repositories/%s/validate", url.PathEscape(opts.Repo)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (r *RepositoryStandard) ListWriteRepositories() (*RepositoryList, error) {
	var result RepositoryList
	resp, err := r.client.R().
		SetResult(&result).
		Get("/api/v1/write-repositories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) GetWrite(repo string) (*RepositoryModel, error) {
	var result RepositoryModel
	resp, err := r.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/write-repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) CreateWriteRepository(repo *RepositoryModel) (*RepositoryModel, error) {
	var result RepositoryModel
	resp, err := r.client.R().
		SetBody(repo).
		SetResult(&result).
		Post("/api/v1/write-repositories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) UpdateWriteRepository(repo *RepositoryModel) (*RepositoryModel, error) {
	var result RepositoryModel
	resp, err := r.client.R().
		SetBody(repo).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/write-repositories/%s", url.PathEscape(repo.Repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) DeleteWriteRepository(repo string) error {
	resp, err := r.client.R().
		Delete(fmt.Sprintf("/api/v1/write-repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (r *RepositoryStandard) ValidateWriteAccess(opts *RepoAccessQuery) error {
	resp, err := r.client.R().
		SetBody(opts).
		Post(fmt.Sprintf("/api/v1/write-repositories/%s/validate", url.PathEscape(opts.Repo)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

var _ Repository = (*RepositoryStandard)(nil)
