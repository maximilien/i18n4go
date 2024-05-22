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
	"os"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rewrite-package -d dirname -r", func() {
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
		Context("rewrite all templated and interpolated strings", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

				session := Runi18n("-c",
					"rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"--root-path", rootPath,
					"-r",
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("adds T() callExprs wrapping string literals", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
				bytes, err := os.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "test.go")
				bytes, err = os.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)
				Ω(actualOutput).Should(Equal(expectedOutput))
			})

			It("recurses to files in nested dirs", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "nested_dir", "test.go")
				bytes, err := os.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "nested_dir", "test.go")
				bytes, err = os.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)
				Ω(actualOutput).Should(Equal(expectedOutput))
			})

			It("adds a i18n_init.go per package", func() {
				initFile := filepath.Join(outputDir, "i18n_init.go")
				expectedBytes, err := os.ReadFile(initFile)
				Ω(err).ShouldNot(HaveOccurred())
				expected := strings.TrimSpace(string(expectedBytes))

				expectedInitFile := filepath.Join(expectedFilesPath, "i18n_init.go")
				actualBytes, err := os.ReadFile(expectedInitFile)
				Ω(err).ShouldNot(HaveOccurred())
				actual := strings.TrimSpace(string(actualBytes))

				Ω(actual).Should(Equal(expected))

				initFile = filepath.Join(outputDir, "nested_dir", "i18n_init.go")
				expectedBytes, err = os.ReadFile(initFile)
				Ω(err).ShouldNot(HaveOccurred())
				expected = strings.TrimSpace(string(expectedBytes))

				expectedInitFile = filepath.Join(expectedFilesPath, "nested_dir", "i18n_init.go")
				actualBytes, err = os.ReadFile(expectedInitFile)
				Ω(err).ShouldNot(HaveOccurred())
				actual = strings.TrimSpace(string(actualBytes))

				Ω(actual).Should(Equal(expected))
			})

			It("adds a generated i18n_resources.go file", func() {
				i18nRC := filepath.Join(outputDir, "i18n_resources.go")
				_, err := os.ReadFile(i18nRC)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("does not translate test files", func() {
				_, doesFileExistErr := os.Stat(filepath.Join(outputDir, "a_really_bad_test.go"))
				Ω(os.IsNotExist(doesFileExistErr)).Should(BeTrue())
			})
		})

		Context("rewrite only templated and interpolated strings from --i18n-strings-filename", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")

				err = os.MkdirAll(filepath.Join(outputDir, "doption"), 0777)
				Ω(err).ShouldNot(HaveOccurred())

				CopyFile(filepath.Join(expectedFilesPath, "doption", "_test.go.en.json"), filepath.Join(outputDir, "doption", "test.go.en.json"))

				session := Runi18n("-c",
					"rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("adds T() callExprs wrapping string literals", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
				generatedOutputFile := filepath.Join(outputDir, "test.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)
			})
		})

		Context("rewrite only templated and interpolated strings from --i18n-strings-filename with multiple files", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")

				err = os.MkdirAll(filepath.Join(outputDir, "doption"), 0777)
				Ω(err).ShouldNot(HaveOccurred())

				CopyFile(filepath.Join(expectedFilesPath, "doption", "_en.all.json"), filepath.Join(outputDir, "doption", "en.all.json"))

				session := Runi18n("-c",
					"rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"--i18n-strings-filename", filepath.Join(expectedFilesPath, "doption", "en.all.json"),
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("not adds T() callExprs wrapping string literals to test3.go since none are in JSON", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test3.go")
				generatedOutputFile := filepath.Join(outputDir, "test3.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)
			})
		})

		Context("rewrite only templated and interpolated strings from --i18n-strings-dirname", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")

				err = os.MkdirAll(filepath.Join(outputDir, "doption"), 0777)
				Ω(err).ShouldNot(HaveOccurred())

				CopyFile(filepath.Join(expectedFilesPath, "doption", "_test.go.en.json"), filepath.Join(outputDir, "doption", "test.go.en.json"))
				CopyFile(filepath.Join(expectedFilesPath, "doption", "_test2.go.en.json"), filepath.Join(outputDir, "doption", "test2.go.en.json"))

				session := Runi18n("-c",
					"rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("adds T() callExprs wrapping string literals for all sources in the directory", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
				generatedOutputFile := filepath.Join(outputDir, "test.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)

				expectedOutputFile = filepath.Join(expectedFilesPath, "test2.go")
				generatedOutputFile = filepath.Join(outputDir, "test2.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)
			})
		})

	})

	Context("Using cobra commands", func() {
		Context("rewrite all templated and interpolated strings", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")

				session := Runi18n("rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"--root-path", rootPath,
					"-r",
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("adds T() callExprs wrapping string literals", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
				bytes, err := os.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "test.go")
				bytes, err = os.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)
				Ω(actualOutput).Should(Equal(expectedOutput))
			})

			It("recurses to files in nested dirs", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "nested_dir", "test.go")
				bytes, err := os.ReadFile(expectedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				expectedOutput := string(bytes)

				generatedOutputFile := filepath.Join(outputDir, "nested_dir", "test.go")
				bytes, err = os.ReadFile(generatedOutputFile)
				Ω(err).ShouldNot(HaveOccurred())

				actualOutput := string(bytes)
				Ω(actualOutput).Should(Equal(expectedOutput))
			})

			It("adds a i18n_init.go per package", func() {
				initFile := filepath.Join(outputDir, "i18n_init.go")
				expectedBytes, err := os.ReadFile(initFile)
				Ω(err).ShouldNot(HaveOccurred())
				expected := strings.TrimSpace(string(expectedBytes))

				expectedInitFile := filepath.Join(expectedFilesPath, "i18n_init.go")
				actualBytes, err := os.ReadFile(expectedInitFile)
				Ω(err).ShouldNot(HaveOccurred())
				actual := strings.TrimSpace(string(actualBytes))

				Ω(actual).Should(Equal(expected))

				initFile = filepath.Join(outputDir, "nested_dir", "i18n_init.go")
				expectedBytes, err = os.ReadFile(initFile)
				Ω(err).ShouldNot(HaveOccurred())
				expected = strings.TrimSpace(string(expectedBytes))

				expectedInitFile = filepath.Join(expectedFilesPath, "nested_dir", "i18n_init.go")
				actualBytes, err = os.ReadFile(expectedInitFile)
				Ω(err).ShouldNot(HaveOccurred())
				actual = strings.TrimSpace(string(actualBytes))

				Ω(actual).Should(Equal(expected))
			})

			It("does not translate test files", func() {
				_, doesFileExistErr := os.Stat(filepath.Join(outputDir, "a_really_bad_test.go"))
				Ω(os.IsNotExist(doesFileExistErr)).Should(BeTrue())
			})
		})

		Context("rewrite only templated and interpolated strings from --i18n-strings-filename", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")

				err = os.MkdirAll(filepath.Join(outputDir, "doption"), 0777)
				Ω(err).ShouldNot(HaveOccurred())

				CopyFile(filepath.Join(expectedFilesPath, "doption", "_test.go.en.json"), filepath.Join(outputDir, "doption", "test.go.en.json"))

				session := Runi18n("rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("adds T() callExprs wrapping string literals", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
				generatedOutputFile := filepath.Join(outputDir, "test.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)
			})
		})

		Context("rewrite only templated and interpolated strings from --i18n-strings-filename with multiple files", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")

				err = os.MkdirAll(filepath.Join(outputDir, "doption"), 0777)
				Ω(err).ShouldNot(HaveOccurred())

				CopyFile(filepath.Join(expectedFilesPath, "doption", "_en.all.json"), filepath.Join(outputDir, "doption", "en.all.json"))

				session := Runi18n("rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"--i18n-strings-filename", filepath.Join(expectedFilesPath, "doption", "en.all.json"),
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("not adds T() callExprs wrapping string literals to test3.go since none are in JSON", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test3.go")
				generatedOutputFile := filepath.Join(outputDir, "test3.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)
			})
		})

		Context("rewrite only templated and interpolated strings from --i18n-strings-dirname", func() {
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).ShouldNot(HaveOccurred())

				rootPath = filepath.Join(dir, "..", "..")

				outputDir, err = os.MkdirTemp(rootPath, "i18n4go_integration")
				Ω(err).ShouldNot(HaveOccurred())

				fixturesPath = filepath.Join("..", "..", "test_fixtures", "rewrite_package")
				inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
				expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")

				err = os.MkdirAll(filepath.Join(outputDir, "doption"), 0777)
				Ω(err).ShouldNot(HaveOccurred())

				CopyFile(filepath.Join(expectedFilesPath, "doption", "_test.go.en.json"), filepath.Join(outputDir, "doption", "test.go.en.json"))
				CopyFile(filepath.Join(expectedFilesPath, "doption", "_test2.go.en.json"), filepath.Join(outputDir, "doption", "test2.go.en.json"))

				session := Runi18n("rewrite-package",
					"-d", inputFilesPath,
					"-o", outputDir,
					"--ignore-regexp", "^[.]\\w+.go$", //Ignoring .*.go files, otherwise it defaults to ignoring *test*.go
					"-v",
				)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("adds T() callExprs wrapping string literals for all sources in the directory", func() {
				expectedOutputFile := filepath.Join(expectedFilesPath, "test.go")
				generatedOutputFile := filepath.Join(outputDir, "test.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)

				expectedOutputFile = filepath.Join(expectedFilesPath, "test2.go")
				generatedOutputFile = filepath.Join(outputDir, "test2.go")
				CompareExpectedOutputToGeneratedOutput(expectedOutputFile, generatedOutputFile)
			})
		})

	})
})
