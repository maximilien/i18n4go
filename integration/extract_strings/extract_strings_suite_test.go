package extract_strings_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var gi18nExec string

func TestExtractStrings(t *testing.T) {
	BeforeSuite(func() {
		var err error
		gi18nExec, err = gexec.Build("./../../main")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "ExtractStrings Suite")
}
