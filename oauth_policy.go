package stormpath

//OAuthPolicy holds the application related OAuth configuration
type OAuthPolicy struct {
	resource
	AccessTokenTTL  string `json:"accessTokenTtl,omitempty"`
	RefreshTokenTTL string `json:"refreshTokenTtl,omitempty"`
}

//GetOAuthPolicy return the application OAuthPolicy
func (app *Application) GetOAuthPolicy() (*OAuthPolicy, error) {
	oauthPolicy := &OAuthPolicy{}

	err := client.get(app.OAuthPolicy.Href, oauthPolicy)

	return oauthPolicy, err
}

//Update OAuthPolicy
func (policy *OAuthPolicy) Update() error {
	return client.post(policy.Href, policy, policy)
}
