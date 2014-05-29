package merge_strings_test

import (
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("merge-strings -d dirName", func() {
	var (
		rootPath          string
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		dir, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())
		rootPath = filepath.Join(dir, "..", "..")

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "merge_strings")
		inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")
	})

	Context("can combine multiple language files", func() {
		BeforeEach(func() {
			session := Runi18n("-merge-strings", "-v", "-d", filepath.Join(inputFilesPath), "-source-language", "en")
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			RemoveAllFiles(
				GetFilePath(inputFilesPath, "en.all.json"),
			)
		})

		It("en.all.json contains translations from both files", func() {
			CompareExpectedToGeneratedTraslationJson(
				GetFilePath(expectedFilesPath, "en.all.json"),
				GetFilePath(inputFilesPath, "en.all.json"),
			)
		})
	})

})
