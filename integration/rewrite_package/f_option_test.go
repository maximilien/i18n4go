package rewrite_package_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rewrite-package -f filename", func() {
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

	Context("no -o option passed, so input file is rewritten", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			CopyFile(filepath.Join(inputFilesPath, "test.go"), filepath.Join(outputDir, "test.go"))

			session := Runi18n("-c", "rewrite-package",
				"-f", filepath.Join(outputDir, "test.go"),
				"--root-path", outputDir,
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("overwrites the input file with T() wrappers around strings", func() {
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

		It("adds T func declaration and i18n init() func in the same directory as input file", func() {
			initFile := filepath.Join(outputDir, "i18n_init.go")
			expectedBytes, err := ioutil.ReadFile(initFile)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(expectedBytes).ShouldNot(Equal(""))
		})
	})

	Context("all strings to rewrite are simple strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			session := Runi18n("-c",
				"rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test.go"),
				"-o", outputDir,
				"--root-path", rootPath,
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("rewrites the input file with T() wrappers around strings", func() {
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

		It("adds T func declaration and i18n init() func", func() {
			initFile := filepath.Join(outputDir, "i18n_init.go")
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

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			session := Runi18n("-c",
				"rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test_templated_strings.go"),
				"-o", filepath.Join(outputDir),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("rewrites the input file with T() wrappers around all (simple and templated) strings", func() {
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

	Context("strings to rewrite contain interpolated strings", func() {
		BeforeEach(func() {
			dir, err := os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			rootPath = filepath.Join(dir, "..", "..")

			outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
			Ω(err).ShouldNot(HaveOccurred())

			fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
			inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
			expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

			session := Runi18n("-c",
				"rewrite-package",
				"-f", filepath.Join(inputFilesPath, "test_interpolated_strings.go"),
				"-o", filepath.Join(outputDir),
				"-v",
			)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("converts interpolated strings to templated and rewrites the input file with T() wrappers around all (simple and templated) strings", func() {
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
	})
})
