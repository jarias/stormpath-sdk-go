package stormpath_test

import (
	. "github.com/jarias/stormpath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {
	Context("Paging", func() {
		Describe("Converting to query params", func() {
			It("should convert a page request to URL query params", func() {
				pr := PageRequest{Offset: 2, Limit: 2}
				q := pr.ToUrlQueryValues()

				Expect(q.Encode()).To(Equal("limit=2&offset=2"))
			})
		})
	})

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
				r := StormpathRequest{Method: "GET", URL: "http://test/test", PageRequest: &PageRequest{}, Filter: DefaultFilter{}, Payload: []byte("")}
				req, err := r.ToHttpRequest()

				Expect(err).NotTo(HaveOccurred())
				Expect(req.URL.Host).To(Equal("test"))
				Expect(req.URL.Path).To(Equal("/test"))
				Expect(req.Method).To(Equal(r.Method))
				Expect(req.URL.Query()).To(BeEmpty())
			})
		})
	})
})
