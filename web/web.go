package stormpathweb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

const (
	ResolvedContentType   = "X-Go-SDK-Resolved-Content-Type"
	TextCSS               = "text/css"
	ApplicationJavascript = "application/javascript"
	AccountKey            = "account"
	ApplicationKey        = "applicaiton"
	NextKey               = "next"
)

var templates = make(map[string]*template.Template, 3)

//StormpathMiddleware the base http.Handler as the base Stormpath web integration
type StormpathMiddleware struct {
	//User configured handler and public paths define by user
	Next        http.Handler
	PublicPaths []string
	//Configured application
	Application *stormpath.Application
	//Integration handlers
	FilterChainHandler      handlerFunc
	LoginHandler            loginHandler
	LogoutHandler           logoutHandler
	RegisterHandler         registerHandler
	ForgotPasswordHandler   forgotPassordHandler
	ChangePasswordHandler   changePasswordHandler
	EmailVerifyHandler      emailVerifyHandler
	FacebookCallbackHandler facebookCallbackHandler
	GoogleCallbackHandler   googleCallbackHandler
	LinkedinCallbackHandler linkedinCallbackHandler
	GithubCallbackHandler   githubCallbackHandler
	CallbackHandler         callbackHandler
	OAuthHandler            oauthHandler
	MeHandler               meHandler
	//User handlers
	PreLoginHandler     UserHandler
	PostLoginHandler    UserHandler
	PreRegisterHandler  UserHandler
	PostRegisterHandler UserHandler
}

type handlerFunc func(http.ResponseWriter, *http.Request, context.Context)

//UserHandler type to handle pre/post register or pre/post login define by the user
type UserHandler func(http.ResponseWriter, *http.Request, context.Context) context.Context

//EmptyUserHandler use as default user handler for pre/post register and pre/post login user handlers
func EmptyUserHandler() UserHandler {
	return UserHandler(func(w http.ResponseWriter, r *http.Request, ctx context.Context) context.Context { return nil })
}

//NewStormpathMiddleware initialize the StormpathMiddleware with the actual user application as a http.Handler
func NewStormpathMiddleware(next http.Handler, publicPaths []string) *StormpathMiddleware {
	loadConfig()
	clientConfig, err := stormpath.LoadConfiguration()

	if err != nil {
		stormpath.Logger.Panicf("[ERROR] Couldn't load Stormpath client configuration: %s", err)
	}

	stormpath.Init(clientConfig, nil)

	application := resolveApplication()
	resolveAccountStores(application)

	h := &StormpathMiddleware{
		Next:                  next,
		PublicPaths:           publicPaths,
		Application:           application,
		LogoutHandler:         logoutHandler{application},
		ForgotPasswordHandler: forgotPassordHandler{application},
		ChangePasswordHandler: changePasswordHandler{application},
		EmailVerifyHandler:    emailVerifyHandler{application},
		OAuthHandler:          oauthHandler{application},
		MeHandler:             meHandler{},
	}

	h.LoginHandler = loginHandler{h, application, Config.LoginForm}
	h.RegisterHandler = registerHandler{h, application, Config.RegisterForm}
	h.FacebookCallbackHandler = facebookCallbackHandler{defaultSocialHandler{application, h.LoginHandler}}
	h.GoogleCallbackHandler = googleCallbackHandler{defaultSocialHandler{application, h.LoginHandler}}
	h.LinkedinCallbackHandler = linkedinCallbackHandler{defaultSocialHandler{application, h.LoginHandler}}
	h.GithubCallbackHandler = githubCallbackHandler{defaultSocialHandler{application, h.LoginHandler}}
	h.CallbackHandler = callbackHandler{h, application}

	h.configureFilterChainHandler()
	return h
}

func resolveAccountStores(application *stormpath.Application) {
	//see https://github.com/stormpath/stormpath-framework-spec/blob/master/configuration.md
	mappings, err := application.GetAccountStoreMappings(stormpath.MakeAccountStoreMappingsCriteria())
	if err != nil || len(mappings.Items) == 0 {
		panic(fmt.Errorf("No account stores are mapped to the specified application. Account stores are required for login and registration. \n"))
	}

	if application.DefaultAccountStoreMapping == nil && Config.RegisterEnabled {
		panic(fmt.Errorf("No default account store is mapped to the specified application. A default account store is required for registration. \n"))
	}
}

