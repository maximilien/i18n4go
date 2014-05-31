package extract_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -f fileName", func() {
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

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "extract_strings", "f_option")
		inputFilesPath = filepath.Join(fixturesPath, "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "expected_output")
	})

	Context("compare generated and expected file", func() {
		BeforeEach(func() {
			session := Runi18n("-extract-strings", "-v", "-p", "-meta", "-f", filepath.Join(inputFilesPath, "app.go"))
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			RemoveAllFiles(
				GetFilePath(inputFilesPath, "app.go.en.json"),
				GetFilePath(inputFilesPath, "app.go.en.po"),
				GetFilePath(inputFilesPath, "app.go.extracted.json"),
			)
		})

		It("app.go.en.json", func() {
			CompareExpectedToGeneratedTraslationJson(
				GetFilePath(expectedFilesPath, "app.go.en.json"),
				GetFilePath(inputFilesPath, "app.go.en.json"),
			)
		})

		It("app.go.extracted.json", func() {
			CompareExpectedToGeneratedExtendedJson(
				GetFilePath(expectedFilesPath, "app.go.extracted.json"),
				GetFilePath(inputFilesPath, "app.go.extracted.json"),
			)
		})

		It("app.go.en.po", func() {
			CompareExpectedToGeneratedPo(
				GetFilePath(expectedFilesPath, "app.go.en.po"),
				GetFilePath(inputFilesPath, "app.go.en.po"),
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

			session := Runi18n("-extract-strings", "-f", filepath.Join(inputFilesPath, "no_strings.go"), "-o", OUTPUT_PATH)
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
