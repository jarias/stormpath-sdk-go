package stormpath

import (
	"net/url"
)

const (
	TenantBaseUrl  = "https://api.stormpath.com/v1/tenants"
	LocationHeader = "Location"
)

type Tenant struct {
	Href         string `json:"href"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	Applications Link   `json:"applications"`
	Directories  Link   `json:"directories"`
}

func CurrentTenant() (*Tenant, error) {
	tenant := &Tenant{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:              Get,
		URL:                 TenantBaseUrl + "/current",
		DontFollowRedirects: true,
	}, tenant)

	return tenant, err
}

func (tenant *Tenant) CreateApplication(app *Application) error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return Client.DoWithResult(&StormpathRequest{
		Method:      Post,
		URL:         ApplicationBaseUrl,
		Payload:     app,
		ExtraParams: extraParams,
	}, app)
}

func (tenant *Tenant) CreateDirectory(dir *Directory) error {
	return Client.DoWithResult(&StormpathRequest{
		Method:  Post,
		URL:     DirectoryBaseUrl,
		Payload: dir,
	}, dir)
}

func (tenant *Tenant) GetApplications(pageRequest PageRequest, filters DefaultFilter) (*Applications, error) {
	apps := &Applications{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      Get,
		URL:         tenant.Applications.Href,
		PageRequest: pageRequest,
		Filter:      filters,
	}, apps)

	return apps, err
}

func (tenant *Tenant) GetDirectories(pageRequest PageRequest, filters DefaultFilter) (*Directories, error) {
	directories := &Directories{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      Get,
		URL:         tenant.Directories.Href,
		PageRequest: pageRequest,
		Filter:      filters,
	}, directories)

	return directories, err
}
