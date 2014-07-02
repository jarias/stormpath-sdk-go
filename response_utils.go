package stormpath

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jarias/stormpath-sdk-go/logger"
)

func extractResponseData(resp *http.Response) ([]byte, error) {
	return ioutil.ReadAll(resp.Body)
}

func unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func handleResponseError(resp *http.Response, err error) ([]byte, error) {
	//Error from the request execution
	if err != nil {
		logger.ERROR.Printf("%s [%s]", err, resp.Request.URL.String())
		return []byte{}, err
	}
	//Check for Stormpath specific errors
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 302 {
		stormpathError := &StormpathError{}

		data, err := extractResponseData(resp)
		if err != nil {
			return []byte{}, err
		}

		unmarshal(data, stormpathError)

		logger.ERROR.Printf("%s [%s]", stormpathError.Message, resp.Request.URL.String())
		return []byte{}, errors.New(stormpathError.Message)
	}
	//No errors from the request execution
	return extractResponseData(resp)
}
