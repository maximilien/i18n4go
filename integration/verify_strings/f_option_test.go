package verify_strings_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("verify-strings -f fileName -languages \"[lang,?]+\"", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("f_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
	)

	Context("valid input file provided", func() {
		BeforeEach(func() {
			session := Runi18n("-verify-strings", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "app.go.en.json", "-languages", "fr,de"))
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			RemoveAllFiles(
				GetFilePath(INPUT_FILES_PATH, "app.go.en.json"),
			)
		})

		It("passes verification with language file with valid keys", func() {
			Ω(true).Should(Equal(false))
		})

		It("passes verification with multiple language files with valid keys", func() {
		})

		It("fails verification with language files with invalid keys", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})

		It("fails verification with language files with missing keys", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})

		It("warns verification with language files with additional keys", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})

		It("warns verification when no language file is found", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})
	})

	Context("invalid input file provided", func() {
		It("does not exist", func() {
		})

		It("cannot be processed", func() {
		})

		It("does not have any keys", func() {
		})
	})
})
