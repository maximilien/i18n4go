package extract_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4cf/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -d dirName", func() {
	var (
		INPUT_DIR_PATH    = filepath.Join("d_option", "input_files/quota")
		EXPECTED_DIR_PATH = filepath.Join("d_option", "expected_output")
		OUTPUT_PATH       string
	)

	Context("When gi18n4cf is run with the -d flag", func() {
		BeforeEach(func() {
			var err error
			OUTPUT_PATH, err = ioutil.TempDir("", "gi18n4cf")
			Ω(err).ToNot(HaveOccurred())

			session := Runi18n("-extract-strings", "-v", "-p", "-d", INPUT_DIR_PATH, "-o", OUTPUT_PATH)
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			os.RemoveAll(OUTPUT_PATH)
		})

		It("Walks input directory and compares each group of generated output to expected output", func() {
			filepath.Walk(INPUT_DIR_PATH, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(EXPECTED_DIR_PATH, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
					filepath.Join(OUTPUT_PATH, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
				)

				CompareExpectedToGeneratedExtendedJson(
					filepath.Join(EXPECTED_DIR_PATH, strings.Join([]string{filepath.Base(path), "extracted.json"}, ".")),
					filepath.Join(OUTPUT_PATH, strings.Join([]string{filepath.Base(path), "extracted.json"}, ".")),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(EXPECTED_DIR_PATH, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
					filepath.Join(OUTPUT_PATH, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
				)

				return nil
			})
		})
	})

	Context("When gi18n4cf is run with the -d -r flags", func() {
		BeforeEach(func() {
			var err error
			INPUT_DIR_PATH = filepath.Join(INPUT_DIR_PATH, "..")

			OUTPUT_PATH, err = ioutil.TempDir("", "gi18n4cf")
			Ω(err).ToNot(HaveOccurred())

			session := Runi18n("-extract-strings", "-v", "-p", "-d", INPUT_DIR_PATH, "-o", OUTPUT_PATH, "-r")
			Ω(session.ExitCode()).Should(Equal(0))
		})

		AfterEach(func() {
			os.RemoveAll(OUTPUT_PATH)
		})

		It("Walks input directories and compares each group of generated output to expected output", func() {
			filepath.Walk(INPUT_DIR_PATH, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(EXPECTED_DIR_PATH, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
					filepath.Join(OUTPUT_PATH, strings.Join([]string{filepath.Base(path), "en.json"}, ".")),
				)

				CompareExpectedToGeneratedExtendedJson(
					filepath.Join(EXPECTED_DIR_PATH, strings.Join([]string{filepath.Base(path), "extracted.json"}, ".")),
					filepath.Join(OUTPUT_PATH, strings.Join([]string{filepath.Base(path), "extracted.json"}, ".")),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(EXPECTED_DIR_PATH, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
					filepath.Join(OUTPUT_PATH, strings.Join([]string{filepath.Base(path), "en.po"}, ".")),
				)

				return nil
			})
		})
	})
})
