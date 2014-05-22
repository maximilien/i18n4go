package extract_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -f fileName -o outputDir", func() {
	var (
		INPUT_FILES_PATH  = filepath.Join("f_option", "input_files")
		EXPECTED_DIR_PATH = filepath.Join("f_option", "expected_output")
		OUTPUT_PATH       string
	)

	Context("-o outputDir --output-flat (default)", func() {
		BeforeEach(func() {
			var err error
			OUTPUT_PATH, err = ioutil.TempDir("", "gi18n4cf")
			Ω(err).ToNot(HaveOccurred())

			session := Runi18n("-extract-strings", "-v", "-p", "-f", filepath.Join(INPUT_FILES_PATH, "app.go"), "-o", OUTPUT_PATH)
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			os.RemoveAll(OUTPUT_PATH)
		})

		It("Walks input directory and compares each group of generated output to expected output", func() {

			CompareExpectedToGeneratedTraslationJson(
				filepath.Join(EXPECTED_DIR_PATH, "app.go.en.json"),
				filepath.Join(OUTPUT_PATH, "app.go.en.json"),
			)

			CompareExpectedToGeneratedExtendedJson(
				filepath.Join(EXPECTED_DIR_PATH, "app.go.extracted.json"),
				filepath.Join(OUTPUT_PATH, "app.go.extracted.json"),
			)

			CompareExpectedToGeneratedPo(
				filepath.Join(EXPECTED_DIR_PATH, "app.go.en.po"),
				filepath.Join(OUTPUT_PATH, "app.go.en.po"),
			)
		})
	})

	Context("-o outputDir --output-match-package", func() {
		BeforeEach(func() {
			var err error
			OUTPUT_PATH, err = ioutil.TempDir("", "gi18n4cf")
			Ω(err).ToNot(HaveOccurred())

			session := Runi18n("-extract-strings", "-v", "-p", "-f", filepath.Join(INPUT_FILES_PATH, "app.go"), "-o", OUTPUT_PATH, "-output-match-package")
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			os.RemoveAll(OUTPUT_PATH)
		})

		It("Walks input directory and compares each group of generated output to expected output and package subdirectories", func() {
			EXPECTED_DIR_PATH = filepath.Join(EXPECTED_DIR_PATH, "app")
			OUTPUT_PATH = filepath.Join(OUTPUT_PATH, "app")
			CompareExpectedToGeneratedTraslationJson(
				filepath.Join(EXPECTED_DIR_PATH, "app.go.en.json"),
				filepath.Join(OUTPUT_PATH, "app.go.en.json"),
			)

			CompareExpectedToGeneratedExtendedJson(
				filepath.Join(EXPECTED_DIR_PATH, "app.go.extracted.json"),
				filepath.Join(OUTPUT_PATH, "app.go.extracted.json"),
			)

			CompareExpectedToGeneratedPo(
				filepath.Join(EXPECTED_DIR_PATH, "app.go.en.po"),
				filepath.Join(OUTPUT_PATH, "app.go.en.po"),
			)
		})
	})
})
