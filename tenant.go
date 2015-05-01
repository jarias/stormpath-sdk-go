package stormpath

import (
	"net/url"
)

//Tenant represents a Stormpath tennat see http://docs.stormpath.com/rest/product-guide/#tenants
type Tenant struct {
	Href         string `json:"href"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	CustomData   link   `json:"customData"`
	Applications link   `json:"applications"`
	Directories  link   `json:"directories"`
}

//CurrentTenant returns the current tenant see http://docs.stormpath.com/rest/product-guide/#retrieve-the-current-tenant
func CurrentTenant() (*Tenant, error) {
	tenant := &Tenant{}

	err := client.doWithResult(
		client.newRequest(
			"GET",
			buildRelativeURL("tenants", "current"),
			emptyPayload(),
		), tenant)

	return tenant, err
}

//CreateApplication creates a new application for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-applications
func (tenant *Tenant) CreateApplication(app *Application) error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return client.post(buildRelativeURL("applications", requestParams(nil, nil, extraParams)), app, app)
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
func (tenant *Tenant) GetApplications(pageRequest url.Values, filter url.Values) (*Applications, error) {
	apps := &Applications{}

	err := client.get(buildAbsoluteURL(tenant.Applications.Href, requestParams(pageRequest, filter, url.Values{})), emptyPayload(), apps)

	return apps, err
}

//GetDirectories returns all the directories for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-directories
func (tenant *Tenant) GetDirectories(pageRequest url.Values, filter url.Values) (*Directories, error) {
	directories := &Directories{}

	err := client.get(buildAbsoluteURL(tenant.Directories.Href, requestParams(pageRequest, filter, url.Values{})), emptyPayload(), directories)

	return directories, err
}

//UpdateCustomData updates the tenant custom data and returns that updated custom data as a map[string]interface
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (tenant *Tenant) UpdateCustomData(customData map[string]interface{}) (map[string]interface{}, error) {
	updatedCustomData := map[string]interface{}{}

	err := client.post(tenant.CustomData.Href, customData, &updatedCustomData)

	return updatedCustomData, err
}

//DeleteCustomData deletes all the tenants custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (tenant *Tenant) DeleteCustomData() error {
	return client.delete(tenant.CustomData.Href, emptyPayload())
}

//GetCustomData gets the tenant custom data map
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (tenant *Tenant) GetCustomData() (map[string]interface{}, error) {
	customData := map[string]interface{}{}

	err := client.get(tenant.CustomData.Href, emptyPayload(), &customData)

	return customData, err
}
