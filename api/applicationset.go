package api

import (
	"context"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
)

type ApplicationSet interface {
	List() ([]*ApplicationSetModel, error)
	Get(name string) (*ApplicationSetModel, error)
	Create(appSet *ApplicationSetModel) (*ApplicationSetModel, error)
	Delete(name string) error
	Generate(appSet *ApplicationSetModel) ([]*ApplicationSetModel, error)
	ResourceTree(name string) (*ApplicationTree, error)
	ListResourceEvents(name string) (*ResourceEventList, error)
	Watch(ctx context.Context) (<-chan *ApplicationSetWatchEvent, error)
}

type ApplicationSetModel struct {
	TypeMeta
	ObjectMeta
	Spec   ApplicationSetSpec   `json:"spec"`
	Status ApplicationSetStatus `json:"status,omitempty"`
}

type ApplicationSetSpec struct {
	Generators   []ApplicationSetGenerator `json:"generators"`
	Template     ApplicationSetTemplate    `json:"template"`
	SyncPolicy   *ApplicationSetSyncPolicy `json:"syncPolicy,omitempty"`
	Strategy     *ApplicationSetStrategy   `json:"strategy,omitempty"`
	PreservedFields  []string               `json:"preservedFields,omitempty"`
	GoTemplate       *bool                 `json:"goTemplate,omitempty"`
	GoTemplateOptions []string             `json:"goTemplateOptions,omitempty"`
}

type ApplicationSetGenerator struct {
	List                    *ListGenerator           `json:"list,omitempty"`
	Clusters                *ClusterGenerator        `json:"clusters,omitempty"`
	Git                     *GitGenerator            `json:"git,omitempty"`
	SCMProvider             *SCMProviderGenerator    `json:"scmProvider,omitempty"`
	ClusterDecisionResource *DuckTypeGenerator       `json:"clusterDecisionResource,omitempty"`
	PullRequest             *PullRequestGenerator    `json:"pullRequest,omitempty"`
	Matrix                  *MatrixGenerator         `json:"matrix,omitempty"`
	Merge                   *MergeGenerator          `json:"merge,omitempty"`
	Selector                *Selector                `json:"selector,omitempty"`
	Plugin                  *PluginGenerator         `json:"plugin,omitempty"`
}

type ListGenerator struct {
	Elements    []map[string]string `json:"elements"`
	ElementsYAML string             `json:"elementsYaml,omitempty"`
	Template    ApplicationSetTemplate `json:"template,omitempty"`
}

type ClusterGenerator struct {
	Selector Selector              `json:"selector,omitempty"`
	Values   map[string]string     `json:"values,omitempty"`
	Template ApplicationSetTemplate `json:"template,omitempty"`
}

type GitGenerator struct {
	RepoURL             string                  `json:"repoURL"`
	Directories         []GitDirectoryGeneratorItem `json:"directories,omitempty"`
	Files               []GitFileGeneratorItem  `json:"files,omitempty"`
	Revision            string                  `json:"revision"`
	RequeueAfterSeconds *int64                  `json:"requeueAfterSeconds,omitempty"`
	Template            ApplicationSetTemplate  `json:"template,omitempty"`
	PathParamPrefix     string                  `json:"pathParamPrefix,omitempty"`
	Values              map[string]string       `json:"values,omitempty"`
}

type GitDirectoryGeneratorItem struct {
	Path    string `json:"path"`
	Exclude *bool  `json:"exclude,omitempty"`
}

type GitFileGeneratorItem struct {
	Path string `json:"path"`
}

