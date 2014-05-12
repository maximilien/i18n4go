package extract_strings_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -f fileName", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("f_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
	)

	Context("compare generated and expected file", func() {
		BeforeEach(func() {
			session := Runi18n("-extract-strings", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "app.go"))
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			RemoveAllFiles(
				GetFilePath(INPUT_FILES_PATH, "app.go.en.json"),
				GetFilePath(INPUT_FILES_PATH, "app.go.en.po"),
				GetFilePath(INPUT_FILES_PATH, "app.go.extracted.json"),
			)
		})

		It("app.go.en.json", func() {
			CompareExpectedToGeneratedTraslationJson(
				GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json"),
				GetFilePath(INPUT_FILES_PATH, "app.go.en.json"),
			)
		})

		It("app.go.extracted.json", func() {
			CompareExpectedToGeneratedExtendedJson(
				GetFilePath(EXPECTED_FILES_PATH, "app.go.extracted.json"),
				GetFilePath(INPUT_FILES_PATH, "app.go.extracted.json"),
			)
		})

		It("app.go.en.po", func() {
			CompareExpectedToGeneratedPo(
				GetFilePath(EXPECTED_FILES_PATH, "app.go.en.po"),
				GetFilePath(INPUT_FILES_PATH, "app.go.en.po"),
			)
		})

	})
})
