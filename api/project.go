package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Project interface {
	List() (*ProjectList, error)
	Get(name string) (*ProjectModel, error)
	Create(project *ProjectModel, opts *ProjectCreateOptions) (*ProjectModel, error)
	Update(project *ProjectModel) (*ProjectModel, error)
	Delete(name string) error
	GetDetailed(name string) (*ProjectDetailed, error)
	GetGlobalProjects(name string) (*ProjectGlobalResponse, error)
	CreateToken(project, role string, opts *ProjectTokenCreateOptions) (*TokenResponse, error)
	DeleteToken(project, role string, iat int64, id string) error
	ListEvents(name string) (*ResourceEventList, error)
	GetSyncWindowsState(name string) (*SyncWindows, error)
	ListLinks(name string) (*LinksResponse, error)
}

type ProjectModel struct {
	ObjectMeta
	Spec   ProjectSpec   `json:"spec"`
	Status ProjectStatus `json:"status,omitempty"`
}

type ProjectSpec struct {
	SourceRepos                     []string                          `json:"sourceRepos,omitempty"`
	Destinations                    []ApplicationDestination          `json:"destinations,omitempty"`
	Description                     string                            `json:"description,omitempty"`
	Roles                           []ProjectRole                     `json:"roles,omitempty"`
	ClusterResourceWhitelist        []GroupKind                       `json:"clusterResourceWhitelist,omitempty"`
	NamespaceResourceBlacklist      []GroupKind                       `json:"namespaceResourceBlacklist,omitempty"`
	OrphanedResources               *OrphanedResourcesMonitorSettings `json:"orphanedResources,omitempty"`
	SyncWindows                     SyncWindows                       `json:"syncWindows,omitempty"`
	SignatureKeys                   []SignatureKey                    `json:"signatureKeys,omitempty"`
	ClusterResourceBlacklist        []GroupKind                       `json:"clusterResourceBlacklist,omitempty"`
	NamespaceResourceWhitelist      []GroupKind                       `json:"namespaceResourceWhitelist,omitempty"`
	SourceNamespaces                []string                          `json:"sourceNamespaces,omitempty"`
	PermitOnlyProjectScopedClusters bool                              `json:"permitOnlyProjectScopedClusters,omitempty"`
}

type GroupKind struct {
	Group string `json:"group,omitempty"`
	Kind  string `json:"kind"`
}

type ProjectRole struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Policies    []string   `json:"policies,omitempty"`
	JWTTokens   []JWTToken `json:"jwtTokens,omitempty"`
	Groups      []string   `json:"groups,omitempty"`
}

type JWTToken struct {
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp,omitempty"`
	ID        string `json:"id,omitempty"`
}

type OrphanedResourcesMonitorSettings struct {
	Warn *bool `json:"warn,omitempty"`
}

type SignatureKey struct {
	KeyID string `json:"keyID"`
}

type ProjectStatus struct {
	JWTTokensByRole map[string]JWTToken `json:"jwtTokensByRole,omitempty"`
}

type ProjectList struct {
	ListMeta
	Items []*ProjectModel `json:"items"`
}

type ProjectDetailed struct {
	ProjectModel
	GlobalProjects []*ProjectModel    `json:"globalProjects,omitempty"`
	Repositories   []*RepositoryModel `json:"repositories,omitempty"`
	Clusters       []*ClusterModel    `json:"clusters,omitempty"`
}

type ProjectGlobalResponse struct {
	ProjectModel
}

type ProjectTokenCreateOptions struct {
	ID          string `json:"id"`
	ExpiresIn   int64  `json:"expiresIn"`
	Description string `json:"description,omitempty"`
}

type ProjectCreateOptions struct {
	Upsert bool `json:"upsert,omitempty"`
}

type ProjectStandard struct {
	client *resty.Client
}

func NewProject(client *resty.Client) Project {
	return &ProjectStandard{client: client}
}

func (p *ProjectStandard) List() (*ProjectList, error) {
	var result ProjectList
	resp, err := p.client.R().
		SetResult(&result).
		Get("/api/v1/projects")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) Get(name string) (*ProjectModel, error) {
	var result ProjectModel
	resp, err := p.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/projects/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) Create(project *ProjectModel, opts *ProjectCreateOptions) (*ProjectModel, error) {
	var result ProjectModel
	req := p.client.R().
		SetBody(project).
		SetResult(&result)
	if opts != nil && opts.Upsert {
		req.SetQueryParam("upsert", "true")
	}
	resp, err := req.Post("/api/v1/projects")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) Update(project *ProjectModel) (*ProjectModel, error) {
	var result ProjectModel
	resp, err := p.client.R().
		SetBody(project).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/projects/%s", project.Name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) Delete(name string) error {
	resp, err := p.client.R().
		Delete(fmt.Sprintf("/api/v1/projects/%s", name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (p *ProjectStandard) GetDetailed(name string) (*ProjectDetailed, error) {
	var result ProjectDetailed
	resp, err := p.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/projects/%s/detailed", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) GetGlobalProjects(name string) (*ProjectGlobalResponse, error) {
	var result ProjectGlobalResponse
	resp, err := p.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/projects/%s/globalprojects", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) CreateToken(project, role string, opts *ProjectTokenCreateOptions) (*TokenResponse, error) {
	var result TokenResponse
	resp, err := p.client.R().
		SetBody(opts).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/projects/%s/roles/%s/token", project, role))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) DeleteToken(project, role string, iat int64, id string) error {
	resp, err := p.client.R().
		Delete(fmt.Sprintf("/api/v1/projects/%s/roles/%s/token/%d/%s", project, role, iat, id))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (p *ProjectStandard) ListEvents(name string) (*ResourceEventList, error) {
	var result ResourceEventList
	resp, err := p.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/projects/%s/events", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) GetSyncWindowsState(name string) (*SyncWindows, error) {
	var result SyncWindows
	resp, err := p.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/projects/%s/syncwindows", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (p *ProjectStandard) ListLinks(name string) (*LinksResponse, error) {
	var result LinksResponse
	resp, err := p.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/projects/%s/links", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

var _ Project = (*ProjectStandard)(nil)
