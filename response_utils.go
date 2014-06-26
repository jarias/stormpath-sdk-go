package stormpath

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jarias/stormpath/logger"
)

func unmarshal(resp *http.Response, v interface{}) error {
	data, err := ioutil.ReadAll(resp.Body)

	logger.INFO.Println(string(data))

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)

	if err != nil {
		return err
	}

	return nil
}

func handleStormpathErrors(resp *http.Response) error {
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 302 {
		stormpathError := &StormpathError{}

		unmarshal(resp, stormpathError)

		logger.ERROR.Println(stormpathError.Message)
		return errors.New(stormpathError.Message)
	}
	return nil
}
