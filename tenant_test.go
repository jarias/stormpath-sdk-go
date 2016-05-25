package stormpath_test

import (
	"testing"

	. "github.com/jarias/stormpath-sdk-go"
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
	application := newTestApplication()
	defer application.Purge()

	for i := 0; i < b.N; i++ {
		err := tenant.CreateApplication(application)
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

func TestTenantCreateApplication(t *testing.T) {
	t.Parallel()

	application := newTestApplication()
	defer application.Purge()

	err := tenant.CreateApplication(application)

	assert.NoError(t, err)
	assert.NotEmpty(t, application.Href)
}

func TestTenantGetApplications(t *testing.T) {
	t.Parallel()

	applications, err := tenant.GetApplications(MakeApplicationsCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, applications.Items)
	assert.NotEmpty(t, applications.Href)
	assert.Equal(t, 0, applications.Offset)
	assert.Equal(t, 25, applications.Limit)
}

func TestTenantGetApplicationsByPage(t *testing.T) {
	t.Parallel()

	applications, err := tenant.GetApplications(MakeApplicationsCriteria().Limit(1))

	assert.NoError(t, err)
	assert.Len(t, applications.Items, 1)
	assert.NotEmpty(t, applications.Href)
	assert.Equal(t, 0, applications.Offset)
	assert.Equal(t, 1, applications.Limit)
}

func TestTenantGetApplicationsFiltered(t *testing.T) {
	t.Parallel()

	applications, err := tenant.GetApplications(MakeApplicationsCriteria().NameEq("stormpath"))

	assert.NoError(t, err)
	assert.Len(t, applications.Items, 1)
	assert.NotEmpty(t, applications.Href)
	assert.Equal(t, 0, applications.Offset)
	assert.Equal(t, 25, applications.Limit)

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
	t.Parallel()

	directories, err := tenant.GetDirectories(MakeDirectoriesCriteria())

	assert.NoError(t, err)
	assert.NotEmpty(t, directories.Items)
	assert.NotEmpty(t, directories.Href)
	assert.Equal(t, 0, directories.Offset)
	assert.Equal(t, 25, directories.Limit)
}

func TestTenantGetDirectoriesByPage(t *testing.T) {
	t.Parallel()

	directories, err := tenant.GetDirectories(MakeDirectoriesCriteria().Limit(1))

	assert.NoError(t, err)
	assert.NotEmpty(t, directories.Items)
	assert.NotEmpty(t, directories.Href)
	assert.Equal(t, 0, directories.Offset)
	assert.Equal(t, 1, directories.Limit)
}

func TestTenantGetDirectoriesFiltered(t *testing.T) {
	t.Parallel()

	directories, err := tenant.GetDirectories(MakeDirectoriesCriteria().NameEq("Stormpath Administrators"))

	assert.NoError(t, err)
	assert.NotEmpty(t, directories.Items)
	assert.NotEmpty(t, directories.Href)
	assert.Equal(t, 0, directories.Offset)
	assert.Equal(t, 25, directories.Limit)

	err = directories.Items[0].Refresh()

	assert.NoError(t, err)
	assert.Equal(t, "Stormpath Administrators", directories.Items[0].Name)
}
