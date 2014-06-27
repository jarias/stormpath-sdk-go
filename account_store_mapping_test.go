package stormpath_test

import (
	"encoding/json"
	. "github.com/jarias/stormpath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountStoreMapping", func() {
	Describe("JSON", func() {
		It("should marshall a minimum JSON with only the required fields", func() {
			accountStoreMapping := NewAccountStoreMapping("http://appurl", "http://storeUrl")

			jsonData, _ := json.Marshal(accountStoreMapping)

			Expect(string(jsonData)).To(Equal("{\"application\":{\"href\":\"http://appurl\"},\"accountStore\":{\"href\":\"http://storeUrl\"}}"))
		})
	})

	Describe("Save", func() {
		PIt("should create a new account store mapping", func() {
		})
	})
})
