package stormpath

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

	resp, err := Client.Do(&StormpathRequest{
		Method:              GET,
		URL:                 TenantBaseUrl + "/current",
		DontFollowRedirects: true,
	})

	if err != nil {
		return nil, err
	}

	location := resp.Header.Get(LocationHeader)

	err = Client.DoWithResult(&StormpathRequest{
		Method: GET,
		URL:    location,
	}, tenant)

	return tenant, err
}

func (tenant *Tenant) GetApplications(pageRequest PageRequest, filters DefaultFilter) (*Applications, error) {
	apps := &Applications{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         tenant.Applications.Href,
		PageRequest: &pageRequest,
		Filter:      filters,
	}, apps)

	return apps, err
}

func (tenant *Tenant) GetDirectories(pageRequest PageRequest, filters DefaultFilter) (*Directories, error) {
	directories := &Directories{}

	err := Client.DoWithResult(&StormpathRequest{
		Method:      GET,
		URL:         tenant.Directories.Href,
		PageRequest: &pageRequest,
		Filter:      filters,
	}, directories)

	return directories, err
}
