package verify_strings_test

import (
	"path/filepath"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("verify-strings -f fileName", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("duplicate_keys", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("duplicate_keys", "expected_output")
	)

	Context("checks for duplicate keys", func() {
		It("should error", func() {
			session := Runi18n("-verify-strings", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "quota.go.en.json"), "-languages", "\"fr\"", "-o", EXPECTED_FILES_PATH, "-source-language", "en")
			Ω(session.ExitCode()).Should(Equal(1))
		})
	})
})
