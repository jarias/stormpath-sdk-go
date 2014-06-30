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

func handleStormpathErrors(resp *http.Response) error {
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 302 {
		stormpathError := &StormpathError{}

		data, err := extractResponseData(resp)
		if err != nil {
			return err
		}

		unmarshal(data, stormpathError)

		logger.ERROR.Println(stormpathError.Message)
		return errors.New(stormpathError.Message)
	}
	return nil
}
