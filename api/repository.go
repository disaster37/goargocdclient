package api

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type Repository interface {
	List(opts *RepositoryQueryOptions) (*RepositoryList, error)
	Get(repo string, opts *RepositoryQueryOptions) (*RepositoryModel, error)
	Create(repo *RepositoryModel, opts *RepositoryCreateOptions) (*RepositoryModel, error)
	Update(repo *RepositoryModel) (*RepositoryModel, error)
	Delete(repo string, opts *RepositoryQueryOptions) error
	ListApps(repo string, opts *RepoAppsQueryOptions) (*RepositoryAppsList, error)
	GetAppDetails(opts *RepoAppDetailsQuery) (*RepoAppDetails, error)
	GetHelmCharts(repo string, opts *RepositoryQueryOptions) (*HelmChartsResponse, error)
	ListRefs(repo string, opts *RepositoryQueryOptions) (*RefsResponse, error)
	ListOCITags(repo string, opts *RepositoryQueryOptions) (*OCITagsResponse, error)
	ValidateAccess(opts *RepoAccessQuery) error
	ListWriteRepositories(opts *RepositoryQueryOptions) (*RepositoryList, error)
	GetWrite(repo string, opts *RepositoryQueryOptions) (*RepositoryModel, error)
	CreateWriteRepository(repo *RepositoryModel, opts *RepositoryCreateOptions) (*RepositoryModel, error)
	UpdateWriteRepository(repo *RepositoryModel) (*RepositoryModel, error)
	DeleteWriteRepository(repo string, opts *RepositoryQueryOptions) error
	ValidateWriteAccess(opts *RepoAccessQuery) error
}

type RepositoryModel struct {
	TypeMeta
	ObjectMeta
	Repo                 string          `json:"repo"`
	Username             string          `json:"username,omitempty"`
	Password             string          `json:"password,omitempty"`
	SSHPrivateKey        string          `json:"sshPrivateKey,omitempty"`
	Insecure             bool            `json:"insecure,omitempty"`
	EnableLFS            bool            `json:"enableLfs,omitempty"`
	EnableOCI            bool            `json:"enableOCI,omitempty"`
	Type                 string          `json:"type,omitempty"`
	Name                 string          `json:"name,omitempty"`
	Project              string          `json:"project,omitempty"`
	InheritedCreds       bool            `json:"inheritedCreds,omitempty"`
	TLSClientCertData    string          `json:"tlsClientCertData,omitempty"`
	TLSClientCertKey     string          `json:"tlsClientCertKey,omitempty"`
	GCPServiceAccountKey string          `json:"gcpServiceAccountKey,omitempty"`
	Proxy                string          `json:"proxy,omitempty"`
	ForceHTTPBasicAuth   bool            `json:"forceHttpBasicAuth,omitempty"`
	ConnectionState      ConnectionState `json:"connectionState,omitempty"`
}

type RepositoryList struct {
	Items []*RepositoryModel `json:"items"`
}

type RepositoryAppsList struct {
	Items []*RepoApp `json:"items"`
}

type RepoApp struct {
	RepoURL     string `json:"repoUrl"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	AppName     string `json:"appName,omitempty"`
	ClusterName string `json:"clusterName,omitempty"`
	NameSpace   string `json:"namespace,omitempty"`
}

type RepoAppDetailsQuery struct {
	Source      ApplicationSource `json:"source"`
	AppName     string            `json:"appName,omitempty"`
	AppProject  string            `json:"appProject,omitempty"`
	SourceIndex *int32            `json:"sourceIndex,omitempty"`
	VersionId   *int32            `json:"versionId,omitempty"`
}

type RepoAppDetails struct {
	Type       string            `json:"type"`
	Path       string            `json:"path"`
	Kustomize  *KustomizeAppSpec `json:"kustomize,omitempty"`
	Directory  *DirectoryAppSpec `json:"directory,omitempty"`
	Helm       *HelmAppSpec      `json:"helm,omitempty"`
	Plugin     *PluginAppSpec    `json:"plugin,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type KustomizeAppSpec struct {
	Path   string   `json:"path,omitempty"`
	Images []string `json:"images,omitempty"`
}

