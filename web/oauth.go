package stormpathweb

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/dgrijalva/jwt-go"
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
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["iat"] = time.Now().Unix()
	token.Claims["iss"] = application.Href
	token.Claims["sub"] = account.Href
	token.Claims["exp"] = time.Now().Add(1 * time.Minute).Unix()
	token.Claims["status"] = "AUTHENTICATED"
	token.Claims["aud"] = stormpath.GetClient().ClientConfiguration.APIKeyID
	token.Header["kid"] = stormpath.GetClient().ClientConfiguration.APIKeyID

	tokenString, err := token.SignedString([]byte(stormpath.GetClient().ClientConfiguration.APIKeySecret))
	if err != nil {
		return nil, err
	}

	return stormpath.NewOAuthStormpathTokenAuthenticator(application).Authenticate(tokenString)
}
