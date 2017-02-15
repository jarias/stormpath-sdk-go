package stormpathweb

import (
	"fmt"
	"net/http"

	"bytes"
	"encoding/json"
	"github.com/jarias/stormpath-sdk-go"
)

type forgotPasswordHandler struct {
	application *stormpath.Application
}

func (h forgotPasswordHandler) serveHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if ctx.account != nil {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			Path:        Config.IDSiteForgotURI,
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
	//No GET for application/json
	if r.Method == http.MethodGet && ctx.contentType == stormpath.TextHTML {
		h.doGET(w, r, ctx)
		return
	}

	methodNotAllowed(w, r, ctx)
}

func (h forgotPasswordHandler) doGET(w http.ResponseWriter, r *http.Request, ctx webContext) {
	model := map[string]interface{}{
		"loginURI": Config.LoginURI,
		"status":   resolveForgotPasswordStatus(r.URL.Query().Get("status")),
		"error":    ctx.webError,
	}

	respondHTML(w, model, Config.ForgotPasswordView)
}

func (h forgotPasswordHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.contentType

	data, originalData := getPostedData(r)

	if data["email"] == "" {
		handleError(w, r, ctx.withError(nil, fmt.Errorf("email is required")), h.doGET)
		return
	}

	accStoreHref := getAccStoreHref(originalData)

	h.application.SendPasswordResetEmail(data["email"], accStoreHref)

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, nil, http.StatusOK)
		return
	}

	if contentType == stormpath.TextHTML {
		http.Redirect(w, r, Config.ForgotPasswordNextURI, http.StatusFound)
		return
	}
}

func resolveForgotPasswordStatus(status string) string {
	statusMessage := ""

	switch status {
	case "invalid_sptoken":
		statusMessage = "The password reset link you tried to use is no longer valid. Please request a new link from the form below."
		break
	}

	return statusMessage
}

func getAccStoreHref(originalData []byte) string {
	data := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewBuffer(originalData)).Decode(&data); err != nil {
		//log if necessary
		return ""
	}

	if v, ok := data["accountStore"]; !ok {
		//log if necessary
		return ""
	} else {
		if accStoreMap, ok := v.(map[string]string); ok {
			if href, ok := accStoreMap["href"]; ok {
				return href
			} else {
				return ""
			}
		} else {
			return ""
		}
	}
}
