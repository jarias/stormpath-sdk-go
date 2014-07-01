package stormpath_test

import (
	"encoding/json"
	. "github.com/jarias/stormpath-sdk-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the account required fields", func() {
			acc := NewAccount("test@test.org", "123", "test", "test")

			jsonData, _ := json.Marshal(acc)

			Expect(string(jsonData)).To(Equal("{\"email\":\"test@test.org\",\"password\":\"123\",\"givenName\":\"test\",\"surname\":\"test\"}"))
		})
	})
})
