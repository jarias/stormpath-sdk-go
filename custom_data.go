package stormpath

type customDataAwareResource struct {
	resource
	CustomData *CustomData `json:"customData,omitempty"`
}

//CustomData represents Stormpath's custom data resouce
type CustomData map[string]interface{}

func (customData CustomData) IsCacheable() bool {
	return true
}

//GetCustomData returns the given resource custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (r *customDataAwareResource) GetCustomData() (CustomData, error) {
	customData := make(CustomData)

	err := client.get(buildAbsoluteURL(r.Href, "customData"), &customData)

	if err != nil {
		return nil, err
	}

	return customData, nil
}

//UpdateCustomData sets or updates the given resource custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (r *customDataAwareResource) UpdateCustomData(customData CustomData) (CustomData, error) {
	customData = cleanCustomData(customData)

	err := client.post(buildAbsoluteURL(r.Href, "customData"), customData, &customData)

	if err != nil {
		return nil, err
	}

	return customData, nil
}

//DeleteCustomData deletes all the resource custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (r *customDataAwareResource) DeleteCustomData() error {
	return client.delete(buildAbsoluteURL(r.Href, "customData"))
}
