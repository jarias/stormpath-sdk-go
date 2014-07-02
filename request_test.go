package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {
	Context("Default values", func() {
		Describe("DontFollowRedirects", func() {
			It("should be false if not set", func() {
				r := StormpathRequest{}

				Expect(r.DontFollowRedirects).To(Equal(false))
			})
		})
	})

	Context("HTTP request", func() {
		Describe("Converting to an HTTP request", func() {
			It("should convert a stormpath request to a http.Request", func() {
				r := StormpathRequest{Method: "GET", URL: "http://test/test", PageRequest: PageRequest{}, Filter: DefaultFilter{}, Payload: []byte("")}
				req, err := r.ToHTTPRequest()

				Expect(err).NotTo(HaveOccurred())
				Expect(req.URL.Host).To(Equal("test"))
				Expect(req.URL.Path).To(Equal("/test"))
				Expect(req.Method).To(Equal(r.Method))
				Expect(req.URL.Query()).To(BeEmpty())
			})
		})
	})
})
