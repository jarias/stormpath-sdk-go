package stormpathweb

import (
	"encoding/json"
	"net/http"
	"strings"

	"io/ioutil"

	"bytes"

	"github.com/jarias/stormpath-sdk-go"
	"github.com/spf13/viper"
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

func buildForm(formName string) form {
	form := form{}

	for _, fieldName := range getConfiguredFormFieldNames(formName) {
		field := field{
			Name:        fieldName,
			Label:       viper.GetString("stormpath.web." + formName + ".form.fields." + fieldName + ".label"),
			PlaceHolder: viper.GetString("stormpath.web." + formName + ".form.fields." + fieldName + ".placeHolder"),
			Visible:     viper.GetBool("stormpath.web." + formName + ".form.fields." + fieldName + ".visible"),
			Enabled:     viper.GetBool("stormpath.web." + formName + ".form.fields." + fieldName + ".enabled"),
			Required:    viper.GetBool("stormpath.web." + formName + ".form.fields." + fieldName + ".required"),
			Type:        viper.GetString("stormpath.web." + formName + ".form.fields." + fieldName + ".type"),
		}
		if field.Enabled {
			form.Fields = append(form.Fields, field)
		}
	}

	return form
}

func getConfiguredFormFieldNames(formName string) []string {
	configuredFields := viper.GetStringMapString("stormpath.web." + formName + ".form.fields")
	fieldOrder := viper.GetStringSlice("stormpath.web." + formName + ".form.fieldOrder")

	for fieldName := range configuredFields {
		if !contains(fieldOrder, fieldName) {
			fieldOrder = append(fieldOrder, fieldName)
		}
	}
	return fieldOrder
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