func resolveApplication() *stormpath.Application {
	//see https://github.com/stormpath/stormpath-framework-spec/blob/master/configuration.md
	applicationHref := Config.ApplicationHref
	applicationName := Config.ApplicationName

	tenant, err := stormpath.CurrentTenant()
	if err != nil {
		panic(fmt.Errorf("Fatal couldn't get current tenant: %s \n", err))
	}

	if applicationHref != "" {
		if !strings.Contains(applicationHref, "/applications/") {
			panic(fmt.Errorf("(%s) is not a valid Stormpath Application href \n", applicationHref))
		}

		application, err := stormpath.GetApplication(applicationHref, stormpath.MakeApplicationCriteria().WithDefaultAccountStoreMapping())
		if err != nil {
			panic(fmt.Errorf("The provided application could not be found. The provided application href was: %s \n", applicationHref))
		}
		return application
	}

	if applicationName != "" {
		applications, err := tenant.GetApplications(stormpath.MakeApplicationsCriteria().NameEq(applicationName).WithDefaultAccountStoreMapping())
		if err != nil || len(applications.Items) == 0 {
			panic(fmt.Errorf("The provided application could not be found. The provided application name was: %s \n", applicationName))
		}

		return &applications.Items[0]
	}

	//Get all apps if size > 1 && <= 2 return the one that's not name "Stormpath" else error

	applications, err := tenant.GetApplications(stormpath.MakeApplicationsCriteria().WithDefaultAccountStoreMapping())

	if len(applications.Items) > 2 || len(applications.Items) == 1 {
		panic(fmt.Errorf("Could not automatically resolve a Stormpath Application. Please specify your Stormpath Application in your configuration \n"))
	}

	var application *stormpath.Application

	for _, app := range applications.Items {
		if app.Name != "Stormpath" {
			application = &app
		}
	}

	return application
}

