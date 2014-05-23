package merge_strings_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("merge-strings -d dirName", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("d_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("d_option", "expected_output")
	)

	Context("can combine multiple language files", func() {
		BeforeEach(func() {
			session := Runi18n("-merge-strings", "-v", "-d", filepath.Join(INPUT_FILES_PATH), "-source-language", "en")
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			RemoveAllFiles(
				GetFilePath(INPUT_FILES_PATH, "en.all.json"),
			)
		})

		It("en.all.json contains translations from both files", func() {
			CompareExpectedToGeneratedTraslationJson(
				GetFilePath(EXPECTED_FILES_PATH, "en.all.json"),
				GetFilePath(INPUT_FILES_PATH, "en.all.json"),
			)
		})
	})

})
