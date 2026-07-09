package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Application interface {
	List(opts *ApplicationListOptions) (*ApplicationList, error)
	Get(name string, opts *ApplicationGetOptions) (*ApplicationModel, error)
	Create(app *ApplicationModel, opts *ApplicationCreateOptions) (*ApplicationModel, error)
	Update(app *ApplicationModel, opts *ApplicationUpdateOptions) (*ApplicationModel, error)
	Delete(name string, opts *ApplicationDeleteOptions) error
	Sync(name string, opts *SyncOptions) error
	Rollback(name string, opts *RollbackOptions) error
	TerminateOperation(name string, opts *TerminateOperationOptions) error
	Patch(name string, patch any, patchType string, opts *ApplicationPatchOptions) (*ApplicationModel, error)
	GetResource(opts *ApplicationResourceRequest) (*ApplicationResourceResponse, error)
	PatchResource(opts *ApplicationResourcePatchRequest) (*ApplicationResourceResponse, error)
	DeleteResource(opts *ApplicationResourceDeleteRequest) error
	ListResourceActions(opts *ApplicationResourceRequest) (*ResourceActionsList, error)
	RunResourceAction(opts *ApplicationResourceActionRequest) error
	GetManifests(name string, opts *ApplicationManifestQuery) (*ManifestResponse, error)
	ResourceTree(name string, opts *ResourcesQuery) (*ApplicationTree, error)
	ManagedResources(name string, opts *ResourcesQuery) (*ManagedResourcesResponse, error)
	RevisionMetadata(name, revision string, opts *RevisionMetadataOptions) (*RevisionMetadata, error)
	GetSyncWindows(name string, opts *SyncWindowsOptions) (*SyncWindows, error)
	ListResourceEvents(name string, opts *ApplicationResourceEventsQuery) (*ResourceEventList, error)
	ListLinks(name string, opts *ListLinksOptions) (*LinksResponse, error)
	ListResourceLinks(opts *ApplicationResourceRequest) (*LinksResponse, error)
	Watch(ctx context.Context, opts *WatchOptions) (<-chan *ApplicationWatchEvent, error)
	WatchResourceTree(ctx context.Context, name string, opts *ResourcesQuery) (<-chan *ApplicationTree, error)
	PodLogs(ctx context.Context, opts *PodLogsOptions) (<-chan *LogEntry, error)
}

type ApplicationModel struct {
	TypeMeta
	ObjectMeta
	Spec      ApplicationSpec   `json:"spec"`
	Status    ApplicationStatus `json:"status,omitempty"`
	Operation *Operation        `json:"operation,omitempty"`
}

type ApplicationSpec struct {
	Source               *ApplicationSource          `json:"source,omitempty"`
	Sources              []ApplicationSource         `json:"sources,omitempty"`
	Destination          ApplicationDestination      `json:"destination"`
	Project              string                      `json:"project"`
	SyncPolicy           *SyncPolicy                 `json:"syncPolicy,omitempty"`
	IgnoreDifferences    []ResourceIgnoreDifferences `json:"ignoreDifferences,omitempty"`
	Info                 []Info                      `json:"info,omitempty"`
	RevisionHistoryLimit *int64                      `json:"revisionHistoryLimit,omitempty"`
}

type ApplicationSource struct {
	RepoURL        string                      `json:"repoURL"`
	Path           string                      `json:"path,omitempty"`
	TargetRevision string                      `json:"targetRevision,omitempty"`
	Helm           *ApplicationSourceHelm      `json:"helm,omitempty"`
	Kustomize      *ApplicationSourceKustomize `json:"kustomize,omitempty"`
	Directory      *ApplicationSourceDirectory `json:"directory,omitempty"`
	Plugin         *ApplicationSourcePlugin    `json:"plugin,omitempty"`
	Chart          string                      `json:"chart,omitempty"`
	Ref            string                      `json:"ref,omitempty"`
}

type ApplicationSourceHelm struct {
	ValueFiles     []string            `json:"valueFiles,omitempty"`
	Values         string              `json:"values,omitempty"`
	ValuesObject   any                 `json:"valuesObject,omitempty"`
	ReleaseName    string              `json:"releaseName,omitempty"`
	Parameters     []HelmParameter     `json:"parameters,omitempty"`
	FileParameters []HelmFileParameter `json:"fileParameters,omitempty"`
}

type HelmParameter struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	ForceString bool   `json:"forceString,omitempty"`
}

