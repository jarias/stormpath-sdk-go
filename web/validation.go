package stormpathweb

import "fmt"

func validateForm(form form, data map[string]string) error {
	for _, formField := range form.Fields {
		//Validate required fields
		if formField.Enabled && formField.Required && data[formField.Name] == "" {
			return fmt.Errorf("%s is required.", formField.Label)
		}
	}

	for postedField := range data {
		//Validate extra fields
		if form.getField(postedField) == nil && postedField != "customData" {
			return fmt.Errorf("%s is not a configured field", postedField)
		}
	}

	return nil
}