type DirectoryAppSpec struct {
}

type HelmAppSpec struct {
	Name           string              `json:"name,omitempty"`
	ValueFiles     []string            `json:"valueFiles,omitempty"`
	Parameters     []HelmParameter     `json:"parameters,omitempty"`
	FileParameters []HelmFileParameter `json:"fileParameters,omitempty"`
	Values         string              `json:"values,omitempty"`
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
	Repo                              string `json:"repo"`
	Username                          string `json:"username,omitempty"`
	Password                          string `json:"password,omitempty"`
	SSHPrivateKey                     string `json:"sshPrivateKey,omitempty"`
	TLSClientCertData                 string `json:"tlsClientCertData,omitempty"`
	TLSClientCertKey                  string `json:"tlsClientCertKey,omitempty"`
	Type                              string `json:"type,omitempty"`
	Name                              string `json:"name,omitempty"`
	Insecure                          bool   `json:"insecure,omitempty"`
	EnableOCI                         bool   `json:"enableOCI,omitempty"`
	Proxy                             string `json:"proxy,omitempty"`
	ForceHTTPBasicAuth                bool   `json:"forceHttpBasicAuth,omitempty"`
	GitHubAppPrivateKey               string `json:"githubAppPrivateKey,omitempty"`
	GitHubAppID                       int64  `json:"githubAppID,omitempty"`
	GitHubAppInstallationID           int64  `json:"githubAppInstallationID,omitempty"`
	GitHubAppEnterpriseBaseUrl        string `json:"githubAppEnterpriseBaseUrl,omitempty"`
	GCPServiceAccountKey              string `json:"gcpServiceAccountKey,omitempty"`
	BearerToken                       string `json:"bearerToken,omitempty"`
	InsecureOCIForceHttp              bool   `json:"insecureOCIForceHttp,omitempty"`
	AzureServicePrincipalClientId     string `json:"azureServicePrincipalClientId,omitempty"`
	AzureServicePrincipalClientSecret string `json:"azureServicePrincipalClientSecret,omitempty"`
	AzureServicePrincipalTenantId     string `json:"azureServicePrincipalTenantId,omitempty"`
	AzureActiveDirectoryEndpoint      string `json:"azureActiveDirectoryEndpoint,omitempty"`
}

type RepositoryQueryOptions struct {
	Repo         string `json:"repo,omitempty"`
	ForceRefresh bool   `json:"forceRefresh,omitempty"`
	AppProject   string `json:"appProject,omitempty"`
}

type RepositoryCreateOptions struct {
	Upsert    bool `json:"upsert,omitempty"`
	CredsOnly bool `json:"credsOnly,omitempty"`
}

type RepoAppsQueryOptions struct {
	Revision   string `json:"revision,omitempty"`
	AppName    string `json:"appName,omitempty"`
	AppProject string `json:"appProject,omitempty"`
}

type RepositoryStandard struct {
	client *resty.Client
}

func NewRepository(client *resty.Client) Repository {
	return &RepositoryStandard{client: client}
}

func addRepositoryQueryOptions(req *resty.Request, opts *RepositoryQueryOptions) *resty.Request {
	if opts == nil {
		return req
	}
	if opts.Repo != "" {
		req.SetQueryParam("repo", opts.Repo)
	}
	if opts.ForceRefresh {
		req.SetQueryParam("forceRefresh", "true")
	}
	if opts.AppProject != "" {
		req.SetQueryParam("appProject", opts.AppProject)
	}
	return req
}

