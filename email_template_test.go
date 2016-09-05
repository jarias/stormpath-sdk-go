package stormpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmailTemplate(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory(t)
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()

	emailTemplates, _ := policy.GetVerificationEmailTemplates()

	emailTemplate, err := GetEmailTemplate(emailTemplates.Items[0].Href)

	assert.NoError(t, err)
	assert.Equal(t, "Verify your account", emailTemplate.Subject)
}

func TestUpdateEmailTemplate(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory(t)
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()

	emailTemplates, _ := policy.GetVerificationEmailTemplates()

	emailTemplate, _ := GetEmailTemplate(emailTemplates.Items[0].Href)

	emailTemplate.Subject = "New Subject"
	err := emailTemplate.Update()

	assert.NoError(t, err)

	et, _ := GetEmailTemplate(emailTemplates.Items[0].Href)

	assert.Equal(t, "New Subject", et.Subject)
}

func TestRefreshEmailTemplate(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory(t)
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()

	emailTemplates, _ := policy.GetVerificationEmailTemplates()

	emailTemplate, _ := GetEmailTemplate(emailTemplates.Items[0].Href)

	emailTemplate.Subject = "New Subject"
	err := emailTemplate.Refresh()

	assert.NoError(t, err)
	assert.Equal(t, "Verify your account", emailTemplate.Subject)
}
