package extract_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -s filePath", func() {
	var (
		outputPath            string
		fixturesPath          string
		inputFilesPath        string
		expectedFilesPath     string
		matchingGroupFilePath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		outputPath, err = ioutil.TempDir("", "i18n4go4go")
		Ω(err).ToNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "extract_strings")
		inputFilesPath = filepath.Join(fixturesPath, "s_option", "input_files", "app", "app.go")
		expectedFilesPath = filepath.Join(fixturesPath, "s_option", "expected_output")
		matchingGroupFilePath = filepath.Join(fixturesPath, "s_option", "input_files", "matching_group.json")
	})

	AfterEach(func() {
		os.RemoveAll(outputPath)
	})

	Context("When i18n4go4go is run with the -s flag", func() {
		BeforeEach(func() {
			session := Runi18n("-c", "extract-strings", "-v", "-f", inputFilesPath, "-o", outputPath, "-s", matchingGroupFilePath)

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("use the regexFile to parse substring", func() {
			expectedFilePath := filepath.Join(expectedFilesPath, "app.go.en.json")
			actualFilePath := filepath.Join(outputPath, "app.go.en.json")

			expectedBytes, err := ioutil.ReadFile(expectedFilePath)
			Ω(err).Should(BeNil())
			Ω(expectedBytes).ShouldNot(BeNil())

			actualBytes, err := ioutil.ReadFile(actualFilePath)
			Ω(err).Should(BeNil())
			Ω(actualBytes).ShouldNot(BeNil())

			Ω(string(expectedBytes)).Should(MatchJSON(string(actualBytes)))
		})
	})
})
