package stormpath

import (
	"net/url"
)

type Tenant struct {
	Href         string `json:"href"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	Applications link   `json:"applications"`
	Directories  link   `json:"directories"`
}

func CurrentTenant() (*Tenant, error) {
	tenant := &Tenant{}

	err := Client.doWithResult(
		Client.newRequestWithoutRedirects(
			"GET",
			buildURL("tenants", "current"),
			emptyPayload(),
		), tenant)

	return tenant, err
}

func (tenant *Tenant) CreateApplication(app *Application) error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return Client.doWithResult(
		Client.newRequest(
			"POST",
			buildURL("applications", requestParams(nil, nil, extraParams)),
			app,
		), app)
}

func (tenant *Tenant) CreateDirectory(dir *Directory) error {
	return Client.doWithResult(
		Client.newRequest(
			"POST",
			buildURL("directories"),
			dir,
		), dir)
}

func (tenant *Tenant) GetApplications(pageRequest PageRequest, filter Filter) (*Applications, error) {
	apps := &Applications{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		tenant.Applications.Href+requestParams(&pageRequest, filter, url.Values{}),
		emptyPayload(),
	), apps)

	return apps, err
}

func (tenant *Tenant) GetDirectories(pageRequest PageRequest, filter Filter) (*Directories, error) {
	directories := &Directories{}

	err := Client.doWithResult(Client.newRequest(
		"GET",
		tenant.Directories.Href+requestParams(&pageRequest, filter, url.Values{}),
		emptyPayload(),
	), directories)

	return directories, err
}
