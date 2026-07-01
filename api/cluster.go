package api

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type Cluster interface {
	List() (*ClusterList, error)
	Get(server string) (*ClusterModel, error)
	Create(cluster *ClusterModel) (*ClusterModel, error)
	Update(cluster *ClusterModel) (*ClusterModel, error)
	Delete(server string) error
	RotateAuth(server string) error
	InvalidateCache(server string) error
}

type ClusterModel struct {
	TypeMeta
	ObjectMeta
	Server          string           `json:"server"`
	Name            string           `json:"name"`
	Config          ClusterConfig    `json:"config"`
	ConnectionState ConnectionState `json:"connectionState,omitempty"`
	ServerVersion   string           `json:"serverVersion,omitempty"`
	Namespaces      []string         `json:"namespaces,omitempty"`
	ClusterResources bool            `json:"clusterResources,omitempty"`
	Project         string           `json:"project,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	RefreshRequestedAt string          `json:"refreshRequestedAt,omitempty"`
	Info            ClusterInfo      `json:"info,omitempty"`
	Shard           *int64           `json:"shard,omitempty"`
}

type ClusterConfig struct {
	Username            string                     `json:"username,omitempty"`
	Password            string                     `json:"password,omitempty"`
	BearerToken         string                     `json:"bearerToken,omitempty"`
	TLSClientConfig     TLSClientConfig            `json:"tlsClientConfig"`
	AWSAuthConfig       *AWSAuthConfig             `json:"awsAuthConfig,omitempty"`
	ExecProviderConfig  *ExecProviderConfig        `json:"execProviderConfig,omitempty"`
}

type TLSClientConfig struct {
	Insecure   bool   `json:"insecure"`
	ServerName string `json:"serverName,omitempty"`
	CertData   string `json:"certData,omitempty"`
	KeyData    string `json:"keyData,omitempty"`
	CAData     string `json:"caData,omitempty"`
}

type AWSAuthConfig struct {
	ClusterName string `json:"clusterName,omitempty"`
	RoleARN     string `json:"roleARN,omitempty"`
	Profile     string `json:"profile,omitempty"`
}

type ExecProviderConfig struct {
	Command     string            `json:"command,omitempty"`
	Args        []string         `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	APIVersion  string            `json:"apiVersion,omitempty"`
	InstallHint string            `json:"installHint,omitempty"`
}

type ConnectionState struct {
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	ModifiedAt string `json:"modifiedAt,omitempty"`
}

type ClusterInfo struct {
	ConnectionState   ConnectionState `json:"connectionState"`
	ServerVersion     string          `json:"serverVersion,omitempty"`
	CacheInfo         ClusterCacheInfo `json:"cacheInfo"`
	ApplicationsCount int64           `json:"applicationsCount"`
	APIVersions       []string        `json:"apiVersions,omitempty"`
}

type ClusterCacheInfo struct {
	ResourcesCount      int64 `json:"resourcesCount,omitempty"`
	APIsCount           int64 `json:"apisCount,omitempty"`
	LastCacheSyncTime   string `json:"lastCacheSyncTime,omitempty"`
}

type ClusterList struct {
	Items []*ClusterModel `json:"items"`
}

type ClusterStandard struct {
	client *resty.Client
}

func NewCluster(client *resty.Client) Cluster {
	return &ClusterStandard{client: client}
}

func (c *ClusterStandard) List() (*ClusterList, error) {
	var result ClusterList
	resp, err := c.client.R().
		SetResult(&result).
		Get("/api/v1/clusters")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (c *ClusterStandard) Get(server string) (*ClusterModel, error) {
	var result ClusterModel
	resp, err := c.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/clusters/%s", encodeClusterServer(server)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (c *ClusterStandard) Create(cluster *ClusterModel) (*ClusterModel, error) {
	var result ClusterModel
	resp, err := c.client.R().
		SetBody(cluster).
		SetResult(&result).
		Post("/api/v1/clusters")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (c *ClusterStandard) Update(cluster *ClusterModel) (*ClusterModel, error) {
	var result ClusterModel
	resp, err := c.client.R().
		SetBody(cluster).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/clusters/%s", encodeClusterServer(cluster.Server)))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (c *ClusterStandard) Delete(server string) error {
	resp, err := c.client.R().
		Delete(fmt.Sprintf("/api/v1/clusters/%s", encodeClusterServer(server)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (c *ClusterStandard) RotateAuth(server string) error {
	resp, err := c.client.R().
		Post(fmt.Sprintf("/api/v1/clusters/%s/rotate-auth", encodeClusterServer(server)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func (c *ClusterStandard) InvalidateCache(server string) error {
	resp, err := c.client.R().
		Post(fmt.Sprintf("/api/v1/clusters/%s/invalidate-cache", encodeClusterServer(server)))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return parseError(resp)
	}
	return nil
}

func encodeClusterServer(server string) string {
	b := []byte(server)
	for i, ch := range b {
		switch ch {
		case '/':
			b[i] = '_'
		}
	}
	return url.QueryEscape(string(b))
}

var _ Cluster = (*ClusterStandard)(nil)
