package stormpathweb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/jarias/stormpath-sdk-go"
)

const (
	TextCSS               = "text/css"
	ApplicationJavascript = "application/javascript"
	NextKey               = "next"
)

var templates = make(map[string]*template.Template, 3)

//StormpathMiddleware the base http.Handler as the base Stormpath web integration
type StormpathMiddleware struct {
	//User configured handler and public paths define by user
	next http.Handler
	//Configured application
	Application *stormpath.Application
	//Integration handlers
	meHandler               meHandler
	loginHandler            loginHandler
	logoutHandler           logoutHandler
	registerHandler         registerHandler
	forgotPasswordHandler   forgotPasswordHandler
	changePasswordHandler   changePasswordHandler
	emailVerifyHandler      emailVerifyHandler
	facebookCallbackHandler facebookCallbackHandler
	googleCallbackHandler   googleCallbackHandler
	linkedinCallbackHandler linkedinCallbackHandler
	githubCallbackHandler   githubCallbackHandler
	callbackHandler         callbackHandler
	oauthHandler            oauthHandler
}

//UserHandler type to handle pre/post register or pre/post login define by the user
type UserHandler func(http.ResponseWriter, *http.Request, *stormpath.Account) bool

type internalHandler func(http.ResponseWriter, *http.Request, webContext)

//EmptyUserHandler use as default user handler for pre/post register and pre/post login user handlers
func EmptyUserHandler() UserHandler {
	return UserHandler(func(w http.ResponseWriter, r *http.Request, account *stormpath.Account) bool { return true })
}

//NewStormpathMiddleware initialize the StormpathMiddleware with the actual user application as a http.Handler
func NewStormpathMiddleware(next http.Handler, cache stormpath.Cache) *StormpathMiddleware {
	loadConfig()
	clientConfig, err := stormpath.LoadConfiguration()

	if err != nil {
		stormpath.Logger.Panicf("[ERROR] Couldn't load Stormpath client configuration: %s", err)
	}

	stormpath.Init(clientConfig, cache)

	application := resolveApplication()
	resolveAccountStores(application)

	h := &StormpathMiddleware{
		next:                  next,
		Application:           application,
		meHandler:             meHandler{Application: application},
		registerHandler:       registerHandler{Application: application},
		loginHandler:          loginHandler{Application: application},
		logoutHandler:         logoutHandler{Application: application},
		forgotPasswordHandler: forgotPasswordHandler{Application: application},
		changePasswordHandler: changePasswordHandler{Application: application},
		emailVerifyHandler:    emailVerifyHandler{Application: application},
		oauthHandler:          oauthHandler{Application: application},
		callbackHandler:       callbackHandler{Application: application},
	}

	h.facebookCallbackHandler = facebookCallbackHandler{defaultSocialHandler{application, h.loginHandler}}
	h.googleCallbackHandler = googleCallbackHandler{defaultSocialHandler{application, h.loginHandler}}
	h.linkedinCallbackHandler = linkedinCallbackHandler{defaultSocialHandler{application, h.loginHandler}}
	h.githubCallbackHandler = githubCallbackHandler{defaultSocialHandler{application, h.loginHandler}}

	return h
}

func (h *StormpathMiddleware) SetPreLoginHandler(uh UserHandler) {
	h.loginHandler.PreLoginHandler = uh
}

func (h *StormpathMiddleware) SetPostLoginHandler(uh UserHandler) {
	h.loginHandler.PostLoginHandler = uh
}

func (h *StormpathMiddleware) SetPreRegisterHandler(uh UserHandler) {
	h.registerHandler.PreRegisterHandler = uh
}

func (h *StormpathMiddleware) SetPostRegisterHandler(uh UserHandler) {
	h.registerHandler.PostRegisterHandler = uh
}

