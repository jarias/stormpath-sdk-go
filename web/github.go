package stormpathweb

import (
	"net/http"
	"net/url"

	"encoding/json"
	"strings"

	"fmt"

	"github.com/jarias/stormpath-sdk-go"
)

const githubAccessTokenURL = "https://github.com/login/oauth/access_token"

type githubCallbackHandler struct {
	defaultSocialHandler
}

func (h githubCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if r.Method == http.MethodGet {
		code := r.URL.Query().Get("code")

		accessToken, err := h.exchangeCode(code)
		if err != nil {
			h.LoginHandler.doGET(w, r, ctx.withError(nil, err))
		}

		socialAccount := &stormpath.SocialAccount{
			Data: stormpath.ProviderData{
				ProviderID:  "github",
				AccessToken: accessToken,
			},
		}

		h.authenticateSocial(w, r, ctx, socialAccount)
		return
	}

	methodNotAllowed(w, r, ctx)
}

func (h githubCallbackHandler) exchangeCode(code string) (string, error) {
	for _, accountStore := range getApplicationAccountStores(h.Application) {
		if accountStore.Provider.ProviderID == "github" {
			values := url.Values{
				"code":          {code},
				"client_id":     {accountStore.Provider.ClientID},
				"client_secret": {accountStore.Provider.ClientSecret},
			}

			request, err := http.NewRequest(http.MethodPost, githubAccessTokenURL, strings.NewReader(values.Encode()))
			if err != nil {
				return "", err
			}

			request.Header.Add("Accept", stormpath.ApplicationJSON)
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				return "", err
			}

			result := map[string]string{}

			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				return "", err
			}

			return result["access_token"], nil
		}
	}

	return "", fmt.Errorf("No GitHub account store for the configured application.")
}
