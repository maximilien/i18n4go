package create_translations_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("create-translations -f fileName --languages \"[lang,?]+\"", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "create_translations", "f_option")
		inputFilesPath = filepath.Join(fixturesPath, "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "expected_output")
	})

	Context("when valid input file is provided", func() {
		Context("and a single language is specified", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.fr.json"),
				)
			})

			It("creates the language file", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.fr.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.fr.json"))
			})
		})

		Context("and multiple languages are specified", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr,de\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.fr.json"),
					GetFilePath(expectedFilesPath, "quota.go.de.json"),
				)
			})

			It("creates a language file for each language specified", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.fr.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.fr.json"))

				fileInfo, err = os.Stat(GetFilePath(expectedFilesPath, "quota.go.de.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.de.json"))
			})
		})

	})

	Context("when invalid input file is provided", func() {
		Context("and file does not exist", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.de.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "zh_TW")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.fr.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})

		Context("and file is empty", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.ja.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "ja")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.ja.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})
	})
})
