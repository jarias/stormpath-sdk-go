package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filters", func() {
	Describe("NewDefaultFilter", func() {
		It("should create a url.Values map with the given field values", func() {
			filter := NewDefaultFilter("name", "description", "status")

			Expect(filter.Get("name")).To(Equal("name"))
			Expect(filter.Get("description")).To(Equal("description"))
			Expect(filter.Get("status")).To(Equal("status"))
		})
	})

	Describe("NewAccountFilter", func() {
		It("should create a url.Values map with the given field values", func() {
			filter := NewAccountFilter("givenName", "middleName", "surname", "username", "email")

			Expect(filter.Get("givenName")).To(Equal("givenName"))
			Expect(filter.Get("middleName")).To(Equal("middleName"))
			Expect(filter.Get("surname")).To(Equal("surname"))
			Expect(filter.Get("username")).To(Equal("username"))
			Expect(filter.Get("email")).To(Equal("email"))
		})
	})

	Describe("NewEmptyFilter", func() {
		It("should return an empty url.Values", func() {
			filter := NewEmptyFilter()

			Expect(filter).To(HaveLen(0))
		})
	})
})
