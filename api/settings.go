package api

import "github.com/go-resty/resty/v2"

type Settings interface {
	Get() (*SettingsModel, error)
	GetPlugins() (*PluginsModel, error)
}

type SettingsModel struct {
	URL                     string                   `json:"url"`
	AppLabelKey             string                   `json:"appLabelKey"`
	ControllerNamespace     string                   `json:"controllerNamespace"`
	ExecEnabled             *bool                    `json:"execEnabled,omitempty"`
	ConfigManagementPlugins []ConfigManagementPlugin `json:"configManagementPlugins,omitempty"`
	KustomizeOptions        *KustomizeOptions        `json:"kustomizeOptions,omitempty"`
	StatusBadge             string                   `json:"statusBadge,omitempty"`
	StatusBadgeRootURL      string                   `json:"statusBadgeRootUrl,omitempty"`
	UserLoginsDisabled      bool                     `json:"userLoginsDisabled"`
	PasswordPattern         string                   `json:"passwordPattern,omitempty"`
	TrackingMethod          string                   `json:"trackingMethod,omitempty"`
	GoogleAnalytics         *GoogleAnalyticsConfig   `json:"googleAnalytics,omitempty"`
	Plugins                 []Plugin                 `json:"plugins,omitempty"`
	Help                    *Help                    `json:"help,omitempty"`
}

type ConfigManagementPlugin struct {
	Name     string   `json:"name"`
	Generate Command  `json:"generate,omitempty"`
	Init     *Command `json:"init,omitempty"`
}

type Command struct {
	Command []string `json:"command"`
	Args    []string `json:"args,omitempty"`
}

type KustomizeOptions struct {
	BinaryPath   string `json:"binaryPath,omitempty"`
	BuildOptions string `json:"buildOptions,omitempty"`
}

type GoogleAnalyticsConfig struct {
	TrackingID     string `json:"trackingID"`
	AnonymizeUsers bool   `json:"anonymizeUsers"`
}

type Plugin struct {
	Name   string `json:"name"`
	Client string `json:"client,omitempty"`
}

type Help struct {
	ChatText   string            `json:"chatText,omitempty"`
	ChatURL    string            `json:"chatUrl,omitempty"`
	BinaryURLs map[string]string `json:"binaryUrls,omitempty"`
}

type PluginsModel struct {
	Plugins []*PluginInfo `json:"plugins"`
}

type PluginInfo struct {
	Name             string   `json:"name"`
	ShortDescription string   `json:"shortDescription,omitempty"`
	Description      string   `json:"description,omitempty"`
	Icon             string   `json:"icon,omitempty"`
	Source           string   `json:"source,omitempty"`
	Support          []string `json:"support,omitempty"`
	Repository       string   `json:"repository,omitempty"`
	StarCount        int      `json:"starCount,omitempty"`
	NumDownloads     int      `json:"numDownloads,omitempty"`
}

type SettingsStandard struct {
	client *resty.Client
}

func NewSettings(client *resty.Client) Settings {
	return &SettingsStandard{client: client}
}

func (s *SettingsStandard) Get() (*SettingsModel, error) {
	var result SettingsModel
	resp, err := s.client.R().
		SetResult(&result).
		Get("/api/v1/settings")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

func (s *SettingsStandard) GetPlugins() (*PluginsModel, error) {
	var result PluginsModel
	resp, err := s.client.R().
		SetResult(&result).
		Get("/api/v1/settings/plugins")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, parseError(resp)
	}
	return &result, nil
}

var _ Settings = (*SettingsStandard)(nil)
