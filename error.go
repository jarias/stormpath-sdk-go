package stormpath

import (
	"encoding/json"
	"net/http"
)

//Error maps a Stormpath API JSON error object which implements Go error interface
type Error struct {
	RequestID        string
	Status           int    `json:"status"`
	Code             int    `json:"code"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developerMessage"`
	MoreInfo         string `json:"moreInfo"`
	OAuth2Error      string `json:"error"`
}

func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	return e.Message + " " + e.DeveloperMessage
}

func handleResponseError(req *http.Request, resp *http.Response, err error) error {
	//Error from the request execution
	if err != nil {
		if req != nil {
			Logger.Printf("[ERROR] %s [%s]", err, req.URL.String())
		}
		return err
	}

	//Check for Stormpath specific errors
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusNoContent &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusFound {
		spError := &Error{}

		err := json.NewDecoder(resp.Body).Decode(spError)
		if err != nil {
			return err
		}
		spError.RequestID = resp.Header.Get("Stormpath-Request-Id")

		return *spError
	}
	//No errors from the request execution
	return nil
}
