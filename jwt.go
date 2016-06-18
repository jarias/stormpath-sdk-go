package stormpath

import (
	"gopkg.in/dgrijalva/jwt-go.v3"
)

//SSOTokenClaims are the JWT for initiating an IDSite workflow
//
//see: http://docs.stormpath.com/guides/using-id-site/
type SSOTokenClaims struct {
	jwt.StandardClaims
	CallbackURI           string `json:"cb_uri,omitempty"`
	Path                  string `json:"path,omitempty"`
	State                 string `json:"state,omitempty"`
	OrganizationNameKey   string `json:"organizationNameKey,omitempty"`
	ShowOrganiztaionField bool   `json:"showOrganiztaionField,omitempty"`
}

//IDSiteAssertionTokenClaims are the JWT claims of an Stormpath Assertion type authentication
//this could originage from an IDSite workflow
type IDSiteAssertionTokenClaims struct {
	jwt.StandardClaims
	State  string `json:"state,omitempty"`
	Status string `json:"status,omitempty"`
}

//SAMLAssertionTokenClaims are the JWT claims of an Stormpath Assertion type authentication
//this could originage from an SAML workflow
type SAMLAssertionTokenClaims struct {
	jwt.StandardClaims
	State    string `json:"state,omitempty"`
	Status   string `json:"status,omitempty"`
	IsNewSub string `json:"isNewSub,omitempty"`
	IRT      string `json:"irt,omitempty"`
}

//SAMLAuthenticationTokenClaims are the JWT claims needed to start a Stormpath SAML workflow
type SAMLAuthenticationTokenClaims struct {
	jwt.StandardClaims
	CallbackURI string `json:"cb_uri,omitempty"`
	State       string `json:"state,omitempty"`
	ASH         string `json:"ash,omitempty"`
	ONK         string `json:"onk,omitempty"`
}

//GrantTypeStormpathTokenClaims are the JWT claims for a Stormpath OAuth2 authentication using
//the stormpath_token grant type
type GrantTypeStormpathTokenClaims struct {
	jwt.StandardClaims
	Status string `json:"status,omitempty"`
}

//GrantTypeClientCredentialsTokenClaims are the JWT claims use for the client credentials OAuth2 grant type
//authentication
type GrantTypeClientCredentialsTokenClaims struct {
	jwt.StandardClaims
	Scope string `json:"scope,omitempty"`
}

//AccessTokenClaims are the JWT for a Stormpath OAuth2 access token
type AccessTokenClaims struct {
	jwt.StandardClaims
	RefreshTokenID string `json:"rti,omitempty"`
}

//JWT helper function to create JWT token strings with the given claims, extra header values,
//and sign with client API Key Secret using SigningMethodHS256 algorithm
func JWT(claims jwt.Claims, extraHeaders map[string]interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	for key, value := range extraHeaders {
		token.Header[key] = value
	}

	encodedJWT, _ := token.SignedString(client.ClientConfiguration.GetJWTSigningKey())
	return encodedJWT
}

func ParseJWT(token string, claims jwt.Claims) *jwt.Token {
	decodedJWT, _ := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return client.ClientConfiguration.GetJWTSigningKey(), nil
	})

	return decodedJWT
}
