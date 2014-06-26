package stormpath

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jarias/stormpath/logger"
)

func Unmarshal(resp *http.Response, v interface{}) error {
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
