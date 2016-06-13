package stormpathweb

import (
	"fmt"
	"html/template"
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type loginHandler struct {
	Parent      *StormpathMiddleware
	Application *stormpath.Application
	Form        form
}

func (h loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	ctx = context.WithValue(ctx, NextKey, r.URL.Query().Get(NextKey))

	if _, ok := isAuthenticated(w, r, ctx); ok {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			Path:        Config.IDSiteLoginURI,
			CallbackURL: baseURL(r) + Config.CallbackURI,
		}

		idSiteURL, err := h.Application.CreateIDSiteURL(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, idSiteURL, http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		h.doPOST(w, r, ctx)
		return
	}
	if r.Method == http.MethodGet {
		h.doGET(w, r, ctx)
		return
	}

	methodNotAllowed(w, r, ctx)
}

func (h loginHandler) doGET(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	model := map[string]interface{}{
		"form":          h.Form,
		"accountStores": getApplicationAccountStores(h.Application),
	}

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, model, http.StatusOK)
		return
	}
	if contentType == stormpath.TextHTML {
		model["registerURI"] = Config.RegisterURI
		if IsVerifyEnabled(h.Application) {
			model["verifyURI"] = Config.VerifyURI
		}
		if IsForgotPasswordEnabled(h.Application) {
			model["forgotURI"] = Config.ForgotPasswordURI
		}
		//Social
		model["googleCallbackUri"] = Config.GoogleCallbackURI
		model["googleScope"] = Config.GoogleScope
		model["githubCallbackUri"] = Config.GithubCallbackURI
		model["githubScope"] = Config.GithubScope
		model["facebookCallbackUri"] = Config.FacebookCallbackURI
		model["facebookScope"] = Config.FacebookScope
		model["linkedinCallbackUri"] = Config.LinkedinCallbackURI
		model["linkedinScope"] = Config.LinkedinScope
		//End Social
		model["postedData"] = ctx.Value("postedData")
		model["baseURL"] = fmt.Sprintf("http://%s/%s", r.Host, Config.BasePath)
		model["status"] = resolveLoginStatus(r.URL.Query().Get("status"))
		model["error"] = ctx.Value("error")

		respondHTML(w, model, Config.LoginView)
	}
}

func (h loginHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var authenticationResult stormpath.AuthResult

	if h.Parent.PreLoginHandler != nil {
		pre := h.Parent.PreLoginHandler(w, r, ctx)
		if pre != nil {
			//User halted execution so we return
			return
		}
	}

	contentType := ctx.Value(ResolvedContentType)

	postedData, originalData := getPostedData(r)

	if _, exists := postedData["providerData"]; exists {
		//Social account
		socialAccount := &stormpath.SocialAccount{}

		json.NewDecoder(bytes.NewBuffer(originalData)).Decode(socialAccount)

		account, err := h.Application.RegisterSocialAccount(socialAccount)
		if err != nil {
			h.handlePostError(w, r, ctx, err, postedData)
			return
		}
		authenticationResult = transientAuthenticationResult(account)
	} else {
		err := validateForm(h.Form, postedData)
		if err != nil {
			h.handlePostError(w, r, ctx, err, postedData)
			return
		}

		authenticationResult, err = stormpath.NewOAuthPasswordAuthenticator(h.Application).Authenticate(postedData["login"], postedData["password"])
		if err != nil {
			h.handlePostError(w, r, ctx, err, postedData)
			return
		}
	}

	err := saveAuthenticationResult(w, r, authenticationResult, h.Application)
	if err != nil {
		h.handlePostError(w, r, ctx, err, postedData)
		return
	}
	account := authenticationResult.GetAccount()
	if account == nil {
		h.handlePostError(w, r, ctx, fmt.Errorf("can't get account from authentication result"), postedData)
		return
	}

	if h.Parent.PostLoginHandler != nil {
		post := h.Parent.PostLoginHandler(w, r, context.WithValue(ctx, AccountKey, account))
		if post != nil {
			//User halted execution so we return
			return
		}
	}

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, accountModel(account), http.StatusOK)
		return
	}

	redirectURL := Config.LoginNextURI
	if ctx.Value(NextKey) != "" {
		redirectURL = ctx.Value(NextKey).(string)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h loginHandler) handlePostError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error, postedData map[string]string) {
	contentType := ctx.Value(ResolvedContentType)

	if contentType == stormpath.TextHTML {
		//Sanitize postedData
		postedData["password"] = ""
		postedData["confirmPassword"] = ""
		h.doGET(w, r, contextWithError(ctx, err, postedData))
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, err)
	}
}

func resolveLoginStatus(status string) template.HTML {
	statusMessage := ""

	switch status {
	case "unverified":
		statusMessage = fmt.Sprintf("Your account verification email has been sent! Before you can log into your account, you need to activate your account by clicking the link we sent to your inbox. Didn't get the email? <a href=\"%s\">Click Here</a>", Config.VerifyURI)
		break
	case "verified":
		statusMessage = "Your Account Has Been Verified. You may now login."
		break
	case "created":
		statusMessage = "Your Account Has Been Created. You may now login."
		break
	case "forgot":
		statusMessage = "Password Reset Requested. If an account exists for the email provided, you will receive an email shortly."
		break
	case "reset":
		statusMessage = "Password Reset Successfully. You can now login with your new password."
	}

	return template.HTML(statusMessage)
}
