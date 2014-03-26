package stormpath_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStormpath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stormpath Suite")
}
