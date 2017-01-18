package stormpathweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type loginHandler struct {
	preLoginHandler  UserHandler
	postLoginHandler UserHandler
	application      *stormpath.Application
}

func (h loginHandler) serveHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	ctx.next = r.URL.Query().Get(NextKey)

	if ctx.account != nil {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			Path:        Config.IDSiteLoginURI,
			CallbackURL: baseURL(r) + Config.CallbackURI,
		}

		idSiteURL, err := h.application.CreateIDSiteURL(options)
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

func (h loginHandler) doGET(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.contentType

	model := map[string]interface{}{
		"form":          Config.LoginForm,
		"accountStores": getApplicationAccountStores(h.application),
	}

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, model, http.StatusOK)
		return
	}
	if contentType == stormpath.TextHTML {
		model["registerURI"] = Config.RegisterURI
		if isVerifyEnabled(h.application) {
			model["verifyURI"] = Config.VerifyURI
		}
		if isForgotPasswordEnabled(h.application) {
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
		model["postedData"] = ctx.postedData
		model["baseURL"] = fmt.Sprintf("http://%s/%s", r.Host, Config.BasePath)
		model["status"] = resolveLoginStatus(r.URL.Query().Get("status"))
		model["error"] = ctx.webError

		respondHTML(w, model, Config.LoginView)
	}
}

func (h loginHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx webContext) {
	var authenticationResult stormpath.AuthResult

	if h.preLoginHandler != nil {
		pre := h.preLoginHandler(w, r, nil)
		if !pre {
			//User halted execution so we return
			return
		}
	}

	contentType := ctx.contentType

	postedData, originalData := getPostedData(r)

	if _, exists := postedData["providerData"]; exists {
		//Social account
		socialAccount := &stormpath.SocialAccount{}

		json.NewDecoder(bytes.NewBuffer(originalData)).Decode(socialAccount)

		account, err := h.application.RegisterSocialAccount(socialAccount)
		if err != nil {
			handleError(w, r, ctx.withError(postedData, err), h.doGET)
			return
		}
		authenticationResult = transientAuthenticationResult(account)
	} else {
		err := validateForm(Config.LoginForm, postedData)
		if err != nil {
			handleError(w, r, ctx.withError(postedData, err), h.doGET)
			return
		}

		authenticationResult, err = stormpath.NewOAuthPasswordAuthenticator(h.application).Authenticate(postedData["login"], postedData["password"])
		if err != nil {
			handleError(w, r, ctx.withError(postedData, err), h.doGET)
			return
		}
	}

	err := saveAuthenticationResult(w, r, authenticationResult, h.application)
	if err != nil {
		handleError(w, r, ctx.withError(postedData, err), h.doGET)
		return
	}
	account := authenticationResult.GetAccount()
	if account == nil {
		handleError(w, r, ctx.withError(postedData, fmt.Errorf("can't get account from authentication result")), h.doGET)
		return
	}

	if h.postLoginHandler != nil {
		post := h.postLoginHandler(w, r, account)
		if !post {
			//User halted execution so we return
			return
		}
	}

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, accountModel(account), http.StatusOK)
		return
	}

	redirectURL := Config.LoginNextURI
	if ctx.next != "" {
		redirectURL = ctx.next
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
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
