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
	routes      map[string]stormpathHandler
}

type stormpathHandler interface {
	serveHTTP(http.ResponseWriter, *http.Request, webContext)
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

	routes := map[string]stormpathHandler{}

	if Config.LoginEnabled {
		lh := loginHandler{application: application}
		routes[Config.LoginURI] = lh
		routes[Config.FacebookCallbackURI] = facebookCallbackHandler{defaultSocialHandler{application, lh}}
		routes[Config.GoogleCallbackURI] = googleCallbackHandler{defaultSocialHandler{application, lh}}
		routes[Config.LinkedinCallbackURI] = linkedinCallbackHandler{defaultSocialHandler{application, lh}}
		routes[Config.GithubCallbackURI] = githubCallbackHandler{defaultSocialHandler{application, lh}}

	}
	if Config.LogoutEnabled {
		lgh := logoutHandler{application: application}
		routes[Config.LogoutURI] = lgh
	}
	if Config.RegisterEnabled {
		rh := registerHandler{application: application}
		routes[Config.RegisterURI] = rh
	}
	if Config.MeEnabled {
		mh := meHandler{application: application}
		routes[Config.MeURI] = mh
	}
	if isForgotPasswordEnabled(application) {
		cph := changePasswordHandler{application: application}
		routes[Config.ChangePasswordURI] = cph
	}
	if isForgotPasswordEnabled(application) {
		fph := forgotPasswordHandler{application: application}
		routes[Config.ForgotPasswordURI] = fph
	}
	if isVerifyEnabled(application) {
		evh := emailVerifyHandler{application: application}
		routes[Config.VerifyURI] = evh
	}
	if Config.CallbackEnabled {
		ch := callbackHandler{application: application}
		routes[Config.CallbackURI] = ch
	}
	if Config.OAuth2Enabled {
		oh := oauthHandler{application: application}
		routes[Config.OAuth2URI] = oh
	}

	h := &StormpathMiddleware{
		next:        next,
		routes:      routes,
		Application: application,
	}

	return h
}

func (h *StormpathMiddleware) SetPreLoginHandler(uh UserHandler) {
	if Config.LoginEnabled {
		lh := h.routes[Config.LoginURI].(loginHandler)
		lh.preLoginHandler = uh
		h.routes[Config.LoginURI] = lh
	}
}

func (h *StormpathMiddleware) SetPostLoginHandler(uh UserHandler) {
	if Config.LoginEnabled {
		lh := h.routes[Config.LoginURI].(loginHandler)
		lh.postLoginHandler = uh
		h.routes[Config.LoginURI] = lh
	}
}

func (h *StormpathMiddleware) SetPreRegisterHandler(uh UserHandler) {
	if Config.RegisterEnabled {
		rh := h.routes[Config.RegisterURI].(registerHandler)
		rh.preRegisterHandler = uh
		h.routes[Config.RegisterURI] = rh
	}
}

func (h *StormpathMiddleware) SetPostRegisterHandler(uh UserHandler) {
	if Config.RegisterEnabled {
		rh := h.routes[Config.RegisterURI].(registerHandler)
		rh.postRegisterHandler = uh
		h.routes[Config.RegisterURI] = rh
	}
}

func (h *StormpathMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resolvedContentType := resolveContentType(r)

	if resolvedContentType == "" {
		h.next.ServeHTTP(w, r)
		return
	}

	xStormpathAgent := r.Header.Get("X-Stormpath-Agent")

	stormpath.GetClient().WebSDKToken = xStormpathAgent

	if strings.HasPrefix(r.URL.Path, "/stormpath/assets/") {
		assetsHandler(w, r)
		return
	}

	account := isAuthenticated(w, r, h.Application)

	if handler, ok := h.routes[r.URL.Path]; ok {
		handler.serveHTTP(w, r, newContext(resolvedContentType, account))
		return
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
