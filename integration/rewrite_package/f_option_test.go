package rewrite_package_test

import (
	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ = Describe("rewrite-package -f filename", func() {
	var (
		rootPath            = "."
		INPUT_FILES_PATH    = filepath.Join("f_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
	)

	BeforeEach(func() {
		dir, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())
		rootPath = filepath.Join(dir, "..", "..")

		session := Runi18n(
			"-rewrite-package",
			"-f", filepath.Join(INPUT_FILES_PATH, "test.go"),
			"-o", filepath.Join(rootPath, "tmp"),
			"-v",
		)

		Ω(session.ExitCode()).Should(Equal(0))
	})

	It("rewrites the input file with T() wrappers around strings", func() {
		expectedOutputFile := filepath.Join(EXPECTED_FILES_PATH, "test.go")
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

		expectedInitFile := filepath.Join(EXPECTED_FILES_PATH, "i18n_init.go")
		actualBytes, err := ioutil.ReadFile(expectedInitFile)
		Ω(err).ShouldNot(HaveOccurred())
		actual := strings.TrimSpace(string(actualBytes))

		Ω(actual).Should(Equal(expected))
	})
})