func (h *StormpathMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resolvedContentType := resolveContentType(r)

	if resolvedContentType == "" {
		h.next.ServeHTTP(w, r)
		return
	}

	xStormpathAgent := r.Header.Get("X-Stormpath-Agent")

	stormpath.GetClient().WebSDKToken = xStormpathAgent

	path := r.URL.Path

	if strings.HasPrefix(r.URL.Path, "/stormpath/assets/") {
		assetsHandler(w, r)
		return
	}

	account := isAuthenticated(w, r, h.Application)

	switch path {
	case Config.MeURI:
		if Config.MeEnabled {
			h.meHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.LoginURI:
		if Config.LoginEnabled {
			h.loginHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.RegisterURI:
		if Config.RegisterEnabled {
			h.registerHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.LogoutURI:
		if Config.LogoutEnabled {
			h.logoutHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.ForgotPasswordURI:
		if IsForgotPasswordEnabled(h.Application) {
			h.forgotPasswordHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.ChangePasswordURI:
		if IsForgotPasswordEnabled(h.Application) {
			h.changePasswordHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.VerifyURI:
		if IsVerifyEnabled(h.Application) {
			h.emailVerifyHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.FacebookCallbackURI:
		if Config.LoginEnabled {
			h.facebookCallbackHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.GoogleCallbackURI:
		if Config.LoginEnabled {
			h.googleCallbackHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.LinkedinCallbackURI:
		if Config.LoginEnabled {
			h.linkedinCallbackHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.GithubCallbackURI:
		if Config.LoginEnabled {
			h.githubCallbackHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.CallbackURI:
		if Config.CallbackEnabled {
			h.callbackHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	case Config.OAuth2URI:
		if Config.OAuth2Enabled {
			h.oauthHandler.ServeHTTP(w, r, newContext(resolvedContentType, account))
			return
		}
	}

	h.next.ServeHTTP(w, r)
	return
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Path[11:]

	data, err := Asset(location)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	if strings.HasSuffix(location, ".css") {
		w.Header().Set(stormpath.ContentTypeHeader, TextCSS)
	}
	if strings.HasSuffix(location, ".js") {
		w.Header().Set(stormpath.ContentTypeHeader, ApplicationJavascript)
	}
	w.Write(data)
}

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

func (h *StormpathMiddleware) GetAuthenticatedAccount(w http.ResponseWriter, r *http.Request) *stormpath.Account {
	return isAuthenticated(w, r, h.Application)
}

func resolveContentType(r *http.Request) string {
	produces := Config.Produces

	accept := r.Header.Get(stormpath.AcceptHeader)

	if accept == "*/*" || accept == "" || strings.HasPrefix(accept, TextCSS) || strings.HasPrefix(accept, ApplicationJavascript) {
		return produces[0]
	}

	if strings.HasPrefix(accept, stormpath.TextHTML) && contains(produces, stormpath.TextHTML) {
		return stormpath.TextHTML
	}

	if strings.HasPrefix(accept, stormpath.ApplicationJSON) && contains(produces, stormpath.ApplicationJSON) {
		return stormpath.ApplicationJSON
	}

	return ""
}

func respondJSON(w http.ResponseWriter, model interface{}, status int) {
	w.Header().Set(stormpath.ContentTypeHeader, stormpath.ApplicationJSON)
	w.WriteHeader(status)

	if model == nil {
		return
	}

	json.NewEncoder(w).Encode(model)
}

func respondHTML(w http.ResponseWriter, model interface{}, view string) {
	if templates[view] == nil {
		templateData, err := Asset("templates/" + view + ".html")
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		t, err := template.New(view).Parse(string(templateData))
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		templates[view] = t
	}

	t := templates[view]

	err := t.Execute(w, model)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
}

func accountModel(account *stormpath.Account) map[string]stormpath.Account {
	if account == nil {
		return map[string]stormpath.Account{
			"account": stormpath.Account{},
		}
	}
	accountModel := stormpath.Account{
		Username:   account.Username,
		Email:      account.Email,
		GivenName:  account.GivenName,
		MiddleName: account.MiddleName,
		FullName:   account.FullName,
		Surname:    account.Surname,
		Status:     account.Status,
	}
	accountModel.Href = account.Href
	accountModel.CreatedAt = account.CreatedAt
	accountModel.ModifiedAt = account.ModifiedAt

	for attribute, shouldExpand := range Config.MeExpand {
		switch attribute {
		case "apiKeys":
			if shouldExpand.(bool) {
				accountModel.APIKeys = account.APIKeys
			}
			break
		case "applications":
			if shouldExpand.(bool) {
				accountModel.Applications = account.Applications
			}
			break
		case "customData":
			if shouldExpand.(bool) {
				customData, err := account.GetCustomData()
				if err == nil {
					accountModel.CustomData = &customData
				}
			}
			break
		case "directory":
			if shouldExpand.(bool) {
				accountModel.Directory = account.Directory
			}
			break
		case "groupMemberships":
			if shouldExpand.(bool) {
				accountModel.GroupMemberships = account.GroupMemberships
			}
			break
		case "groups":
			if shouldExpand.(bool) {
				accountModel.Groups = account.Groups
			}
			break
		case "providerData":
			if shouldExpand.(bool) {
				accountModel.ProviderData = account.ProviderData
			}
			break
		case "tenant":
			if shouldExpand.(bool) {
				accountModel.Tenant = account.Tenant
			}
			break
		}
	}

	return map[string]stormpath.Account{
		"account": accountModel,
	}
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func transientAuthenticationResult(account *stormpath.Account) *stormpath.AuthenticationResult {
	return &stormpath.AuthenticationResult{account}
}

func handleError(w http.ResponseWriter, r *http.Request, ctx webContext, h internalHandler) {
	contentType := ctx.ContentType

	if contentType == stormpath.TextHTML {
		if r.Method == http.MethodGet {
			r.URL.RawQuery = ""
		}
		h(w, r, ctx)
		return
	}

	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, ctx.getError())
	}
}