type HelmFileParameter struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ApplicationSourceKustomize struct {
	NamePrefix        string            `json:"namePrefix,omitempty"`
	NameSuffix        string            `json:"nameSuffix,omitempty"`
	Images            []KustomizeImage  `json:"images,omitempty"`
	CommonLabels      map[string]string `json:"commonLabels,omitempty"`
	Version           string            `json:"version,omitempty"`
	CommonAnnotations map[string]string `json:"commonAnnotations,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
}

type KustomizeImage string

type ApplicationSourceDirectory struct {
	Recurse bool                      `json:"recurse,omitempty"`
	JSONnet *ApplicationSourceJSONnet `json:"jsonnet,omitempty"`
	Exclude string                    `json:"exclude,omitempty"`
	Include string                    `json:"include,omitempty"`
}

type ApplicationSourceJSONnet struct {
	ExtVars []JSONnetVar `json:"extVars,omitempty"`
	TLAs    []JSONnetVar `json:"tlas,omitempty"`
	Libs    []string     `json:"libs,omitempty"`
}

type JSONnetVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Code  bool   `json:"code,omitempty"`
}

type ApplicationSourcePlugin struct {
	Name       string           `json:"name,omitempty"`
	Env        []EnvEntry       `json:"env,omitempty"`
	Parameters []ParameterEntry `json:"parameters,omitempty"`
}

type EnvEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ParameterEntry struct {
	Name    string            `json:"name"`
	Array   []string          `json:"array,omitempty"`
	Map     map[string]string `json:"map,omitempty"`
	String_ string            `json:"string,omitempty"`
}

type ApplicationDestination struct {
	Server    string `json:"server,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

type SyncPolicy struct {
	Automated                *SyncPolicyAutomated      `json:"automated,omitempty"`
	SyncOptions              []string                  `json:"syncOptions,omitempty"`
	Retry                    *RetryStrategy            `json:"retry,omitempty"`
	ManagedNamespaceMetadata *ManagedNamespaceMetadata `json:"managedNamespaceMetadata,omitempty"`
}

type SyncPolicyAutomated struct {
	Prune      bool `json:"prune,omitempty"`
	SelfHeal   bool `json:"selfHeal,omitempty"`
	AllowEmpty bool `json:"allowEmpty,omitempty"`
}

type RetryStrategy struct {
	Limit   int64    `json:"limit"`
	Backoff *Backoff `json:"backoff,omitempty"`
}

type Backoff struct {
	Duration    string `json:"duration,omitempty"`
	Factor      *int64 `json:"factor,omitempty"`
	MaxDuration string `json:"maxDuration,omitempty"`
}

type ManagedNamespaceMetadata struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ResourceIgnoreDifferences struct {
	Group                 string   `json:"group,omitempty"`
	Kind                  string   `json:"kind"`
	Name                  string   `json:"name,omitempty"`
	Namespace             string   `json:"namespace,omitempty"`
	JSONPointers          []string `json:"jsonPointers,omitempty"`
	JQPathExpressions     []string `json:"jqPathExpressions,omitempty"`
	ManagedFieldsManagers []string `json:"managedFieldsManagers,omitempty"`
}

type Info struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ApplicationStatus struct {
	Resources            []ResourceStatus       `json:"resources,omitempty"`
	Sync                 SyncStatus             `json:"sync,omitempty"`
	Health               HealthStatus           `json:"health,omitempty"`
	History              []RevisionHistory      `json:"history,omitempty"`
	Conditions           []ApplicationCondition `json:"conditions,omitempty"`
	ReconciledAt         string                 `json:"reconciledAt,omitempty"`
	OperationState       *OperationState        `json:"operationState,omitempty"`
	ObservedAt           string                 `json:"observedAt,omitempty"`
	SourceType           ApplicationSourceType  `json:"sourceType,omitempty"`
	Summary              *ApplicationSummary    `json:"summary,omitempty"`
	ResourceHealthSource string                 `json:"resourceHealthSource,omitempty"`
}

type ResourceStatus struct {
	Group           string         `json:"group,omitempty"`
	Version         string         `json:"version"`
	Kind            string         `json:"kind"`
	Namespace       string         `json:"namespace"`
	Name            string         `json:"name"`
	Status          SyncStatusCode `json:"status,omitempty"`
	Health          *HealthStatus  `json:"health,omitempty"`
	Hook            bool           `json:"hook,omitempty"`
	RequiresPruning bool           `json:"requiresPruning,omitempty"`
	SyncWave        int64          `json:"syncWave,omitempty"`
}

type SyncStatus struct {
	Status     SyncStatusCode `json:"status"`
	ComparedTo ComparedTo     `json:"comparedTo,omitempty"`
	Revision   string         `json:"revision,omitempty"`
	Revisions  []string       `json:"revisions,omitempty"`
}

type SyncStatusCode string

const (
	SyncStatusCodeUnknown   SyncStatusCode = "Unknown"
	SyncStatusCodeSynced    SyncStatusCode = "Synced"
	SyncStatusCodeOutOfSync SyncStatusCode = "OutOfSync"
)

type ComparedTo struct {
	Source      ApplicationSource      `json:"source,omitempty"`
	Sources     []ApplicationSource    `json:"sources,omitempty"`
	Destination ApplicationDestination `json:"destination"`
}

type HealthStatus struct {
	Status  HealthStatusCode `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
}

type HealthStatusCode string

const (
	HealthStatusUnknown     HealthStatusCode = "Unknown"
	HealthStatusProgressing HealthStatusCode = "Progressing"
	HealthStatusHealthy     HealthStatusCode = "Healthy"
	HealthStatusSuspended   HealthStatusCode = "Suspended"
	HealthStatusDegraded    HealthStatusCode = "Degraded"
	HealthStatusMissing     HealthStatusCode = "Missing"
)

type RevisionHistory struct {
	Revision        string            `json:"revision"`
	DeployedAt      string            `json:"deployedAt"`
	ID              int64             `json:"id"`
	Source          ApplicationSource `json:"source,omitempty"`
	DeployStartedAt string            `json:"deployStartedAt,omitempty"`
}

type ApplicationCondition struct {
	Type               string `json:"type"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
}

type OperationState struct {
	Operation  Operation            `json:"operation"`
	Phase      OperationPhase       `json:"phase"`
	Message    string               `json:"message,omitempty"`
	SyncResult *SyncOperationResult `json:"syncResult,omitempty"`
	StartedAt  string               `json:"startedAt"`
	FinishedAt string               `json:"finishedAt,omitempty"`
}

type Operation struct {
	Sync        *SyncOperation     `json:"sync,omitempty"`
	InitiatedBy OperationInitiator `json:"initiatedBy,omitempty"`
	Info        []*Info            `json:"info,omitempty"`
}

type SyncOperation struct {
	Revision     string                  `json:"revision,omitempty"`
	Revisions    []string                `json:"revisions,omitempty"`
	Prune        bool                    `json:"prune,omitempty"`
	DryRun       bool                    `json:"dryRun,omitempty"`
	SyncStrategy *SyncStrategy           `json:"syncStrategy,omitempty"`
	Resources    []SyncOperationResource `json:"resources,omitempty"`
	Source       *ApplicationSource      `json:"source,omitempty"`
	Sources      []ApplicationSource     `json:"sources,omitempty"`
	Manifests    []string                `json:"manifests,omitempty"`
	SyncOptions  []string                `json:"syncOptions,omitempty"`
}

type SyncOperationResource struct {
	Group     string `json:"group,omitempty"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type SyncStrategy struct {
	Apply *SyncStrategyApply `json:"apply,omitempty"`
	Hook  *SyncStrategyHook  `json:"hook,omitempty"`
}

type SyncStrategyApply struct {
	Force bool `json:"force,omitempty"`
}

type SyncStrategyHook struct {
	SyncStrategyApply
}

type OperationInitiator struct {
	Username  string `json:"username,omitempty"`
	Automated bool   `json:"automated,omitempty"`
}

type OperationPhase string

const (
	OperationRunning     OperationPhase = "Running"
	OperationTerminating OperationPhase = "Terminating"
	OperationFailed      OperationPhase = "Failed"
	OperationError       OperationPhase = "Error"
	OperationSucceeded   OperationPhase = "Succeeded"
)

type SyncOperationResult struct {
	Resources []*ResourceResult `json:"resources,omitempty"`
	Revision  string            `json:"revision"`
	Revisions []string          `json:"revisions,omitempty"`
	Source    ApplicationSource `json:"source,omitempty"`
}

type ResourceResult struct {
	Group     string         `json:"group"`
	Version   string         `json:"version"`
	Namespace string         `json:"namespace"`
	Name      string         `json:"name"`
	Status    ResultCode     `json:"status,omitempty"`
	Message   string         `json:"message,omitempty"`
	HookType  string         `json:"hookType,omitempty"`
	HookPhase OperationPhase `json:"hookPhase,omitempty"`
	SyncPhase string         `json:"syncPhase,omitempty"`
}

type ResultCode string

const (
	ResultCodeSynced       ResultCode = "Synced"
	ResultCodeSyncFailed   ResultCode = "SyncFailed"
	ResultCodePruned       ResultCode = "Pruned"
	ResultCodePruneSkipped ResultCode = "PruneSkipped"
)

type ApplicationSourceType string

const (
	ApplicationSourceTypeHelm      ApplicationSourceType = "Helm"
	ApplicationSourceTypeKustomize ApplicationSourceType = "Kustomize"
	ApplicationSourceTypeDirectory ApplicationSourceType = "Directory"
	ApplicationSourceTypePlugin    ApplicationSourceType = "Plugin"
)

type ApplicationSummary struct {
	ExternalURLs []string `json:"externalURLs,omitempty"`
	Images       []string `json:"images,omitempty"`
}

type ApplicationList struct {
	ListMeta
	Items []*ApplicationModel `json:"items"`
}

type ApplicationDeleteOptions struct {
	Cascade           *bool  `json:"cascade,omitempty"`
	PropagationPolicy string `json:"propagationPolicy,omitempty"`
	AppNamespace      string `json:"appNamespace,omitempty"`
	Project           string `json:"project,omitempty"`
}

type SyncOptions struct {
	Revision        string                  `json:"revision,omitempty"`
	Prune           bool                    `json:"prune,omitempty"`
	DryRun          bool                    `json:"dryRun,omitempty"`
	Strategy        *SyncStrategy           `json:"strategy,omitempty"`
	Resources       []SyncOperationResource `json:"resources,omitempty"`
	Manifests       []string                `json:"manifests,omitempty"`
	Infos           []Info                  `json:"infos,omitempty"`
	RetryStrategy   *RetryStrategy          `json:"retryStrategy,omitempty"`
	SyncOptions     []string                `json:"syncOptions,omitempty"`
	AppNamespace    string                  `json:"appNamespace,omitempty"`
	Project         string                  `json:"project,omitempty"`
	SourcePositions []int64                 `json:"sourcePositions,omitempty"`
	Revisions       []string                `json:"revisions,omitempty"`
}

type RollbackOptions struct {
	ID           int64  `json:"id"`
	DryRun       bool   `json:"dryRun,omitempty"`
	Prune        bool   `json:"prune,omitempty"`
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
}

type ApplicationResourceRequest struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	ResourceName string `json:"resourceName"`
	Version      string `json:"version"`
	Group        string `json:"group,omitempty"`
	Kind         string `json:"kind"`
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
}

type ApplicationResourceResponse struct {
	Manifest string `json:"manifest,omitempty"`
}

type ApplicationResourcePatchRequest struct {
	ApplicationResourceRequest
	Patch     string `json:"patch"`
	PatchType string `json:"patchType"`
}

type ApplicationResourceDeleteRequest struct {
	ApplicationResourceRequest
	Force  *bool `json:"force,omitempty"`
	Orphan *bool `json:"orphan,omitempty"`
}

type ResourceActionsList struct {
	Actions []ResourceAction `json:"actions"`
}

type ResourceAction struct {
	Name     string                `json:"name"`
	Params   []ResourceActionParam `json:"params,omitempty"`
	Disabled bool                  `json:"disabled,omitempty"`
}

type ResourceActionParam struct {
	Name    string `json:"name"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Default string `json:"default,omitempty"`
}

type ApplicationResourceActionRequest struct {
	ApplicationResourceRequest
	Action string `json:"action"`
}

type ApplicationManifestQuery struct {
	Revision        string   `json:"revision,omitempty"`
	AppNamespace    string   `json:"namespace,omitempty"`
	Project         string   `json:"project,omitempty"`
	SourcePositions []int64  `json:"sourcePositions,omitempty"`
	Revisions       []string `json:"revisions,omitempty"`
	NoCache         *bool    `json:"noCache,omitempty"`
}

type ManifestResponse struct {
	Manifests []string `json:"manifests"`
}

type ApplicationTree struct {
	Nodes             []ResourceNode `json:"nodes,omitempty"`
	OrphanedResources []ResourceNode `json:"orphanedResources,omitempty"`
	Hosts             []string       `json:"hosts,omitempty"`
}

type ResourceNode struct {
	Group           string                  `json:"group,omitempty"`
	Version         string                  `json:"version"`
	Kind            string                  `json:"kind"`
	Namespace       string                  `json:"namespace"`
	Name            string                  `json:"name"`
	UID             string                  `json:"uid,omitempty"`
	ResourceVersion string                  `json:"resourceVersion,omitempty"`
	Status          SyncStatusCode          `json:"status,omitempty"`
	Health          *HealthStatus           `json:"health,omitempty"`
	CreatedAt       string                  `json:"createdAt,omitempty"`
	Images          []string                `json:"images,omitempty"`
	NetworkingInfo  *ResourceNetworkingInfo `json:"networkingInfo,omitempty"`
	Info            []InfoItem              `json:"info,omitempty"`
	ParentRefs      []ResourceRef           `json:"parentRefs,omitempty"`
}

type ResourceNetworkingInfo struct {
	TargetLabels map[string]string     `json:"targetLabels,omitempty"`
	TargetRefs   []ResourceRef         `json:"targetRefs,omitempty"`
	Labels       map[string]string     `json:"labels,omitempty"`
	Ingress      []LoadBalancerIngress `json:"ingress,omitempty"`
	ExternalURLs []string              `json:"externalURLs,omitempty"`
}

type ResourceRef struct {
	Group     string `json:"group,omitempty"`
	Version   string `json:"version,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	UID       string `json:"uid,omitempty"`
}

type InfoItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LoadBalancerIngress struct {
	Hostname string `json:"hostname,omitempty"`
	IP       string `json:"ip,omitempty"`
}

type ManagedResourcesResponse struct {
	Items []*ResourceDiff `json:"items"`
}

type ResourceDiff struct {
	Group               string `json:"group,omitempty"`
	Kind                string `json:"kind"`
	Namespace           string `json:"namespace"`
	Name                string `json:"name"`
	TargetState         string `json:"targetState,omitempty"`
	LiveState           string `json:"liveState,omitempty"`
	Diff                string `json:"diff,omitempty"`
	Hook                bool   `json:"hook,omitempty"`
	NormalizedLiveState string `json:"normalizedLiveState,omitempty"`
	PredictedLiveState  string `json:"predictedLiveState,omitempty"`
	ResourceVersion     string `json:"resourceVersion,omitempty"`
	Modified            bool   `json:"modified,omitempty"`
}

type RevisionMetadata struct {
	Author        string   `json:"author,omitempty"`
	Date          string   `json:"date"`
	Tags          []string `json:"tags,omitempty"`
	Message       string   `json:"message,omitempty"`
	SignatureInfo string   `json:"signatureInfo,omitempty"`
}

type SyncWindows struct {
	Windows []*SyncWindow `json:"windows,omitempty"`
}

type SyncWindow struct {
	Kind         string   `json:"kind,omitempty"`
	Schedule     string   `json:"schedule"`
	Duration     string   `json:"duration"`
	Applications []string `json:"applications,omitempty"`
	Namespaces   []string `json:"namespaces,omitempty"`
	Clusters     []string `json:"clusters,omitempty"`
	ManualSync   bool     `json:"manualSync"`
	TimeZone     string   `json:"timeZone,omitempty"`
}

type ApplicationResourceEventsQuery struct {
	Name              string `json:"name,omitempty"`
	Namespace         string `json:"namespace,omitempty"`
	ResourceName      string `json:"resourceName,omitempty"`
	ResourceUID       string `json:"resourceUID,omitempty"`
	ResourceNamespace string `json:"resourceNamespace,omitempty"`
	AppNamespace      string `json:"appNamespace,omitempty"`
	Project           string `json:"project,omitempty"`
}

type ResourceEventList struct {
	Items []ResourceEvent `json:"items"`
}

type ResourceEvent struct {
	Action         string            `json:"action"`
	EventTime      string            `json:"eventTime"`
	Reason         string            `json:"reason"`
	Note           string            `json:"note"`
	Type           string            `json:"type"`
	InvolvedObject ObjectReference   `json:"involvedObject"`
	Source         EventSource       `json:"source"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type ObjectReference struct {
	Kind            string `json:"kind,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
	Name            string `json:"name,omitempty"`
	UID             string `json:"uid,omitempty"`
	APIVersion      string `json:"apiVersion,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	FieldPath       string `json:"fieldPath,omitempty"`
}

type EventSource struct {
	Component string `json:"component,omitempty"`
	Host      string `json:"host,omitempty"`
}

type LinksResponse struct {
	Items []LinkItem `json:"items,omitempty"`
}

type LinkItem struct {
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
}

type WatchOptions struct {
	Revision        string   `json:"revision,omitempty"`
	Namespace       string   `json:"namespace,omitempty"`
	ResourceVersion string   `json:"resourceVersion,omitempty"`
	Selector        string   `json:"selector,omitempty"`
	Repo            string   `json:"repo,omitempty"`
	Project         []string `json:"project,omitempty"`
	Refresh         string   `json:"refresh,omitempty"`
}

type ApplicationWatchEvent struct {
	Type        SyncStatusCode    `json:"type"`
	Application *ApplicationModel `json:"application"`
}

type PodLogsOptions struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace,omitempty"`
	ResourceName string `json:"resourceName"`
	Group        string `json:"group,omitempty"`
	Kind         string `json:"kind"`
	Container    string `json:"container,omitempty"`
	TailLines    int64  `json:"tailLines,omitempty"`
	Follow       bool   `json:"follow,omitempty"`
	SinceSeconds int64  `json:"sinceSeconds,omitempty"`
	SinceTime    string `json:"sinceTime,omitempty"`
	Previous     bool   `json:"previous,omitempty"`
	PodName      string `json:"podName,omitempty"`
	UntilTime    string `json:"untilTime,omitempty"`
	Filter       string `json:"filter,omitempty"`
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
	MatchCase    *bool  `json:"matchCase,omitempty"`
}

type LogEntry struct {
	Content      string `json:"content"`
	TimeStamp    string `json:"timeStamp,omitempty"`
	Last         bool   `json:"last,omitempty"`
	TimeStampStr string `json:"timeStampStr,omitempty"`
	PodName      string `json:"podName,omitempty"`
}

type ApplicationListOptions struct {
	Refresh         string   `json:"refresh,omitempty"`
	ResourceVersion string   `json:"resourceVersion,omitempty"`
	Selector        string   `json:"selector,omitempty"`
	Repo            string   `json:"repo,omitempty"`
	AppNamespace    string   `json:"appNamespace,omitempty"`
	Project         []string `json:"project,omitempty"`
}

type ApplicationGetOptions struct {
	Refresh         string   `json:"refresh,omitempty"`
	ResourceVersion string   `json:"resourceVersion,omitempty"`
	Selector        string   `json:"selector,omitempty"`
	Repo            string   `json:"repo,omitempty"`
	AppNamespace    string   `json:"appNamespace,omitempty"`
	Project         []string `json:"project,omitempty"`
}

type ApplicationCreateOptions struct {
	Upsert   *bool `json:"upsert,omitempty"`
	Validate *bool `json:"validate,omitempty"`
}

type ApplicationUpdateOptions struct {
	Validate *bool  `json:"validate,omitempty"`
	Project  string `json:"project,omitempty"`
}

type ApplicationPatchOptions struct {
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
}

type TerminateOperationOptions struct {
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
}

type ResourcesQuery struct {
	Namespace    string `json:"namespace,omitempty"`
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
	Group        string `json:"group,omitempty"`
	Kind         string `json:"kind,omitempty"`
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
}

type RevisionMetadataOptions struct {
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
	SourceIndex  *int32 `json:"sourceIndex,omitempty"`
	VersionId    *int32 `json:"versionId,omitempty"`
}

type SyncWindowsOptions struct {
	AppNamespace string `json:"appNamespace,omitempty"`
	Project      string `json:"project,omitempty"`
}

type ListLinksOptions struct {
	Namespace string `json:"namespace,omitempty"`
	Project   string `json:"project,omitempty"`
}

type ApplicationStandard struct {
	client *resty.Client
}

func NewApplication(client *resty.Client) Application {
	return &ApplicationStandard{client: client}
}

func (a *ApplicationStandard) List(opts *ApplicationListOptions) (*ApplicationList, error) {
	var result ApplicationList
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Refresh != "" {
			params["refresh"] = opts.Refresh
		}
		if opts.ResourceVersion != "" {
			params["resourceVersion"] = opts.ResourceVersion
		}
		if opts.Selector != "" {
			params["selector"] = opts.Selector
		}
		if opts.Repo != "" {
			params["repo"] = opts.Repo
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
		for _, p := range opts.Project {
			req = req.SetQueryParam("project", p)
		}
	}
	resp, err := req.Get("/api/v1/applications")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) Get(name string, opts *ApplicationGetOptions) (*ApplicationModel, error) {
	var result ApplicationModel
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Refresh != "" {
			params["refresh"] = opts.Refresh
		}
		if opts.ResourceVersion != "" {
			params["resourceVersion"] = opts.ResourceVersion
		}
		if opts.Selector != "" {
			params["selector"] = opts.Selector
		}
		if opts.Repo != "" {
			params["repo"] = opts.Repo
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
		for _, p := range opts.Project {
			req = req.SetQueryParam("project", p)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) Create(app *ApplicationModel, opts *ApplicationCreateOptions) (*ApplicationModel, error) {
	var result ApplicationModel
	req := a.client.R().SetBody(app).SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Upsert != nil {
			params["upsert"] = fmt.Sprintf("%v", *opts.Upsert)
		}
		if opts.Validate != nil {
			params["validate"] = fmt.Sprintf("%v", *opts.Validate)
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Post("/api/v1/applications")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) Update(app *ApplicationModel, opts *ApplicationUpdateOptions) (*ApplicationModel, error) {
	var result ApplicationModel
	req := a.client.R().SetBody(app).SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Validate != nil {
			params["validate"] = fmt.Sprintf("%v", *opts.Validate)
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Put(fmt.Sprintf("/api/v1/applications/%s", app.Name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) Delete(name string, opts *ApplicationDeleteOptions) error {
	req := a.client.R()
	if opts != nil {
		params := map[string]string{}
		if opts.Cascade != nil {
			params["cascade"] = fmt.Sprintf("%v", *opts.Cascade)
		}
		if opts.PropagationPolicy != "" {
			params["propagationPolicy"] = opts.PropagationPolicy
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Delete(fmt.Sprintf("/api/v1/applications/%s", name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationStandard) Sync(name string, opts *SyncOptions) error {
	req := a.client.R().SetBody(opts)
	resp, err := req.Post(fmt.Sprintf("/api/v1/applications/%s/sync", name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationStandard) Rollback(name string, opts *RollbackOptions) error {
	req := a.client.R().SetBody(opts)
	resp, err := req.Post(fmt.Sprintf("/api/v1/applications/%s/rollback", name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationStandard) TerminateOperation(name string, opts *TerminateOperationOptions) error {
	req := a.client.R()
	if opts != nil {
		params := map[string]string{}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Delete(fmt.Sprintf("/api/v1/applications/%s/operation", name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationStandard) Patch(name string, patch any, patchType string, opts *ApplicationPatchOptions) (*ApplicationModel, error) {
	var result ApplicationModel
	req := a.client.R().SetBody(patch).SetResult(&result)
	if patchType != "" {
		req = req.SetHeader("Content-Type", patchType)
		req = req.SetQueryParam("patchType", patchType)
	}
	if opts != nil {
		params := map[string]string{}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Patch(fmt.Sprintf("/api/v1/applications/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) GetResource(opts *ApplicationResourceRequest) (*ApplicationResourceResponse, error) {
	var result ApplicationResourceResponse
	resp, err := a.client.R().
		SetBody(opts).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/applications/%s/resource", opts.Name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) PatchResource(opts *ApplicationResourcePatchRequest) (*ApplicationResourceResponse, error) {
	var result ApplicationResourceResponse
	resp, err := a.client.R().
		SetBody(opts).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/applications/%s/resource", opts.Name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) DeleteResource(opts *ApplicationResourceDeleteRequest) error {
	req := a.client.R().
		SetQueryParams(map[string]string{
			"name":         opts.ResourceName,
			"namespace":    opts.Namespace,
			"resourceName": opts.ResourceName,
			"version":      opts.Version,
			"group":        opts.Group,
			"kind":         opts.Kind,
		})
	if opts.Force != nil {
		req = req.SetQueryParam("force", fmt.Sprintf("%v", *opts.Force))
	}
	if opts.Orphan != nil {
		req = req.SetQueryParam("orphan", fmt.Sprintf("%v", *opts.Orphan))
	}
	resp, err := req.Delete(fmt.Sprintf("/api/v1/applications/%s/resource", opts.Name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationStandard) ListResourceActions(opts *ApplicationResourceRequest) (*ResourceActionsList, error) {
	var result ResourceActionsList
	resp, err := a.client.R().
		SetBody(opts).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/applications/%s/resource/actions", opts.Name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) RunResourceAction(opts *ApplicationResourceActionRequest) error {
	resp, err := a.client.R().
		SetBody(opts).
		Post(fmt.Sprintf("/api/v1/applications/%s/resource/actions", opts.Name))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (a *ApplicationStandard) GetManifests(name string, opts *ApplicationManifestQuery) (*ManifestResponse, error) {
	var result ManifestResponse
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Revision != "" {
			params["revision"] = opts.Revision
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if opts.NoCache != nil {
			params["noCache"] = fmt.Sprintf("%v", *opts.NoCache)
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
		for _, v := range opts.SourcePositions {
			req = req.SetQueryParam("sourcePositions", fmt.Sprintf("%d", v))
		}
		for _, v := range opts.Revisions {
			req = req.SetQueryParam("revisions", v)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/manifests", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) ResourceTree(name string, opts *ResourcesQuery) (*ApplicationTree, error) {
	var result ApplicationTree
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Namespace != "" {
			params["namespace"] = opts.Namespace
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Version != "" {
			params["version"] = opts.Version
		}
		if opts.Group != "" {
			params["group"] = opts.Group
		}
		if opts.Kind != "" {
			params["kind"] = opts.Kind
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/resource-tree", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) ManagedResources(name string, opts *ResourcesQuery) (*ManagedResourcesResponse, error) {
	var result ManagedResourcesResponse
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Namespace != "" {
			params["namespace"] = opts.Namespace
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Version != "" {
			params["version"] = opts.Version
		}
		if opts.Group != "" {
			params["group"] = opts.Group
		}
		if opts.Kind != "" {
			params["kind"] = opts.Kind
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/managed-resources", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) RevisionMetadata(name, revision string, opts *RevisionMetadataOptions) (*RevisionMetadata, error) {
	var result RevisionMetadata
	req := a.client.R().SetResult(&result).SetQueryParam("revision", revision)
	if opts != nil {
		params := map[string]string{}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if opts.SourceIndex != nil {
			params["sourceIndex"] = fmt.Sprintf("%d", *opts.SourceIndex)
		}
		if opts.VersionId != nil {
			params["versionId"] = fmt.Sprintf("%d", *opts.VersionId)
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/revision-metadata", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) GetSyncWindows(name string, opts *SyncWindowsOptions) (*SyncWindows, error) {
	var result SyncWindows
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/syncwindows", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) ListResourceEvents(name string, opts *ApplicationResourceEventsQuery) (*ResourceEventList, error) {
	var result ResourceEventList
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.ResourceName != "" {
			params["resourceName"] = opts.ResourceName
		}
		if opts.ResourceUID != "" {
			params["resourceUID"] = opts.ResourceUID
		}
		if opts.ResourceNamespace != "" {
			params["resourceNamespace"] = opts.ResourceNamespace
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/events", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) ListLinks(name string, opts *ListLinksOptions) (*LinksResponse, error) {
	var result LinksResponse
	req := a.client.R().SetResult(&result)
	if opts != nil {
		params := map[string]string{}
		if opts.Namespace != "" {
			params["namespace"] = opts.Namespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/links", name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) ListResourceLinks(opts *ApplicationResourceRequest) (*LinksResponse, error) {
	var result LinksResponse
	resp, err := a.client.R().
		SetBody(opts).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/applications/%s/resource/links", opts.Name))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (a *ApplicationStandard) Watch(ctx context.Context, opts *WatchOptions) (<-chan *ApplicationWatchEvent, error) {
	ch := make(chan *ApplicationWatchEvent)
	req := a.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Accept", "text/event-stream")
	if opts != nil {
		params := map[string]string{}
		if opts.Revision != "" {
			params["revision"] = opts.Revision
		}
		if opts.ResourceVersion != "" {
			params["resourceVersion"] = opts.ResourceVersion
		}
		if opts.Selector != "" {
			params["selector"] = opts.Selector
		}
		if opts.Repo != "" {
			params["repo"] = opts.Repo
		}
		if opts.Refresh != "" {
			params["refresh"] = opts.Refresh
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
		for _, p := range opts.Project {
			req = req.SetQueryParam("project", p)
		}
	}
	resp, err := req.Get("/api/v1/stream/applications")
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

func (a *ApplicationStandard) WatchResourceTree(ctx context.Context, name string, opts *ResourcesQuery) (<-chan *ApplicationTree, error) {
	ch := make(chan *ApplicationTree)
	req := a.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Accept", "text/event-stream")
	if opts != nil {
		params := map[string]string{}
		if opts.Namespace != "" {
			params["namespace"] = opts.Namespace
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Version != "" {
			params["version"] = opts.Version
		}
		if opts.Group != "" {
			params["group"] = opts.Group
		}
		if opts.Kind != "" {
			params["kind"] = opts.Kind
		}
		if opts.AppNamespace != "" {
			params["appNamespace"] = opts.AppNamespace
		}
		if opts.Project != "" {
			params["project"] = opts.Project
		}
		if len(params) > 0 {
			req = req.SetQueryParams(params)
		}
	}
	resp, err := req.Get(fmt.Sprintf("/api/v1/stream/applications/%s/resource-tree", name))
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

func (a *ApplicationStandard) PodLogs(ctx context.Context, opts *PodLogsOptions) (<-chan *LogEntry, error) {
	if opts == nil {
		return nil, fmt.Errorf("podLogs opts must not be nil")
	}
	ch := make(chan *LogEntry)
	req := a.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Accept", "text/event-stream")
	params := map[string]string{
		"namespace": opts.Namespace,
		"name":      opts.ResourceName,
		"group":     opts.Group,
		"kind":      opts.Kind,
		"follow":    fmt.Sprintf("%v", opts.Follow),
	}
	if opts.Container != "" {
		params["container"] = opts.Container
	}
	if opts.TailLines > 0 {
		params["tailLines"] = fmt.Sprintf("%d", opts.TailLines)
	}
	if opts.SinceSeconds > 0 {
		params["sinceSeconds"] = fmt.Sprintf("%d", opts.SinceSeconds)
	}
	if opts.SinceTime != "" {
		params["sinceTime"] = opts.SinceTime
	}
	if opts.Previous {
		params["previous"] = "true"
	}
	if opts.PodName != "" {
		params["podName"] = opts.PodName
	}
	if opts.UntilTime != "" {
		params["untilTime"] = opts.UntilTime
	}
	if opts.Filter != "" {
		params["filter"] = opts.Filter
	}
	if opts.AppNamespace != "" {
		params["appNamespace"] = opts.AppNamespace
	}
	if opts.Project != "" {
		params["project"] = opts.Project
	}
	if opts.MatchCase != nil {
		params["matchCase"] = fmt.Sprintf("%v", *opts.MatchCase)
	}
	req = req.SetQueryParams(params)
	resp, err := req.Get(fmt.Sprintf("/api/v1/applications/%s/pod-logs", opts.Name))
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

func readSSE[T any](ctx context.Context, ch chan<- T, body io.ReadCloser) {
	defer body.Close()

	go func() {
		<-ctx.Done()
		body.Close()
	}()

	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		var event struct {
			Result T `json:"result"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}
		select {
		case ch <- event.Result:
		case <-ctx.Done():
			return
		}
	}
}

func (o *PodLogsOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	if o.ResourceName == "" {
		return fmt.Errorf("resourceName is required")
	}
	return nil
}

func (o *ApplicationResourceRequest) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	if o.ResourceName == "" {
		return fmt.Errorf("resourceName is required")
	}
	if o.Version == "" {
		return fmt.Errorf("version is required")
	}
	if o.Kind == "" {
		return fmt.Errorf("kind is required")
	}
	return nil
}

func (o *ApplicationResourceDeleteRequest) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	if err := o.ApplicationResourceRequest.Validate(); err != nil {
		return err
	}
	if o.Force != nil && *o.Force && o.Orphan != nil && *o.Orphan {
		return fmt.Errorf("force and orphan cannot both be true")
	}
	return nil
}

func (o *SyncOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	return nil
}

func (o *RollbackOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	if o.ID <= 0 {
		return fmt.Errorf("id must be greater than 0")
	}
	return nil
}

func (o *ResourcesQuery) Validate() error {
	if o == nil {
		return fmt.Errorf("opts must not be nil")
	}
	return nil
}

var _ Application = (*ApplicationStandard)(nil)
