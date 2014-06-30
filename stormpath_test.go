package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

var _ = Describe("Stormpath SAuthc1", func() {
	It("should authenticate a request without query params", func() {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/", nil)

		cred := &Credentials{Id: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte(""), time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

		Expect(req.Header.Get("Authorization")).To(Equal("SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
			"sauthc1SignedHeaders=host;x-stormpath-date, " +
			"sauthc1Signature=990a95aabbcbeb53e48fb721f73b75bd3ae025a2e86ad359d08558e1bbb9411c"))
	})

	It("should authenticate a request with query params", func() {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/directories?orderBy=name+asc", nil)

		cred := &Credentials{Id: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte(""), time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

		Expect(req.Header.Get("Authorization")).To(Equal("SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
			"sauthc1SignedHeaders=host;x-stormpath-date, " +
			"sauthc1Signature=fc04c5187cc017bbdf9c0bb743a52a9487ccb91c0996267988ceae3f10314176"))
	})
})
