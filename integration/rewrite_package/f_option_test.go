package rewrite_package_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rewrite-package -f filename", func() {
	var (
		rootPath          string
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	Context("all strings to rewrite are simple strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			session := Runi18n(
				"-rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test.go"),
				"-o", filepath.Join(rootPath, "tmp"),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("rewrites the input file with T() wrappers around strings", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
			bytes, err := ioutil.ReadFile(expectedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			expectedOutput := string(bytes)

			generatedOutputFile := filepath.Join(rootPath, "tmp", "test.go")
			bytes, err = ioutil.ReadFile(generatedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			actualOutput := string(bytes)

			Ω(actualOutput).Should(Equal(expectedOutput))
		})

		It("adds T func declaration and i18n init() func", func() {
			initFile := filepath.Join(rootPath, "tmp", "i18n_init.go")
			expectedBytes, err := ioutil.ReadFile(initFile)
			Ω(err).ShouldNot(HaveOccurred())
			expected := strings.TrimSpace(string(expectedBytes))

			expectedInitFile := filepath.Join(expectedFilesPath, "i18n_init.go")
			actualBytes, err := ioutil.ReadFile(expectedInitFile)
			Ω(err).ShouldNot(HaveOccurred())
			actual := strings.TrimSpace(string(actualBytes))

			Ω(actual).Should(Equal(expected))
		})
	})

	Context("strings to rewrite contain complex templated strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			session := Runi18n(
				"-rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test_templated_strings.go"),
				"-o", filepath.Join(rootPath, "tmp"),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("rewrites the input file with T() wrappers around all (simple and templated) strings", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test_templated_strings.go")
			bytes, err := ioutil.ReadFile(expectedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			expectedOutput := string(bytes)

			generatedOutputFile := filepath.Join(rootPath, "tmp", "test_templated_strings.go")
			bytes, err = ioutil.ReadFile(generatedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			actualOutput := string(bytes)

			Ω(actualOutput).Should(Equal(expectedOutput))
		})
	})

	Context("strings to rewrite contain interpolated strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			session := Runi18n(
				"-rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test_interpolated_strings.go"),
				"-o", filepath.Join(rootPath, "tmp"),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("converts interpolated strings to templated and rewrites the input file with T() wrappers around all (simple and templated) strings", func() {
			expectedOutputFile := filepath.Join(expectedFilesPath, "test_interpolated_strings.go")
			bytes, err := ioutil.ReadFile(expectedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			expectedOutput := string(bytes)

			generatedOutputFile := filepath.Join(rootPath, "tmp", "test_interpolated_strings.go")
			bytes, err = ioutil.ReadFile(generatedOutputFile)
			Ω(err).ShouldNot(HaveOccurred())

			actualOutput := string(bytes)

			Ω(actualOutput).Should(Equal(expectedOutput))
		})
	})
})
