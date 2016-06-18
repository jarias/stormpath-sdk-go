package stormpathweb

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/jarias/stormpath-sdk-go"
)

type oauthHandler struct {
	Application *stormpath.Application
}

func (h oauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if r.Method == http.MethodPost {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		//Check which grant types are enabled
		r.ParseForm()

		if len(r.Form) == 0 {
			respondJSON(w, map[string]string{"error": "invalid_request"}, http.StatusBadRequest)
			return
		}

		grantType := r.Form.Get("grant_type")
		if grantType == "" {
			respondJSON(w, map[string]string{"error": "invalid_request"}, http.StatusBadRequest)
			return
		}

		if (grantType == "client_credentials" && !Config.OAuth2ClientCredentialsGrantTypeEnabled) ||
			(grantType == "password" && !Config.OAuth2PasswordGrantTypeEnabled) {
			respondJSON(w, map[string]string{"error": "unsupported_grant_type"}, http.StatusBadRequest)
			return
		}

		oauthRequestAuthenticator := stormpath.NewOAuthRequestAuthenticator(h.Application)
		oauthRequestAuthenticator.TTL = Config.OAuth2ClientCredentialsGrantTypeAccessTokenTTL

		authenticationResult, err := oauthRequestAuthenticator.Authenticate(r)
		if err != nil {
			h.handleOAuth2Error(w, r, err)
			return
		}

		respondJSON(w, authenticationResult, http.StatusOK)
		return
	}
	methodNotAllowed(w, r, ctx)
}

func (h oauthHandler) handleOAuth2Error(w http.ResponseWriter, r *http.Request, err error) {
	errorModel := map[string]string{
		"error":   err.Error(),
		"message": buildErrorModel(err).Message,
	}
	status := http.StatusBadRequest

	spError, ok := err.(stormpath.Error)
	if ok {
		errorModel["error"] = spError.OAuth2Error
		respondJSON(w, errorModel, status)
		return
	}

	switch err.Error() {
	case "unsupported_grant_type":
		errorModel["error"] = err.Error()
		respondJSON(w, errorModel, status)
		return
	case "invalid_client":
		status = http.StatusUnauthorized
		respondJSON(w, errorModel, status)
		return
	}
	errorModel["error"] = "invalid_request"

	respondJSON(w, errorModel, status)
}

func exchangeToken(account *stormpath.Account, application *stormpath.Application) (*stormpath.OAuthAccessTokenResult, error) {
	claims := stormpath.GrantTypeStormpathTokenClaims{}
	claims.IssuedAt = time.Now().Unix()
	claims.Issuer = application.Href
	claims.Subject = account.Href
	claims.ExpiresAt = time.Now().Add(1 * time.Minute).Unix()
	claims.Status = "AUTHENTICATED"
	claims.Audience = stormpath.GetClient().ClientConfiguration.APIKeyID

	jwtString := stormpath.JWT(
		claims,
		map[string]interface{}{
			"kid": stormpath.GetClient().ClientConfiguration.APIKeyID,
		},
	)

	return stormpath.NewOAuthStormpathTokenAuthenticator(application).Authenticate(jwtString)
}
