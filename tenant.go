package stormpath

import "net/url"

//Tenant represents a Stormpath tennat see http://docs.stormpath.com/rest/product-guide/#tenants
type Tenant struct {
	customDataAwareResource
	Name         string        `json:"name,omitempty"`
	Key          string        `json:"key,omitempty"`
	Applications *Applications `json:"applications,omitempty"`
	Directories  *Directories  `json:"directories,omitempty"`
}

//CurrentTenant returns the current tenant see http://docs.stormpath.com/rest/product-guide/#retrieve-the-current-tenant
func CurrentTenant() (*Tenant, error) {
	tenant := &Tenant{}

	err := client.get(buildRelativeURL("tenants", "current"), tenant)

	return tenant, err
}

//CreateApplication creates a new application for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-applications
func (tenant *Tenant) CreateApplication(app *Application) error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return client.post(buildRelativeURL("applications", requestParams(extraParams)), app, app)
}

//CreateDirectory creates a new directory for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-directories
func (tenant *Tenant) CreateDirectory(dir *Directory) error {
	return client.post(buildRelativeURL("directories"), dir, dir)
}

//GetApplications returns all the applications for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-applications
func (tenant *Tenant) GetApplications(criteria Criteria) (*Applications, error) {
	apps := &Applications{}

	err := client.get(buildAbsoluteURL(tenant.Applications.Href, criteria.ToQueryString()), apps)

	return apps, err
}

//GetDirectories returns all the directories for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-directories
func (tenant *Tenant) GetDirectories(criteria Criteria) (*Directories, error) {
	directories := &Directories{}

	err := client.get(buildAbsoluteURL(tenant.Directories.Href, criteria.ToQueryString()), directories)

	return directories, err
}
