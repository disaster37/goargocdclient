package api

import "github.com/go-resty/resty/v2"

type API interface {
	Account() Account
	Application() Application
	ApplicationSet() ApplicationSet
	Certificate() Certificate
	Cluster() Cluster
	GPGKey() GPGKey
	Notification() Notification
	Project() Project
	RepoCreds() RepoCreds
	Repository() Repository
	Session() Session
	Settings() Settings
	Version() Version
}

type APIStandard struct {
	account        Account
	application    Application
	applicationSet ApplicationSet
	certificate    Certificate
	cluster        Cluster
	gpgKey         GPGKey
	notification   Notification
	project        Project
	repoCreds      RepoCreds
	repository     Repository
	session        Session
	settings       Settings
	version        Version
}

func New(client *resty.Client) API {
	return &APIStandard{
		account:        NewAccount(client),
		application:    NewApplication(client),
		applicationSet: NewApplicationSet(client),
		certificate:    NewCertificate(client),
		cluster:        NewCluster(client),
		gpgKey:         NewGPGKey(client),
		notification:   NewNotification(client),
		project:        NewProject(client),
		repoCreds:      NewRepoCreds(client),
		repository:     NewRepository(client),
		session:        NewSession(client),
		settings:       NewSettings(client),
		version:        NewVersion(client),
	}
}

func (a *APIStandard) Account() Account               { return a.account }
func (a *APIStandard) Application() Application       { return a.application }
func (a *APIStandard) ApplicationSet() ApplicationSet { return a.applicationSet }
func (a *APIStandard) Certificate() Certificate       { return a.certificate }
func (a *APIStandard) Cluster() Cluster               { return a.cluster }
func (a *APIStandard) GPGKey() GPGKey                 { return a.gpgKey }
func (a *APIStandard) Notification() Notification     { return a.notification }
func (a *APIStandard) Project() Project               { return a.project }
func (a *APIStandard) RepoCreds() RepoCreds           { return a.repoCreds }
func (a *APIStandard) Repository() Repository         { return a.repository }
func (a *APIStandard) Session() Session               { return a.session }
func (a *APIStandard) Settings() Settings             { return a.settings }
func (a *APIStandard) Version() Version               { return a.version }
