// Copyright © 2015-2023 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rewrite_package_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rewrite-package [...] --init-code-snippet-filename some-file", func() {
	var (
		outputDir         string
		rootPath          string
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	AfterEach(func() {
		err := os.RemoveAll(outputDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("Using legacy commands", func() {
		Context("invokes rewrite-package command and uses the default i18n init function", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())
				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "expected_output")

				session := Runi18n("-c",
					"rewrite-package",
					"-f", filepath.Join(inputFilesPath, "issue14.go"),
					"-o", outputDir,
					"--root-path", rootPath,
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("rewrites the source go file wrapping strings with T() and generates a i18n_init.go using default", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "issue14.go")
				bytes, err := ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "issue14.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))

				expectedOutputFile = filepath.Join(expectedFilesPath, "i18n_init_default.go")
				bytes, err = ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput = strings.Trim(string(bytes), "\n")

				generatedOutputFile = filepath.Join(outputDir, "i18n_init.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput = string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))
			})
		})

		Context("invokes rewrite-package command and uses the specified --init-code-snippet-filename", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())
				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "expected_output")

				session := Runi18n("-c",
					"rewrite-package",
					"-f", filepath.Join(inputFilesPath, "issue14.go"),
					"-o", outputDir,
					"--init-code-snippet-filename", filepath.Join(inputFilesPath, "init_code_snippet.go.template"),
					"--root-path", rootPath,
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("rewrites the source go file wrapping strings with T() and generates a i18n_init.go using teamplate file", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "issue14.go")
				bytes, err := ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "issue14.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))

				expectedOutputFile = filepath.Join(expectedFilesPath, "i18n_init_from_template.go")
				bytes, err = ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput = string(bytes)

				generatedOutputFile = filepath.Join(outputDir, "i18n_init.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput = string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))
			})
		})

	})

	Context("Using cobra commands", func() {
		Context("invokes rewrite-package command and uses the default i18n init function", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())
				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "expected_output")

				session := Runi18n("rewrite-package",
					"-f", filepath.Join(inputFilesPath, "issue14.go"),
					"-o", outputDir,
					"--root-path", rootPath,
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("rewrites the source go file wrapping strings with T() and generates a i18n_init.go using default", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "issue14.go")
				bytes, err := ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "issue14.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))

				expectedOutputFile = filepath.Join(expectedFilesPath, "i18n_init_default.go")
				bytes, err = ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput = strings.Trim(string(bytes), "\n")

				generatedOutputFile = filepath.Join(outputDir, "i18n_init.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput = string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))
			})
		})

		Context("invokes rewrite-package command and uses the specified --init-code-snippet-filename", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())
				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = ioutil.TempDir(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "init_code_snippet_filename", "expected_output")

				session := Runi18n("rewrite-package",
					"-f", filepath.Join(inputFilesPath, "issue14.go"),
					"-o", outputDir,
					"--init-code-snippet-filename", filepath.Join(inputFilesPath, "init_code_snippet.go.template"),
					"--root-path", rootPath,
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("rewrites the source go file wrapping strings with T() and generates a i18n_init.go using template file", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "issue14.go")
				bytes, err := ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "issue14.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))

				expectedOutputFile = filepath.Join(expectedFilesPath, "i18n_init_from_template.go")
				bytes, err = ioutil.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput = string(bytes)

				generatedOutputFile = filepath.Join(outputDir, "i18n_init.go")
				bytes, err = ioutil.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput = string(bytes)

				Ω(actualOutput).Should(Equal(expectedOutput))
			})
		})

	})
})
