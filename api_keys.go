package stormpath

import "net/url"

//APIKey represents an Account key id/secret pair resource
//
//See: https://docs.stormpath.com/rest/product-guide/latest/reference.html#account-api-keys
type APIKey struct {
	resource
	ID      string   `json:"id"`
	Secret  string   `json:"secret"`
	Status  string   `json:"status"`
	Account *Account `json:"account"`
	Tenant  *Tenant  `json:"tenant"`
}

//APIKeys represents a collection of APIKey resources
type APIKeys struct {
	collectionResource
	Items []APIKey `json:"items,omitempty"`
}

//APIKeyCriteria represents the criteria type for the APIKey resource
type APIKeyCriteria struct {
	baseCriteria
}

//MakeAPIKeyCriteria creates a default APIKeyCriteria for a single APIKey resource
func MakeAPIKeyCriteria() APIKeyCriteria {
	return APIKeyCriteria{baseCriteria{filter: url.Values{}}}
}

//MakeAPIKeysCriteria creates a default APIKeyCriteria for a APIKeys collection resource
func MakeAPIKeysCriteria() APIKeyCriteria {
	return APIKeyCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//GetAPIKey retrives an APIKey resource by href and optional criteria
func GetAPIKey(href string, criteria APIKeyCriteria) (*APIKey, error) {
	apiKey := &APIKey{}

	err := client.get(
		buildAbsoluteURL(href, criteria.toQueryString()),
		apiKey,
	)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

//Delete deletes a given APIKey
func (k *APIKey) Delete() error {
	return client.delete(k.Href)
}

//Update updates the given APIKey against Stormpath
func (k *APIKey) Update() error {
	return client.post(k.Href, map[string]string{"status": k.Status}, k)
}

//WithAccount adds the account expansion to the given APIKeyCriteria
func (c APIKeyCriteria) WithAccount() APIKeyCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "account")
	return c
}

//WithTenant adds the tenant expansion to the given APIKeyCriteria
func (c APIKeyCriteria) WithTenant() APIKeyCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

//IDEq adds the id filter to the given APIKeyCriteria
func (c APIKeyCriteria) IDEq(id string) APIKeyCriteria {
	c.filter.Add("id", id)
	return c
}
