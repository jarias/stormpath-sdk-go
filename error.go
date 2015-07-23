package stormpath

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Error maps a Stormpath API JSON error object which implements Go error interface
type Error struct {
	Status           int
	Code             int
	Message          string
	DeveloperMessage string
	MoreInfo         string
}

func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	return fmt.Sprintf("Stormpath request error \nCode: [ %d ]\nMessage: [ %s ]\nDeveloper Message: [ %s ]\nMore info [ %s ]", e.Code, e.Message, e.DeveloperMessage, e.MoreInfo)
}

func handleResponseError(resp *http.Response, err error) error {
	//Error from the request execution
	if err != nil {
		Logger.Printf("[ERROR] %s [%s]", err, resp.Request.URL.String())
		return err
	}
	//Check for Stormpath specific errors
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 302 {
		spError := &Error{}

		err := json.NewDecoder(resp.Body).Decode(spError)
		if err != nil {
			return err
		}

		Logger.Printf("[ERROR] %s", spError)
		return *spError
	}
	//No errors from the request execution
	return nil
}
