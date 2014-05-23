package create_translations_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("create-translations -f fileName -languages \"[lang,?]+\"", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("f_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
	)

	Context("when valid input file is provided", func() {
		Context("and a single language is specified", func() {
			BeforeEach(func() {
				session := Runi18n("-create-translations", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "quota.go.en.json"), "-languages", "\"fr\"", "-o", EXPECTED_FILES_PATH)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(EXPECTED_FILES_PATH, "quota.go.fr.json"),
				)
			})

			It("creates the language file", func() {
				fileInfo, err := os.Stat(GetFilePath(EXPECTED_FILES_PATH, "quota.go.fr.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.fr.json"))
			})
		})

		Context("and multiple languages are specified", func() {
			BeforeEach(func() {
				session := Runi18n("-create-translations", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "quota.go.en.json"), "-languages", "\"fr,de\"", "-o", EXPECTED_FILES_PATH)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(EXPECTED_FILES_PATH, "quota.go.fr.json"),
					GetFilePath(EXPECTED_FILES_PATH, "quota.go.de.json"),
				)
			})

			It("creates a language file for each language specified", func() {
				fileInfo, err := os.Stat(GetFilePath(EXPECTED_FILES_PATH, "quota.go.fr.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.fr.json"))

				fileInfo, err = os.Stat(GetFilePath(EXPECTED_FILES_PATH, "quota.go.de.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.de.json"))
			})
		})

	})

	Context("when invalid input file is provided", func() {
		Context("and file does not exist", func() {
			BeforeEach(func() {
				session := Runi18n("-create-translations", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "quota.go.de.json"), "-languages", "\"fr\"", "-o", EXPECTED_FILES_PATH, "-source-language", "zh_TW")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(EXPECTED_FILES_PATH, "quota.go.fr.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})

		Context("and file is empty", func() {
			BeforeEach(func() {
				session := Runi18n("-create-translations", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "quota.go.ja.json"), "-languages", "\"fr\"", "-o", EXPECTED_FILES_PATH, "-source-language", "ja")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(EXPECTED_FILES_PATH, "quota.go.ja.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})
	})
})
