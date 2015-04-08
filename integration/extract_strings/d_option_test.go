package extract_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -d dirName", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
		outputPath        string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		outputPath, err = ioutil.TempDir("", "i18n4go4go")
		Ω(err).ToNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "extract_strings")
		inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files", "quota")
		expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")
	})

	AfterEach(func() {
		os.RemoveAll(outputPath)
	})

	Context("When i18n4go4go is run with the -d flag", func() {
		BeforeEach(func() {
			session := Runi18n("-c", "extract-strings", "-v", "--po", "--meta", "-d", inputFilesPath, "-o", outputPath, "--ignore-regexp", "^[.]\\w+.go$")

			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("Walks input directory and compares each group of generated output to expected output", func() {
			filepath.Walk(inputFilesPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(expectedFilesPath, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
					filepath.Join(outputPath, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(expectedFilesPath, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
					filepath.Join(outputPath, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
				)

				return nil
			})
		})
	})

	Context("When i18n4go4go is run with the -d -r flags", func() {
		BeforeEach(func() {
			inputFilesPath = filepath.Join(inputFilesPath, "..")

			session := Runi18n("-c", "extract-strings", "-v", "--po", "--meta", "-d", inputFilesPath, "-o", outputPath, "-r", "--ignore-regexp", "^[.]\\w+.go$")
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("Walks input directories and compares each group of generated output to expected output", func() {
			filepath.Walk(inputFilesPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(expectedFilesPath, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
					filepath.Join(outputPath, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(expectedFilesPath, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
					filepath.Join(outputPath, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
				)

				return nil
			})
		})
	})
})
