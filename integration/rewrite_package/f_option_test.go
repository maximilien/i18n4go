package rewrite_package_test

import (
	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("rewrite-package -f filename", func() {
	var (
		INPUT_FILES_PATH = filepath.Join("f_option", "input_files")
	// 	EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
	)

	It("doesn't fail yet", func() {
		session := Runi18n("-rewrite-package", "-f", filepath.Join(INPUT_FILES_PATH, "something"))
		Ω(session.ExitCode()).Should(Equal(0))
	})
})
