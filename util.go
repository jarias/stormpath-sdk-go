package stormpath

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Unmarshal(resp *http.Response, v interface{}) error {
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)

	if err != nil {
		return err
	}

	return nil
}