func addRepoAppsQueryOptions(req *resty.Request, opts *RepoAppsQueryOptions) *resty.Request {
	if opts == nil {
		return req
	}
	if opts.Revision != "" {
		req.SetQueryParam("revision", opts.Revision)
	}
	if opts.AppName != "" {
		req.SetQueryParam("appName", opts.AppName)
	}
	if opts.AppProject != "" {
		req.SetQueryParam("appProject", opts.AppProject)
	}
	return req
}

func addRepositoryCreateOptions(req *resty.Request, opts *RepositoryCreateOptions) *resty.Request {
	if opts == nil {
		return req
	}
	req.SetQueryParam("upsert", fmt.Sprintf("%v", opts.Upsert))
	req.SetQueryParam("credsOnly", fmt.Sprintf("%v", opts.CredsOnly))
	return req
}

func (r *RepositoryStandard) List(opts *RepositoryQueryOptions) (*RepositoryList, error) {
	var result RepositoryList
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get("/api/v1/repositories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) Get(repo string, opts *RepositoryQueryOptions) (*RepositoryModel, error) {
	var result RepositoryModel
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get(fmt.Sprintf("/api/v1/repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) Create(repo *RepositoryModel, opts *RepositoryCreateOptions) (*RepositoryModel, error) {
	var result RepositoryModel
	req := r.client.R().
		SetBody(repo).
		SetResult(&result)
	req = addRepositoryCreateOptions(req, opts)
	resp, err := req.Post("/api/v1/repositories")
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

func (r *RepositoryStandard) Delete(repo string, opts *RepositoryQueryOptions) error {
	req := r.client.R()
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Delete(fmt.Sprintf("/api/v1/repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (r *RepositoryStandard) ListApps(repo string, opts *RepoAppsQueryOptions) (*RepositoryAppsList, error) {
	var result RepositoryAppsList
	req := r.client.R().
		SetResult(&result)
	req = addRepoAppsQueryOptions(req, opts)
	resp, err := req.Get(fmt.Sprintf("/api/v1/repositories/%s/apps", url.PathEscape(repo)))
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

func (r *RepositoryStandard) GetHelmCharts(repo string, opts *RepositoryQueryOptions) (*HelmChartsResponse, error) {
	var result HelmChartsResponse
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get(fmt.Sprintf("/api/v1/repositories/%s/helmcharts", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) ListRefs(repo string, opts *RepositoryQueryOptions) (*RefsResponse, error) {
	var result RefsResponse
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get(fmt.Sprintf("/api/v1/repositories/%s/refs", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) ListOCITags(repo string, opts *RepositoryQueryOptions) (*OCITagsResponse, error) {
	var result OCITagsResponse
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get(fmt.Sprintf("/api/v1/repositories/%s/oci/tags", url.PathEscape(repo)))
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

func (r *RepositoryStandard) ListWriteRepositories(opts *RepositoryQueryOptions) (*RepositoryList, error) {
	var result RepositoryList
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get("/api/v1/write-repositories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) GetWrite(repo string, opts *RepositoryQueryOptions) (*RepositoryModel, error) {
	var result RepositoryModel
	req := r.client.R().
		SetResult(&result)
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Get(fmt.Sprintf("/api/v1/write-repositories/%s", url.PathEscape(repo)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (r *RepositoryStandard) CreateWriteRepository(repo *RepositoryModel, opts *RepositoryCreateOptions) (*RepositoryModel, error) {
	var result RepositoryModel
	req := r.client.R().
		SetBody(repo).
		SetResult(&result)
	req = addRepositoryCreateOptions(req, opts)
	resp, err := req.Post("/api/v1/write-repositories")
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

func (r *RepositoryStandard) DeleteWriteRepository(repo string, opts *RepositoryQueryOptions) error {
	req := r.client.R()
	req = addRepositoryQueryOptions(req, opts)
	resp, err := req.Delete(fmt.Sprintf("/api/v1/write-repositories/%s", url.PathEscape(repo)))
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