type SCMProviderGenerator struct {
	Github          *SCMProviderGeneratorGithub          `json:"github,omitempty"`
	Gitlab          *SCMProviderGeneratorGitlab          `json:"gitlab,omitempty"`
	Bitbucket       *SCMProviderGeneratorBitbucket       `json:"bitbucket,omitempty"`
	BitbucketServer *SCMProviderGeneratorBitbucketServer `json:"bitbucketServer,omitempty"`
	Gitea           *SCMProviderGeneratorGitea           `json:"gitea,omitempty"`
	AzureDevOps     *SCMProviderGeneratorAzureDevOps     `json:"azureDevOps,omitempty"`
	Filters         []SCMProviderGeneratorFilter         `json:"filters,omitempty"`
	CloneProtocol   string                               `json:"cloneProtocol,omitempty"`
	RequeueAfterSeconds *int64                           `json:"requeueAfterSeconds,omitempty"`
	Template        ApplicationSetTemplate               `json:"template,omitempty"`
	Values          map[string]string                    `json:"values,omitempty"`
}

type SCMProviderGeneratorGithub struct {
	Organization        string                 `json:"organization"`
	API                 string                 `json:"api,omitempty"`
	AppSecretName       string                 `json:"appSecretName,omitempty"`
	TokenRef            *SecretRef             `json:"tokenRef,omitempty"`
	AllBranches         *bool                  `json:"allBranches,omitempty"`
}

type SCMProviderGeneratorGitlab struct {
	Group               string                 `json:"group"`
	IncludeSubgroups    bool                   `json:"includeSubgroups,omitempty"`
	TokenRef            *SecretRef             `json:"tokenRef,omitempty"`
	AllBranches         *bool                  `json:"allBranches,omitempty"`
}

type SCMProviderGeneratorBitbucket struct {
	Owner               string                 `json:"owner"`
	User                string                 `json:"user"`
	AppPasswordRef      *SecretRef             `json:"appPasswordRef,omitempty"`
	AllBranches         *bool                  `json:"allBranches,omitempty"`
}

type SCMProviderGeneratorBitbucketServer struct {
	Project             string                 `json:"project"`
	API                 string                 `json:"api"`
	BasicAuthRef        *SecretRef             `json:"basicAuthRef,omitempty"`
	AllBranches         *bool                  `json:"allBranches,omitempty"`
}

type SCMProviderGeneratorGitea struct {
	Owner               string                 `json:"owner"`
	API                 string                 `json:"api"`
	TokenRef            *SecretRef             `json:"tokenRef,omitempty"`
	Insecure            bool                   `json:"insecure,omitempty"`
	AllBranches         *bool                  `json:"allBranches,omitempty"`
}

type SCMProviderGeneratorAzureDevOps struct {
	Organization        string                 `json:"organization"`
	TeamProject         string                 `json:"teamProject"`
	AccessTokenRef      *SecretRef             `json:"accessTokenRef,omitempty"`
	API                 string                 `json:"api,omitempty"`
	AllBranches         *bool                  `json:"allBranches,omitempty"`
}

type SCMProviderGeneratorFilter struct {
	RepositoryMatch     *string                `json:"repositoryMatch,omitempty"`
	PathsExist          []string               `json:"pathsExist,omitempty"`
	PathsDoNotExist     []string               `json:"pathsDoNotExist,omitempty"`
	LabelMatch          *string                `json:"labelMatch,omitempty"`
	BranchMatch         *string                `json:"branchMatch,omitempty"`
}

type SecretRef struct {
	SecretName string `json:"secretName"`
	Key        string `json:"key"`
}

type DuckTypeGenerator struct {
	ConfigMapRef        string                 `json:"configMapRef"`
	Name                string                 `json:"name,omitempty"`
	LabelSelector       string                 `json:"labelSelector,omitempty"`
	RequeueAfterSeconds *int64                 `json:"requeueAfterSeconds,omitempty"`
	Template            ApplicationSetTemplate `json:"template,omitempty"`
	Values              map[string]string      `json:"values,omitempty"`
}

