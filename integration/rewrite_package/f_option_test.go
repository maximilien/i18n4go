package rewrite_package_test

import (
	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("rewrite-package -f filename", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("f_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
	)

	It("adds an import statement for go-i18n and i18n ", func() {
		dir, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())
		dir = filepath.Join(dir, "..", "..")

		session := Runi18n(
			"-rewrite-package",
			"-f", filepath.Join(INPUT_FILES_PATH, "test.go"),
			"-o", filepath.Join(dir, "tmp"),
			"-v",
		)

		Ω(session.ExitCode()).Should(Equal(0))

		expectedOutputFile := filepath.Join(EXPECTED_FILES_PATH, "test.go")
		bytes, err := ioutil.ReadFile(expectedOutputFile)
		Ω(err).ShouldNot(HaveOccurred())

		expectedOutput := string(bytes)

		generatedOutputFile := filepath.Join(dir, "tmp", "test.go")
		bytes, err = ioutil.ReadFile(generatedOutputFile)
		Ω(err).ShouldNot(HaveOccurred())

		actualOutput := string(bytes)

		Ω(actualOutput).Should(Equal(expectedOutput))
	})
})
