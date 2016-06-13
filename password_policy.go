package stormpath

type PasswordPolicy struct {
	resource
	ResetTokenTTL              int             `json:"resetTokenTtl,omitempty"`
	ResetEmailStatus           string          `json:"resetEmailStatus,omitempty"`
	ResetSuccessEmailStatus    string          `json:"resetSuccessEmailStatus,omitempty"`
	ResetEmailTemplates        *EmailTemplates `json:"resetEmailTemplates,omitempty"`
	ResetSuccessEmailTemplates *EmailTemplates `json:"resetSuccessEmailTemplates,omitempty"`
	//TODO password strength
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (policy *PasswordPolicy) Refresh() error {
	return client.get(policy.Href, policy)
}

//Update updates the given resource, by doing a POST to the resource Href
func (policy *PasswordPolicy) Update() error {
	return client.post(policy.Href, policy, policy)
}

//GetResetEmailTemplates loads the policy ResetEmailTemplates collection and returns it
func (policy *PasswordPolicy) GetResetEmailTemplates() (*EmailTemplates, error) {
	err := client.get(policy.ResetEmailTemplates.Href, policy.ResetEmailTemplates)

	if err != nil {
		return nil, err
	}

	return policy.ResetEmailTemplates, nil
}

//GetResetSuccessEmailTemplates loads the policy ResetSuccessEmailTemplates collection and returns it
func (policy *PasswordPolicy) GetResetSuccessEmailTemplates() (*EmailTemplates, error) {
	err := client.get(policy.ResetSuccessEmailTemplates.Href, policy.ResetSuccessEmailTemplates)

	if err != nil {
		return nil, err
	}

	return policy.ResetSuccessEmailTemplates, nil
}
