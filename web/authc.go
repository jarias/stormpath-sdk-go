package stormpathweb

import (
	"net/http"
	"strings"

	stormpath "github.com/jarias/stormpath-sdk-go"
)

func isAuthenticated(w http.ResponseWriter, r *http.Request, application *stormpath.Application) *stormpath.Account {
	//Cookie
	authResult := isCookieAuthenticated(r, application)
	if authResult != nil {
		saveAuthenticationResult(w, r, authResult, application)
		return authResult.GetAccount()
	}
	//Token
	tokenAuthResult := isTokenBearerAuthenticated(r, application)
	if tokenAuthResult != nil {
		saveAuthenticationResult(w, r, tokenAuthResult, application)
		return tokenAuthResult.GetAccount()
	}
	//Basic
	basicAuthResult := isHTTPBasicAuthenticated(r, application)
	if basicAuthResult != nil {
		saveAuthenticationResult(w, r, basicAuthResult, application)
		return basicAuthResult.GetAccount()
	}

	clearAuthentication(w, r, application)
	return nil
}

func isHTTPBasicAuthenticated(r *http.Request, application *stormpath.Application) stormpath.AuthResult {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil
	}

	authenticationResult, err := stormpath.NewBasicAuthenticator(application).Authenticate(username, password)
	if err != nil {
		return nil
	}

	return authenticationResult
}

func isTokenBearerAuthenticated(r *http.Request, application *stormpath.Application) stormpath.AuthResult {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil
	}

	token := authorizationHeader[strings.Index(authorizationHeader, "bearer ")+7:]

	authenticationResult, err := stormpath.NewOAuthBearerAuthenticator(application).Authenticate(token)
	if err != nil {
		return nil
	}

	return authenticationResult
}

func isCookieAuthenticated(r *http.Request, application *stormpath.Application) stormpath.AuthResult {
	isRefresh := false

	cookie, err := r.Cookie(Config.AccessTokenCookieName)
	if err == http.ErrNoCookie {
		cookie, err = r.Cookie(Config.RefreshTokenCookieName)
		if err != nil {
			return nil
		}
		isRefresh = true
	}

	if isRefresh {
		authenticationResult, err := stormpath.NewOAuthRefreshTokenAuthenticator(application).Authenticate(cookie.Value)
		if err != nil {
			return nil
		}
		return authenticationResult
	}
	//Validate the token to make sure it hasn't expire yet
	authenticationResult, err := stormpath.NewOAuthBearerAuthenticator(application).Authenticate(cookie.Value)
	if err != nil {
		return nil
	}

	return authenticationResult
}

func transientAuthenticationResult(account *stormpath.Account) *stormpath.AuthenticationResult {
	return &stormpath.AuthenticationResult{Account: account}
}

func (m StormpathMiddleware) GetAuthenticatedAccount(w http.ResponseWriter, r *http.Request) *stormpath.Account {
	return isAuthenticated(w, r, m.Application)
}
