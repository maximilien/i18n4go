package rewrite_package_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rewrite-package --i18n-strings-filename some-file", func() {
	var (
		outputDir         string
		rootPath          string
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	AfterEach(func() {
		err := os.RemoveAll(outputDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("input file only contains simple strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "i18n_strings_filename_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "i18n_strings_filename_option", "expected_output")

			session := Runi18n("-c",
				"rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test.go"),
				"-o", outputDir,
				"--i18n-strings-filename", filepath.Join(inputFilesPath, "strings.json"),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("rewrites the input file with T() wrappers around the strings specified in the --i18n-strings-filename flag", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
			bytes, err := ioutil.ReadFile(expectedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			expectedOutput := string(bytes)

			generatedOutputFile := filepath.Join(outputDir, "test.go")
			bytes, err = ioutil.ReadFile(generatedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			actualOutput := string(bytes)

			Ω(actualOutput).Should(Equal(expectedOutput))
		})
	})

	Context("input file contains some templated strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "i18n_strings_filename_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "i18n_strings_filename_option", "expected_output")

			session := Runi18n("-c",
				"rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test_templated_strings.go"),
				"-o", outputDir,
				"--i18n-strings-filename", filepath.Join(inputFilesPath, "test_templated_strings.go.en.json"),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("rewrites the input file with T() wrappers around the strings (templated and not) specified in the --i18n-strings-filename flag", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test_templated_strings.go")
			bytes, err := ioutil.ReadFile(expectedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			expectedOutput := string(bytes)

			generatedOutputFile := filepath.Join(outputDir, "test_templated_strings.go")
			bytes, err = ioutil.ReadFile(generatedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			actualOutput := string(bytes)

			Ω(actualOutput).Should(Equal(expectedOutput))
		})
	})

	Context("input file contains some interpolated strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "i18n_strings_filename_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "i18n_strings_filename_option", "expected_output")

			CopyFile(filepath.Join(inputFilesPath, "_test_interpolated_strings.go.en.json"), filepath.Join(inputFilesPath, "test_interpolated_strings.go.en.json"))

			session := Runi18n("-c",
				"rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test_interpolated_strings.go"),
				"-o", outputDir,
				"--i18n-strings-filename", filepath.Join(inputFilesPath, "test_interpolated_strings.go.en.json"),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			CopyFile(filepath.Join(inputFilesPath, "_test_interpolated_strings.go.en.json"), filepath.Join(inputFilesPath, "test_interpolated_strings.go.en.json"))
		})

		It("converts interpolated strings to templated and rewrites the input file with T() wrappers around the strings (templated and not) specified in the --i18n-strings-filename flag", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test_interpolated_strings.go")
			bytes, err := ioutil.ReadFile(expectedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			expectedOutput := string(bytes)

			generatedOutputFile := filepath.Join(outputDir, "test_interpolated_strings.go")
			bytes, err = ioutil.ReadFile(generatedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			actualOutput := string(bytes)

			Ω(actualOutput).Should(Equal(expectedOutput))
		})

		It("updates the i18n strings JSON file with the converted interpolated JSON strings", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test_interpolated_strings.go.en.json")
			generatedOutputFile := filepath.Join(inputFilesPath, "test_interpolated_strings.go.en.json")
			CompareExpectedToGeneratedTraslationJson(expectedOutputFile, generatedOutputFile)
		})
	})
})
