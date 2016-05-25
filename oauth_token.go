package stormpath

import "net/url"

//OAuthToken represents the Stormpath OAuthToken see: https://docs.stormpath.com/guides/token-management/
type OAuthToken struct {
	resource
	Account     *Account     `json:"account"`
	Application *Application `json:"application"`
	Tenant      *Tenant      `json:"tenant"`
	JWT         string       `json:"jwt"`
	ExpandedJWT ExpandedJWT  `json:"expandedJwt"`
}

//ExpandedJWT represents the OAuth token expanded JWT information
type ExpandedJWT struct {
	Claims    Claims `json:"claims"`
	Header    Header `json:"header"`
	Signature string `json:"signature"`
}

//Claims represents the expanded JWT claims
type Claims struct {
	EXP int64  `json:"exp"`
	IAT int64  `json:"iat"`
	ISS string `json:"iss"`
	JTI string `json:"jti"`
	RTI string `json:"rti"`
	SUB string `json:"sub"`
}

//Header represents the expanded JWT header
type Header struct {
	ALG string `json:"alg"`
	KID string `json:"kid"`
}

//OAuthTokens collection type for OAuthToken
type OAuthTokens struct {
	collectionResource
	Items []OAuthToken `json:"items"`
}

//OAuthResponse represents an OAuth2 response from StormPath
type OAuthResponse struct {
	AccessToken              string `json:"access_token"`
	RefreshToken             string `json:"refresh_token"`
	TokenType                string `json:"token_type"`
	ExpiresIn                int    `json:"expires_in"`
	StormpathAccessTokenHref string `json:"stormpath_access_token_href"`
}

type OAuthTokenCriteria struct {
	baseCriteria
}

func MakeOAuthTokensCriteria() OAuthTokenCriteria {
	return OAuthTokenCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Delete deletes the given OAuthToken
func (t *OAuthToken) Delete() error {
	return client.delete(t.Href, emptyPayload())
}
