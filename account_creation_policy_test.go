package stormpath_test

import (
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestGetVerificationEmailTemplatesPolicyNoExists(t *testing.T) {
	t.Parallel()

	policy := &AccountCreationPolicy{VerificationEmailTemplates: &EmailTemplates{}}
	policy.VerificationEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

	templates, err := policy.GetVerificationEmailTemplates()

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
	assert.Nil(t, templates)
}

func TestGetVerificationEmailTemplates(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()

	templates, err := policy.GetVerificationEmailTemplates()

	assert.NoError(t, err)
	assert.Len(t, templates.Items, 1)
}

func TestGetVerificationSuccessEmailTemplatesPolicyNoExists(t *testing.T) {
	t.Parallel()

	policy := &AccountCreationPolicy{VerificationSuccessEmailTemplates: &EmailTemplates{}}
	policy.VerificationSuccessEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationSuccessEmailTemplates"

	templates, err := policy.GetVerificationSuccessEmailTemplates()

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
	assert.Nil(t, templates)
}

func TestGetVerificationSuccessEmailTemplates(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()

	templates, err := policy.GetVerificationSuccessEmailTemplates()

	assert.NoError(t, err)
	assert.Len(t, templates.Items, 1)
}

func TestGetWelcomeEmailTemplatesPolicyNoExists(t *testing.T) {
	t.Parallel()

	policy := &AccountCreationPolicy{WelcomeEmailTemplates: &EmailTemplates{}}
	policy.WelcomeEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/welcomeEmailTemplates"

	templates, err := policy.GetWelcomeEmailTemplates()

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
	assert.Nil(t, templates)
}

func TestGetWelcomeEmailTemplates(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()

	templates, err := policy.GetWelcomeEmailTemplates()

	assert.NoError(t, err)
	assert.Len(t, templates.Items, 1)
}

func TestUpdateAccountCreationPolicy(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()
	policy.VerificationEmailStatus = Enabled
	err := policy.Update()

	assert.NoError(t, err)
	assert.Equal(t, Enabled, policy.VerificationEmailStatus)
}

func TestUpdateAccountCreationPolicyNoExists(t *testing.T) {
	t.Parallel()

	policy := AccountCreationPolicy{}
	policy.Href = GetClient().ClientConfiguration.BaseURL + "accountCreationPolicies/XXXX"

	err := policy.Update()

	assert.Error(t, err)
	assert.Equal(t, 404, err.(Error).Status)
}

func TestRefreshAccountCreationPolicy(t *testing.T) {
	t.Parallel()

	directory := createTestDirectory()
	defer directory.Delete()

	policy, _ := directory.GetAccountCreationPolicy()
	policy.VerificationEmailStatus = Enabled
	err := policy.Refresh()

	assert.NoError(t, err)
	assert.Equal(t, Disabled, policy.VerificationEmailStatus)
}
