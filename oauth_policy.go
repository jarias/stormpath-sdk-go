package stormpath

//OAuthPolicy holds the application related OAuth configuration
type OAuthPolicy struct {
	resource
	AccessTokenTtl  string `json:"accessTokenTtl"`
	RefreshTokenTtl string `json:"refreshTokenTtl"`
}

func (policy *OAuthPolicy) Update() error {
	return client.post(policy.Href, policy, policy)
}
