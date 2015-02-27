package stormpath_test

import (
	"net/http"
	"time"

	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stormpath SAuthc1", func() {
	It("should authenticate a request without query params", func() {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/", nil)

		cred := Credentials{ID: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

		Expect(req.Header.Get("Authorization")).To(Equal("SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
			"sauthc1SignedHeaders=host;x-stormpath-date, " +
			"sauthc1Signature=990a95aabbcbeb53e48fb721f73b75bd3ae025a2e86ad359d08558e1bbb9411c"))
	})

	It("should authenticate a request with query params", func() {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/directories?orderBy=name+asc", nil)

		cred := Credentials{ID: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

		Expect(req.Header.Get("Authorization")).To(Equal("SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
			"sauthc1SignedHeaders=host;x-stormpath-date, " +
			"sauthc1Signature=fc04c5187cc017bbdf9c0bb743a52a9487ccb91c0996267988ceae3f10314176"))
	})

	It("should authenticate a request with multiple query params", func() {
		req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/applications/77JnfFiREjdfQH0SObMfjI/groups?q=group&limit=25&offset=25", nil)

		cred := Credentials{ID: "MyId", Secret: "Shush!"}

		Authenticate(req, []byte{}, time.Date(2013, 7, 1, 0, 0, 0, 0, time.UTC), cred, "a43a9d25-ab06-421e-8605-33fd1e760825")

		Expect(req.Header.Get("Authorization")).To(Equal("SAuthc1 sauthc1Id=MyId/20130701/a43a9d25-ab06-421e-8605-33fd1e760825/sauthc1_request, " +
			"sauthc1SignedHeaders=host;x-stormpath-date, " +
			"sauthc1Signature=e30a62c0d03ca6cb422e66039786865f3eb6269400941ede6226760553a832d3"))
	})

	//Describe("https://github.com/stormpath/stormpath-sdk-python/issues/101", func() {
	//	It("should authenticate a paginated request of groups", func() {
	//		directory := NewDirectory("new-directory-test")
	//		tenant.CreateDirectory(directory)
	//
	//		for i := 1; i <= 50; i++ {
	//			group := NewGroup(fmt.Sprintf("group-%d", i))
	//			directory.CreateGroup(group)
	//		}
	//
	//		filter := url.Values{}
	//		filter.Add("q", "group")
	//
	//		groups, err := directory.GetGroups(NewPageRequest(25, 25), filter)
	//
	//		Expect(err).NotTo(HaveOccurred())
	//		Expect(groups.Href).NotTo(BeEmpty())
	//		Expect(groups.Offset).To(Equal(25))
	//		Expect(groups.Limit).To(Equal(25))
	//		Expect(groups.Items).NotTo(BeEmpty())
	//		directory.Delete()
	//	})
	//})
})
