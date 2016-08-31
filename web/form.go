package stormpathweb

import (
	"encoding/json"
	"net/http"
	"strings"

	"io/ioutil"

	"bytes"

	"github.com/jarias/stormpath-sdk-go"
)

type form struct {
	Fields []field `json:"fields"`
}

type field struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"-"`
	Visible     bool   `json:"-"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Label       string `json:"label"`
	PlaceHolder string `json:"placeholder"`
}

func (f form) getField(fieldName string) *field {
	for _, field := range f.Fields {
		if field.Name == fieldName {
			return &field
		}
	}
	return nil
}

func sanitizePostedData(postedData map[string]string) map[string]string {
	if postedData == nil {
		return nil
	}
	postedData["password"] = ""
	postedData["confirmPassword"] = ""
	return postedData
}

func getPostedData(r *http.Request) (map[string]string, []byte) {
	data := map[string]string{}

	contentType := r.Header.Get(stormpath.ContentTypeHeader)

	if strings.Contains(contentType, stormpath.ApplicationJSON) {
		originalData, _ := ioutil.ReadAll(r.Body)
		json.NewDecoder(bytes.NewBuffer(originalData)).Decode(&data)
		return data, originalData
	}
	if strings.Contains(contentType, stormpath.ApplicationFormURLencoded) {
		r.ParseForm()
		for param, value := range r.Form {
			if param != "next" && param != "status" {
				data[param] = value[0]
			}
		}
		return data, nil
	}
	return data, nil
}
