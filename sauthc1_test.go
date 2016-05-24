package stormpath_test

import (
	"net/http"
	"testing"
	"time"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSAuthc1WithoutQueryParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/", nil)

		cred := Credentials{ID: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")
	}
}

func BenchmarkSAuthc1WithQueryParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/directories?orderBy=name+asc", nil)

		cred := Credentials{ID: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")
	}
}

func BenchmarkSAuthc1WithMultipleQueryParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/applications/77JnfFiREjdfQH0SObMfjI/groups?q=group&limit=25&offset=25", nil)

		cred := Credentials{ID: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")
	}
}

func TestSAuthc1WithoutQueryParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/", nil)

	cred := Credentials{ID: "MyId", Secret: "Shush!"}

	Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

	expectedAuthHeaderValue := "SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
		"sauthc1SignedHeaders=host;x-stormpath-date, " +
		"sauthc1Signature=990a95aabbcbeb53e48fb721f73b75bd3ae025a2e86ad359d08558e1bbb9411c"
	authHeader := req.Header.Get("Authorization")
	assert.Equal(t, expectedAuthHeaderValue, authHeader)
}

func TestSAuthc1WithQueryParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/directories?orderBy=name+asc", nil)

	cred := Credentials{ID: "MyId", Secret: "Shush!"}

	Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

	expectedAuthHeaderValue := "SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
		"sauthc1SignedHeaders=host;x-stormpath-date, " +
		"sauthc1Signature=fc04c5187cc017bbdf9c0bb743a52a9487ccb91c0996267988ceae3f10314176"
	authHeader := req.Header.Get("Authorization")
	assert.Equal(t, expectedAuthHeaderValue, authHeader)
}

func TestSAuthc1WithMultipleQueryParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/applications/77JnfFiREjdfQH0SObMfjI/groups?q=group&limit=25&offset=25", nil)

	cred := Credentials{ID: "MyId", Secret: "Shush!"}

	Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

	expectedAuthHeaderValue := "SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
		"sauthc1SignedHeaders=host;x-stormpath-date, " +
		"sauthc1Signature=e30a62c0d03ca6cb422e66039786865f3eb6269400941ede6226760553a832d3"
	authHeader := req.Header.Get("Authorization")
	assert.Equal(t, expectedAuthHeaderValue, authHeader)
}
