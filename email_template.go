package stormpath

const (
	TextPlain = "text/plain"
	TextHTML  = "text/html"
)

type EmailTemplate struct {
	resource
	FromEmailAddress string            `json:"fromEmailAddress"`
	FromName         string            `json:"fromName"`
	Subject          string            `json:"subject"`
	HtmlBody         string            `json:"htmlBody"`
	TextBody         string            `json:"textBody"`
	MimeType         string            `json:"mimeType"`
	DefaultModel     map[string]string `json:"defaultModel"`
}

type EmailTemplates struct {
	collectionResource
	Items []EmailTemplate `json:"items"`
}

func GetEmailTemplate(href string) (*EmailTemplate, error) {
	emailTemplate := &EmailTemplate{}

	err := client.get(
		href,
		emptyPayload(),
		emailTemplate,
	)

	return emailTemplate, err
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (template *EmailTemplate) Refresh() error {
	return client.get(template.Href, emptyPayload(), template)
}

//Update updates the given resource, by doing a POST to the resource Href
func (template *EmailTemplate) Update() error {
	return client.post(template.Href, template, template)
}