func (h *StormpathMiddleware) configureFilterChainHandler() {
	h.FilterChainHandler = handlerFunc(func(w http.ResponseWriter, r *http.Request, ctx context.Context) {
		xStormpathAgent := r.Header.Get("X-Stormpath-Agent")

		stormpath.GetClient().WebSDKToken = xStormpathAgent

		path := r.URL.Path

		if strings.HasPrefix(r.URL.Path, "/stormpath/assets/") {
			h.handleAssets(w, r)
			return
		}

		newContext, authenticated := isAuthenticated(w, r, ctx)
		if authenticated {
			ctx = newContext
		}

		switch path {
		case Config.MeURI:
			if Config.MeEnabled {
				h.MeHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.LoginURI:
			if Config.LoginEnabled {
				h.LoginHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.RegisterURI:
			if Config.RegisterEnabled {
				h.RegisterHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.LogoutURI:
			if Config.LogoutEnabled {
				h.LogoutHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.ForgotPasswordURI:
			if IsForgotPasswordEnabled(h.Application) {
				h.ForgotPasswordHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.ChangePasswordURI:
			if IsForgotPasswordEnabled(h.Application) {
				h.ChangePasswordHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.VerifyURI:
			if IsVerifyEnabled(h.Application) {
				h.EmailVerifyHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.FacebookCallbackURI:
			if Config.LoginEnabled {
				h.FacebookCallbackHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.GoogleCallbackURI:
			if Config.LoginEnabled {
				h.GoogleCallbackHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.LinkedinCallbackURI:
			if Config.LoginEnabled {
				h.LinkedinCallbackHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.GithubCallbackURI:
			if Config.LoginEnabled {
				h.GithubCallbackHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.CallbackURI:
			if Config.CallbackEnabled {
				h.CallbackHandler.ServeHTTP(w, r, ctx)
				return
			}
		case Config.OAuth2URI:
			if Config.OAuth2Enabled {
				h.OAuthHandler.ServeHTTP(w, r, ctx)
				return
			}
		}

		//If authenticated and not an integration handler pass to the user app handler
		if authenticated || contains(h.PublicPaths, r.URL.Path) {
			h.Next.ServeHTTP(w, r)
			return
		}

		unauthorizedRequest(w, r, ctx)
		return
	})
}

func (h *StormpathMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resolvedContentType := h.resolveContentType(r)
	if resolvedContentType == "" {
		h.Next.ServeHTTP(w, r)
		return
	}

	ctx := context.WithValue(context.WithValue(context.Background(), ApplicationKey, h.Application), ResolvedContentType, resolvedContentType)

	h.FilterChainHandler(w, r, ctx)
}

func (h *StormpathMiddleware) handleAssets(w http.ResponseWriter, r *http.Request) {
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

func isAuthenticated(w http.ResponseWriter, r *http.Request, ctx context.Context) (context.Context, bool) {
	if ctx.Value(AccountKey) != nil {
		return ctx, true
	}

	application := ctx.Value(ApplicationKey).(*stormpath.Application)

	//Cookie
	authResult, ok := isCookieAuthenticated(r, application)
	if ok {
		saveAuthenticationResult(w, r, authResult, application)
		return context.WithValue(ctx, AccountKey, authResult.GetAccount()), ok
	}
	//Token
	tokenAuthResult, ok := isTokenBearerAuthenticated(r, application)
	if ok {
		saveAuthenticationResult(w, r, tokenAuthResult, application)
		return context.WithValue(ctx, AccountKey, tokenAuthResult.GetAccount()), ok
	}
	//Basic
	basicAuthResult, ok := isHTTPBasicAuthenticated(r, application)
	if ok {
		saveAuthenticationResult(w, r, basicAuthResult, application)
		return context.WithValue(ctx, AccountKey, basicAuthResult.GetAccount()), ok
	}

	clearAuthentication(w, r, application)
	return ctx, false
}

func isHTTPBasicAuthenticated(r *http.Request, application *stormpath.Application) (stormpath.AuthResult, bool) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, false
	}

	authenticationResult, err := stormpath.NewBasicAuthenticator(application).Authenticate(username, password)
	if err != nil {
		return nil, false
	}

	return authenticationResult, true
}

func isTokenBearerAuthenticated(r *http.Request, application *stormpath.Application) (stormpath.AuthResult, bool) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, false
	}

	token := authorizationHeader[strings.Index(authorizationHeader, "bearer ")+7:]

	authenticationResult, err := stormpath.NewOAuthBearerAuthenticator(application).Authenticate(token)
	if err != nil {
		return nil, false
	}

	return authenticationResult, true
}

func isCookieAuthenticated(r *http.Request, application *stormpath.Application) (stormpath.AuthResult, bool) {
	isRefresh := false

	cookie, err := r.Cookie(Config.AccessTokenCookieName)
	if err == http.ErrNoCookie {
		cookie, err = r.Cookie(Config.RefreshTokenCookieName)
		if err != nil {
			return nil, false
		}
		isRefresh = true
	}

	if isRefresh {
		authenticationResult, err := stormpath.NewOAuthRefreshTokenAuthenticator(application).Authenticate(cookie.Value)
		if err != nil {
			return nil, false
		}
		return authenticationResult, true
	}
	//Validate the token to make sure it hasn't expire yet
	authenticationResult, err := stormpath.NewOAuthBearerAuthenticator(application).Authenticate(cookie.Value)
	if err != nil {
		return nil, false
	}

	return authenticationResult, true
}

func (h *StormpathMiddleware) GetAuthenticatedAccount(w http.ResponseWriter, r *http.Request) *stormpath.Account {
	ctx, ok := isAuthenticated(w, r, context.WithValue(context.Background(), ApplicationKey, h.Application))
	if ok {
		return ctx.Value(AccountKey).(*stormpath.Account)
	}
	return nil
}

func (h *StormpathMiddleware) resolveContentType(r *http.Request) string {
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

func contextWithError(ctx context.Context, err error, postedData map[string]string) context.Context {
	return context.WithValue(context.WithValue(ctx, "postedData", postedData), "error", buildErrorModel(err))
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
