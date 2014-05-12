package extract_strings_test

import (
	"github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestExtractStrings(t *testing.T) {
	BeforeSuite(test_helpers.BuildExecutable)

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "ExtractStrings Suite")
}
