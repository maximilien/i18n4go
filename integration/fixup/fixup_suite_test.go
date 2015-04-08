package fixup_test

import (
	"testing"

	"github.com/maximilien/i18n4go/integration/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCheckup(t *testing.T) {
	BeforeSuite(test_helpers.BuildExecutable)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fixup Suite")
}
