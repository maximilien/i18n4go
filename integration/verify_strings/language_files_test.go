package verify_strings_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("verify-strings -f fileName -language-files \"[lang/fileName,?]+\"", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("language_files_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("language_files_option", "expected_output")
	)

	Context("valid language_files provided", func() {
		BeforeEach(func() {
			session := Runi18n("-verify-strings", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "app.go.en.json"), "-language-files", "app.go.fr.json, app.go.de.json")
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			RemoveAllFiles(
				GetFilePath(INPUT_FILES_PATH, "app.go.en.json"),
			)
		})

		It("passes verification with input file, with all valid keys", func() {
			Ω(true).Should(Equal(false))
		})

		It("fails verification with input file, when there are invalid keys", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})

		It("fails verification with input file, when there are missing keys", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})

		It("warns verification with input file, when additional keys are present", func() {
			GetFilePath(EXPECTED_FILES_PATH, "app.go.en.json.diff")
		})
	})

	Context("invalid language_files provided", func() {
		It("does not exist", func() {
		})

		It("cannot be processed", func() {
		})

		It("does not have any keys", func() {
		})
	})
})
