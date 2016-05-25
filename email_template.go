package stormpath

const (
	//TextPlain "text/plain" mime type
	TextPlain = "text/plain"
	//TextHTML "text/html" mime type
	TextHTML = "text/html"
)

//EmailTemplate represents an account creation policy email template
type EmailTemplate struct {
	resource
	FromEmailAddress string            `json:"fromEmailAddress"`
	FromName         string            `json:"fromName"`
	Subject          string            `json:"subject"`
	HTMLBody         string            `json:"htmlBody"`
	TextBody         string            `json:"textBody"`
	MimeType         string            `json:"mimeType"`
	DefaultModel     map[string]string `json:"defaultModel"`
}

//EmailTemplates represents a collection of EmailTemplate
type EmailTemplates struct {
	collectionResource
	Items []EmailTemplate `json:"items"`
}

//GetEmailTemplate loads an email template by href
func GetEmailTemplate(href string) (*EmailTemplate, error) {
	emailTemplate := &EmailTemplate{}

	err := client.get(
		href,
		emailTemplate,
	)

	if err != nil {
		return nil, err
	}

	return emailTemplate, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (template *EmailTemplate) Refresh() error {
	return client.get(template.Href, template)
}

//Update updates the given resource, by doing a POST to the resource Href
func (template *EmailTemplate) Update() error {
	return client.post(template.Href, template, template)
}