type PullRequestGenerator struct {
	Github          *PullRequestGeneratorGithub          `json:"github,omitempty"`
	Gitlab          *PullRequestGeneratorGitlab          `json:"gitlab,omitempty"`
	BitbucketServer *PullRequestGeneratorBitbucketServer `json:"bitbucketServer,omitempty"`
	Gitea           *PullRequestGeneratorGitea           `json:"gitea,omitempty"`
	Bitbucket       *PullRequestGeneratorBitbucket       `json:"bitbucket,omitempty"`
	AzureDevOps     *PullRequestGeneratorAzureDevOps     `json:"azureDevOps,omitempty"`
	Filters         []PullRequestGeneratorFilter         `json:"filters,omitempty"`
	RequeueAfterSeconds *int64                           `json:"requeueAfterSeconds,omitempty"`
	Template        ApplicationSetTemplate               `json:"template,omitempty"`
	BitbucketServerBasicAuth *SecretRef                  `json:"bitbucketServerBasicAuth,omitempty"`
}

type PullRequestGeneratorGithub struct {
	Owner               string                 `json:"owner"`
	Repo                string                 `json:"repo"`
	API                 string                 `json:"api,omitempty"`
	AppSecretName       string                 `json:"appSecretName,omitempty"`
	TokenRef            *SecretRef             `json:"tokenRef,omitempty"`
	Labels              []string               `json:"labels,omitempty"`
}

type PullRequestGeneratorGitlab struct {
	Project             string                 `json:"project"`
	API                 string                 `json:"api,omitempty"`
	TokenRef            *SecretRef             `json:"tokenRef,omitempty"`
	Labels              []string               `json:"labels,omitempty"`
	PullRequestState    string                 `json:"pullRequestState,omitempty"`
}

type PullRequestGeneratorBitbucketServer struct {
	Project             string                 `json:"project"`
	Repo                string                 `json:"repo"`
	API                 string                 `json:"api"`
}

type PullRequestGeneratorGitea struct {
	Owner               string                 `json:"owner"`
	Repo                string                 `json:"repo"`
	API                 string                 `json:"api"`
	TokenRef            *SecretRef             `json:"tokenRef,omitempty"`
	Insecure            bool                   `json:"insecure,omitempty"`
}

type PullRequestGeneratorBitbucket struct {
	Owner               string                 `json:"owner"`
	Repo                string                 `json:"repo"`
	BearerTokenRef      *SecretRef             `json:"bearerTokenRef,omitempty"`
}

type PullRequestGeneratorAzureDevOps struct {
	Organization        string                 `json:"organization"`
	Project             string                 `json:"project"`
	Repo                string                 `json:"repo"`
	AccessTokenRef      *SecretRef             `json:"accessTokenRef,omitempty"`
	API                 string                 `json:"api,omitempty"`
	Labels              []string               `json:"labels,omitempty"`
}

type PullRequestGeneratorFilter struct {
	BranchMatch       *string  `json:"branchMatch,omitempty"`
	TargetBranchMatch *string  `json:"targetBranchMatch,omitempty"`
}

type MatrixGenerator struct {
	Generators []ApplicationSetGenerator `json:"generators"`
	Template   ApplicationSetTemplate    `json:"template,omitempty"`
}

type MergeGenerator struct {
	Generators      []ApplicationSetGenerator `json:"generators"`
	MergeKeys       []string                 `json:"mergeKeys"`
	Template        ApplicationSetTemplate    `json:"template,omitempty"`
}

type PluginGenerator struct {
	ConfigMapRef string            `json:"configMapRef"`
	Input        ApplicationSetTemplate `json:"input,omitempty"`
	Values       map[string]string `json:"values,omitempty"`
	RequeueAfterSeconds *int64     `json:"requeueAfterSeconds,omitempty"`
	Template     ApplicationSetTemplate `json:"template,omitempty"`
}

type Selector struct {
	MatchExpressions []SelectorMatchExpression `json:"matchExpressions,omitempty"`
	MatchLabels      map[string]string         `json:"matchLabels,omitempty"`
}

type SelectorMatchExpression struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values,omitempty"`
}

type ApplicationSetTemplate struct {
	ApplicationSetTemplateMeta `json:"metadata"`
	Spec                       ApplicationSpec `json:"spec"`
}

type ApplicationSetTemplateMeta struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Finalizers  []string          `json:"finalizers,omitempty"`
}

