package stormpath

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetApplicationOAuthPolicy(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	policy, err := application.GetOAuthPolicy()

	assert.NoError(t, err)
	assert.Equal(t, "PT1H", policy.AccessTokenTTL)
	assert.Equal(t, "P60D", policy.RefreshTokenTTL)
}

func TestUpdateApplicationOAuthPolicy(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	policy, _ := application.GetOAuthPolicy()

	policy.AccessTokenTTL = "PT2H"
	policy.RefreshTokenTTL = "P50D"
	err := policy.Update()

	assert.NoError(t, err)

	updatedPolicy, _ := application.GetOAuthPolicy()

	assert.Equal(t, "PT2H", updatedPolicy.AccessTokenTTL)
	assert.Equal(t, "P50D", updatedPolicy.RefreshTokenTTL)
}

func TestUpdateApplicationOAuthPolicyInvalidTTL(t *testing.T) {
	t.Parallel()

	application := createTestApplication(t)
	defer application.Purge()

	policy, _ := application.GetOAuthPolicy()

	policy.AccessTokenTTL = "hello i'm not a valid value deal with me"
	err := policy.Update()

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(Error).Status)
	assert.Equal(t, 2002, err.(Error).Code)

	policy.AccessTokenTTL = "PT2H"
	policy.RefreshTokenTTL = "hello i'm not a valid value deal with me"
	err = policy.Update()

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(Error).Status)
	assert.Equal(t, 2002, err.(Error).Code)
}
