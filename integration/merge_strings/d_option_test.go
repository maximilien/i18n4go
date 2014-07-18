package merge_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("merge-strings -d dirName", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "merge_strings")
		inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")
	})

	Context("can combine multiple language files", func() {
		Context("merging en files in input_files path", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "merge-strings", "-v", "-d", filepath.Join(inputFilesPath), "--source-language", "en")
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "en.all.json"),
				)
			})

			It("creates an en.all.json that contains translations from both files", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "en.all.json"),
					GetFilePath(inputFilesPath, "en.all.json"),
				)
			})

			It("creates an en.all.json for which the translation strings order are stable", func() {
				expectedFilePath := GetFilePath(expectedFilesPath, "en.all.json")
				actualFilePath := GetFilePath(inputFilesPath, "en.all.json")

				expectedBytes, err := ioutil.ReadFile(expectedFilePath)
				Ω(err).Should(BeNil())
				Ω(expectedBytes).ShouldNot(BeNil())

				actualBytes, err := ioutil.ReadFile(actualFilePath)
				Ω(err).Should(BeNil())
				Ω(actualBytes).ShouldNot(BeNil())

				Ω(string(expectedBytes)).Should(Equal(string(actualBytes)))
			})
		})

		Context("merging en files in input_files/reordered path", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "merge-strings", "-v", "-d", filepath.Join(inputFilesPath, "reordered"), "--source-language", "en")
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(filepath.Join(inputFilesPath, "reordered"), "en.all.json"),
				)
			})

			It("creates an en.all.json keeping the stable order", func() {
				expectedFilePath := GetFilePath(expectedFilesPath, "en.all.json")
				actualFilePath := GetFilePath(filepath.Join(inputFilesPath, "reordered"), "en.all.json")

				expectedBytes, err := ioutil.ReadFile(expectedFilePath)
				Ω(err).Should(BeNil())
				Ω(expectedBytes).ShouldNot(BeNil())

				actualBytes, err := ioutil.ReadFile(actualFilePath)
				Ω(err).Should(BeNil())
				Ω(actualBytes).ShouldNot(BeNil())

				Ω(string(expectedBytes)).Should(Equal(string(actualBytes)))
			})
		})
	})
})