type ApplicationSetSyncPolicy struct {
	PreserveResourcesOnDeletion bool `json:"preserveResourcesOnDeletion,omitempty"`
}

type ApplicationSetStrategy struct {
	Type       string                              `json:"type,omitempty"`
	RollingSync *ApplicationSetRolloutStep         `json:"rollingSync,omitempty"`
	RollingUpdate *ApplicationSetRolloutStep       `json:"rollingUpdate,omitempty"`
}

type ApplicationSetRolloutStep struct {
	MaxUpdate *IntOrString `json:"maxUpdate,omitempty"`
}

type IntOrString struct {
	Type   int64  `json:"type,omitempty"`
	IntVal int64  `json:"intVal,omitempty"`
	StrVal string `json:"strVal,omitempty"`
}

type ApplicationSetStatus struct {
	Conditions       []ApplicationSetCondition `json:"conditions,omitempty"`
	Resources        []ResourceStatus          `json:"resources,omitempty"`
	ApplicationStatus []ApplicationStatus      `json:"applicationStatus,omitempty"`
}

type ApplicationSetCondition struct {
	Type               string `json:"type"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
}

type ApplicationSetWatchEvent struct {
	Type           SyncStatusCode        `json:"type"`
	ApplicationSet *ApplicationSetModel  `json:"applicationSet"`
	Application    *ApplicationModel     `json:"application,omitempty"`
	Apps           []*ApplicationModel   `json:"apps,omitempty"`
}

type ApplicationSetStandard struct {
	client *resty.Client
}

func NewApplicationSet(client *resty.Client) ApplicationSet {
	return &ApplicationSetStandard{client: client}
}

func (a *ApplicationSetStandard) List() ([]*ApplicationSetModel, error) {
	var result struct {
		Items []*ApplicationSetModel `json:"items"`
	}
	resp, err := a.client.R().
		SetResult(&result).
		Get("/api/v1/applicationsets")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return result.Items, nil
}

func (a *ApplicationSetStandard) Get(name string) (*ApplicationSetModel, error) {
	var result ApplicationSetModel
	resp, err := a.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/applicationsets/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationSetStandard) Create(appSet *ApplicationSetModel) (*ApplicationSetModel, error) {
	var result ApplicationSetModel
	resp, err := a.client.R().
		SetBody(appSet).
		SetResult(&result).
		Post("/api/v1/applicationsets")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationSetStandard) Delete(name string) error {
	resp, err := a.client.R().
		Delete(fmt.Sprintf("/api/v1/applicationsets/%s", name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationSetStandard) Generate(appSet *ApplicationSetModel) ([]*ApplicationSetModel, error) {
	var result struct {
		Items []*ApplicationSetModel `json:"items"`
	}
	resp, err := a.client.R().
		SetBody(appSet).
		SetResult(&result).
		Post("/api/v1/applicationsets/generate")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return result.Items, nil
}

func (a *ApplicationSetStandard) ResourceTree(name string) (*ApplicationTree, error) {
	var result ApplicationTree
	resp, err := a.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/applicationsets/%s/resource-tree", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationSetStandard) ListResourceEvents(name string) (*ResourceEventList, error) {
	var result ResourceEventList
	resp, err := a.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/applicationsets/%s/events", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationSetStandard) Watch(ctx context.Context) (<-chan *ApplicationSetWatchEvent, error) {
	ch := make(chan *ApplicationSetWatchEvent)
	resp, err := a.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Accept", "text/event-stream").
		Get("/api/v1/stream/applicationsets")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		body, _ := io.ReadAll(resp.RawBody())
		resp.RawBody().Close()
		return nil, ParseErrorFromBody(resp.StatusCode(), body)
	}
	go func() {
		defer close(ch)
		defer resp.RawBody().Close()
		readSSE(ctx, ch, resp.RawBody())
	}()
	return ch, nil
}

var _ ApplicationSet = (*ApplicationSetStandard)(nil)
