package verify_strings_test

import (
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("verify-strings -f fileName", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		Ω(err).ToNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "verify_strings")
		inputFilesPath = filepath.Join(fixturesPath, "duplicate_keys", "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "duplicate_keys", "expected_output")
	})

	Context("checks for duplicate keys", func() {
		It("should error", func() {
			session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "en")
			Ω(session.ExitCode()).Should(Equal(1))
			Ω(session).Should(gbytes.Say("Duplicated key found: Show quota info"))
		})
	})
})
