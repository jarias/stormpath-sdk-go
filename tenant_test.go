package stormpath

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGetCurrentTenant(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := CurrentTenant()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkCreateApplication(b *testing.B) {
	for i := 0; i < b.N; i++ {
		application := newTestApplication()
		err := tenant.CreateApplication(application)
		defer application.Purge()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUpdateCustomData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		customData := map[string]interface{}{
			"testIntField":    1,
			"testStringField": "test",
		}

		_, err := tenant.UpdateCustomData(customData)
		if err != nil {
			panic(err)
		}
	}
}

func ExampleCurrentTenant() {
	tenant, err := CurrentTenant()
	if err != nil {
		log.Panicf("Couldn't get the current tenant %s", err)
	}

	fmt.Printf("tenant = %+v\n", tenant)
}

func TestGetCurrentTenant(t *testing.T) {
	t.Parallel()

	currentTenant, err := CurrentTenant()

	assert.NoError(t, err)
	assert.NotEmpty(t, currentTenant.Href)
	assert.NotEmpty(t, currentTenant.Name)
	assert.NotEmpty(t, currentTenant.Key)
	assert.NotEmpty(t, currentTenant.Applications.Href)
	assert.NotEmpty(t, currentTenant.Directories.Href)
}

//TODO: Fix panic: runtime error: invalid memory address or nil pointer dereference
// This error is caused by Purge() method when there were errors during application creation.
func TestTenantCreateApplication(t *testing.T) {
	t.Parallel()

	application := newTestApplication()
	defer application.Purge()

	err := tenant.CreateApplication(application)

	assert.NoError(t, err)
	assert.NotEmpty(t, application.Href)
}

func TestTenantGetApplications(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")
	t.Parallel()

	applications, err := tenant.GetApplications(MakeApplicationsCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, applications.Items)
	assert.NotEmpty(t, applications.Href)
	assert.Equal(t, 0, applications.GetOffset())
	assert.Equal(t, 25, applications.GetLimit())
}

func TestTenantGetApplicationsByPage(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")
	t.Parallel()

	applications, err := tenant.GetApplications(MakeApplicationsCriteria().Limit(1))

	assert.NoError(t, err)
	assert.Len(t, applications.Items, 1)
	assert.NotEmpty(t, applications.Href)
	assert.Equal(t, 0, applications.GetOffset())
	assert.Equal(t, 1, applications.GetLimit())
}

func TestTenantGetApplicationsFiltered(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")
	t.Parallel()

	applications, err := tenant.GetApplications(MakeApplicationsCriteria().NameEq("stormpath"))

	assert.NoError(t, err)
	assert.Len(t, applications.Items, 1)
	assert.NotEmpty(t, applications.Href)
	assert.Equal(t, 0, applications.GetOffset())
	assert.Equal(t, 25, applications.GetLimit())

	err = applications.Items[0].Refresh()

	assert.NoError(t, err)
	assert.Equal(t, "Stormpath", applications.Items[0].Name)
}

func TestUpdateTenantCustomData(t *testing.T) {
	customData := map[string]interface{}{
		"testIntField":    1,
		"testStringField": "test",
	}

	updatedCustomData, err := tenant.UpdateCustomData(customData)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), updatedCustomData["testIntField"])
	assert.Len(t, updatedCustomData, 5)

	//Clean the tenant custom data
	tenant.DeleteCustomData()
}

func TestGetTenantCustomData(t *testing.T) {
	customData := map[string]interface{}{
		"testIntField":    1,
		"testStringField": "test",
	}

	tenant.UpdateCustomData(customData)

	updatedCustomData, err := tenant.GetCustomData()

	assert.NoError(t, err)
	assert.Equal(t, float64(1), updatedCustomData["testIntField"])
	assert.Len(t, updatedCustomData, 5)

	//Clean the tenant custom data
	tenant.DeleteCustomData()
}

func TestDeleteTenantCustomData(t *testing.T) {
	customData := map[string]interface{}{
		"testIntField":    1,
		"testStringField": "test",
	}

	tenant.UpdateCustomData(customData)

	err := tenant.DeleteCustomData()
	assert.NoError(t, err)

	updatedCustomData, err := tenant.GetCustomData()

	assert.NoError(t, err)
	assert.Len(t, updatedCustomData, 3)
}

func TestTenantCreateDirectory(t *testing.T) {
	t.Parallel()

	dir := newTestDirectory()
	defer dir.Delete()

	err := tenant.CreateDirectory(dir)

	assert.NoError(t, err)
	assert.NotEmpty(t, dir.Href)
	assert.NotEmpty(t, dir.Name)
}

func TestTenantGetDirectories(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")

	t.Parallel()

	directories, err := tenant.GetDirectories(MakeDirectoriesCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, directories.Items)
	assert.NotEmpty(t, directories.Href)
	assert.Equal(t, 0, directories.GetOffset())
	assert.Equal(t, 25, directories.GetLimit())
}

func TestTenantGetDirectoriesByPage(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")

	t.Parallel()

	directories, err := tenant.GetDirectories(MakeDirectoriesCriteria().Limit(1))

	assert.NoError(t, err)
	assert.NotEmpty(t, directories.Items)
	assert.NotEmpty(t, directories.Href)
	assert.Equal(t, 0, directories.GetOffset())
	assert.Equal(t, 1, directories.GetLimit())
}

func TestTenantGetDirectoriesFiltered(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")

	t.Parallel()

	directories, err := tenant.GetDirectories(MakeDirectoriesCriteria().NameEq("Stormpath Administrators"))

	assert.NoError(t, err)
	assert.NotEmpty(t, directories.Items)
	assert.NotEmpty(t, directories.Href)
	assert.Equal(t, 0, directories.GetOffset())
	assert.Equal(t, 25, directories.GetLimit())

	err = directories.Items[0].Refresh()

	assert.NoError(t, err)
	assert.Equal(t, "Stormpath Administrators", directories.Items[0].Name)
}

/*
func TestTenantGetAccountsByCustomData(t *testing.T) {
	t.Parallel()

	cdKey := "customId"
	cdValue := "myCustomDataValue"
	application := createTestApplication(t)
	defer application.Purge()

	account := createTestAccount(application, t)
	customData, err := account.UpdateCustomData(map[string]interface{}{cdKey: cdValue})

	assert.NoError(t, err)
	assert.Equal(t, cdValue, customData[cdKey])

	time.Sleep(5 * time.Second)

	accounts, err := tenant.GetAccounts(MakeAccountCriteria().CustomDataEq(cdKey, cdValue))

	assert.Len(t, accounts.Items, 1)
}
*/

func TestTenantGetOrganizations(t *testing.T) {
	t.Skip("Skiping until I figure out the issue with this API call")

	t.Parallel()

	organizations, err := tenant.GetOrganizations(MakeOrganizationsCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, organizations.Items)
	assert.NotEmpty(t, organizations.Href)
	assert.Equal(t, 0, organizations.GetOffset())
	assert.Equal(t, 25, organizations.GetLimit())
}
