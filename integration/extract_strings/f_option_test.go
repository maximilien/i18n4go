package extract_strings_test

import (
	"io/ioutil"
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
			session := Runi18n("-extract-strings", "-v", "-p", "-f", filepath.Join(INPUT_FILES_PATH, "app.go"))
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

	Context("when the file specified has no strings at all", func() {
		var (
			OUTPUT_PATH string
		)

		BeforeEach(func() {
			var err error
			OUTPUT_PATH, err = ioutil.TempDir("", "gi18n4cf")
			Ω(err).ShouldNot(HaveOccurred())

			session := Runi18n("-extract-strings", "-f", filepath.Join(INPUT_FILES_PATH, "no_strings.go"), "-o", OUTPUT_PATH)
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("does not generate any files", func() {
			println(OUTPUT_PATH)
			files, err := ioutil.ReadDir(OUTPUT_PATH)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(files).Should(BeEmpty())
		})
	})
})
