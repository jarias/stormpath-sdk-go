package stormpath

import (
	"net/url"
)

const tenantBaseUrl = "https://api.stormpath.com/v1/tenants"

type Tenant struct {
	Href         string
	Name         string
	Key          string
	Applications struct {
		Href string
	}
	Directories struct {
		Href string
	}
	Client *StormpathClient
}

type Applications struct {
	Href   string
	Offset int
	Limit  int
	Items  []Application
}

func CurrentTenant(credentials *Credentials) (*Tenant, error) {
	tenant := &Tenant{Client: NewStormpathClient(credentials)}

	resp, err := tenant.Client.Do(NewStormpathRequestNoRedirects("GET", tenantBaseUrl+"/current", url.Values{}))

	if err != nil {
		return nil, err
	}

	location := resp.Header.Get("Location")

	resp, err = tenant.Client.Do(NewStormpathRequest("GET", location, url.Values{}))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, tenant)

	return tenant, err
}

func (tenant *Tenant) GetApplications(pageRequest *PageRequest) (*Applications, error) {
	apps := &Applications{}

	resp, err := tenant.Client.Do(NewStormpathRequest(GET, tenant.Applications.Href, pageRequest.ToUrlQueryValues()))

	if err != nil {
		return nil, err
	}

	err = Unmarshal(resp, apps)

	return apps, err
}
