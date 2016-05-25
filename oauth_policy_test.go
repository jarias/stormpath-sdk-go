package stormpath_test

import (
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestGetApplicationOAuthPolicy(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	policy, err := application.GetOAuthPolicy()

	assert.NoError(t, err)
	assert.Equal(t, "PT1H", policy.AccessTokenTtl)
	assert.Equal(t, "P60D", policy.RefreshTokenTtl)
}

func TestUpdateApplicationOAuthPolicy(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	policy, _ := application.GetOAuthPolicy()

	policy.AccessTokenTtl = "PT2H"
	policy.RefreshTokenTtl = "P50D"
	err := policy.Update()

	assert.NoError(t, err)

	updatedPolicy, _ := application.GetOAuthPolicy()

	assert.Equal(t, "PT2H", updatedPolicy.AccessTokenTtl)
	assert.Equal(t, "P50D", updatedPolicy.RefreshTokenTtl)
}

func TestUpdateApplicationOAuthPolicyInvalidTTL(t *testing.T) {
	t.Parallel()

	application := createTestApplication()
	defer application.Purge()

	policy, _ := application.GetOAuthPolicy()

	policy.AccessTokenTtl = "hello i'm not a valid value deal with me"
	err := policy.Update()
	
	assert.Error(t, err)
	assert.Equal(t, 400, err.(Error).Status)
	assert.Equal(t, 2002, err.(Error).Code)
	
	policy.AccessTokenTtl = "PT2H"
	policy.RefreshTokenTtl = "hello i'm not a valid value deal with me"
	err = policy.Update()
	
	assert.Error(t, err)
	assert.Equal(t, 400, err.(Error).Status)
	assert.Equal(t, 2002, err.(Error).Code)
}
