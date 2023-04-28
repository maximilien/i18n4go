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

package extract_strings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -f fileName", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "extract_strings", "f_option")
		inputFilesPath = filepath.Join(fixturesPath, "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "expected_output")
	})

	Context("Using legacy commands", func() {
		Context("compare generated and expected file", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "--po", "--meta", "-f", filepath.Join(inputFilesPath, "app.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "app.go.en.json"),
					GetFilePath(inputFilesPath, "app.go.en.po"),
					GetFilePath(inputFilesPath, "app.go.extracted.json"),
				)
			})

			It("app.go.en.json", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "app.go.en.json"),
					GetFilePath(inputFilesPath, "app.go.en.json"),
				)
			})

			It("app.go.extracted.json", func() {
				CompareExpectedToGeneratedExtendedJson(
					GetFilePath(expectedFilesPath, "app.go.extracted.json"),
					GetFilePath(inputFilesPath, "app.go.extracted.json"),
				)
			})

			It("app.go.en.po", func() {
				CompareExpectedToGeneratedPo(
					GetFilePath(expectedFilesPath, "app.go.en.po"),
					GetFilePath(inputFilesPath, "app.go.en.po"),
				)
			})
		})

		Context("GitHub issue #4: extracting some character as ascii code, e.g., > as \u003e", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "-f", filepath.Join(inputFilesPath, "issue4.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "issue4.go.en.json"),
				)
			})

			It("issue4.go.en.json", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "issue4.go.en.json"),
					GetFilePath(inputFilesPath, "issue4.go.en.json"),
				)
			})
		})

		Context("GitHub issue #16: Extract Strings should ignore string keys in maps", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "-f", filepath.Join(inputFilesPath, "issue16.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "issue16.go.en.json"),
				)
			})

			It("issue16.go.en.json", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "issue16.go.en.json"),
					GetFilePath(inputFilesPath, "issue16.go.en.json"),
				)
			})
		})

		Context("when the file specified has no strings at all", func() {
			var (
				OUTPUT_PATH string
			)

			BeforeEach(func() {
				var err error
				OUTPUT_PATH, err = ioutil.TempDir("", "i18n4go4go")
				Ω(err).ShouldNot(HaveOccurred())

				session := Runi18n("-c", "extract-strings", "-f", filepath.Join(inputFilesPath, "no_strings.go"), "-o", OUTPUT_PATH)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("does not generate any files", func() {
				println(OUTPUT_PATH)
				files, err := ioutil.ReadDir(OUTPUT_PATH)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(files).Should(BeEmpty())
			})
		})

		Context("GitHub issue #45: Extract Strings should extract strings string embedded inside a func, inside a func in a return", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "-f", filepath.Join(inputFilesPath, "issue45.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "issue45.go.en.json"),
				)
			})

			It("generates issue45.go.en.json correctly", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "issue45.go.en.json"),
					GetFilePath(inputFilesPath, "issue45.go.en.json"),
				)
			})
		})

	})

	Context("Using cobra commands", func() {
		Context("compare generated and expected file", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "--po", "--meta", "-f", filepath.Join(inputFilesPath, "app.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "app.go.en.json"),
					GetFilePath(inputFilesPath, "app.go.en.po"),
					GetFilePath(inputFilesPath, "app.go.extracted.json"),
				)
			})

			It("app.go.en.json", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "app.go.en.json"),
					GetFilePath(inputFilesPath, "app.go.en.json"),
				)
			})

			It("app.go.extracted.json", func() {
				CompareExpectedToGeneratedExtendedJson(
					GetFilePath(expectedFilesPath, "app.go.extracted.json"),
					GetFilePath(inputFilesPath, "app.go.extracted.json"),
				)
			})

			It("app.go.en.po", func() {
				CompareExpectedToGeneratedPo(
					GetFilePath(expectedFilesPath, "app.go.en.po"),
					GetFilePath(inputFilesPath, "app.go.en.po"),
				)
			})
		})

		Context("GitHub issue #4: extracting some character as ascii code, e.g., > as \u003e", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "-f", filepath.Join(inputFilesPath, "issue4.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "issue4.go.en.json"),
				)
			})

			It("issue4.go.en.json", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "issue4.go.en.json"),
					GetFilePath(inputFilesPath, "issue4.go.en.json"),
				)
			})
		})

		Context("GitHub issue #16: Extract Strings should ignore string keys in maps", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "-f", filepath.Join(inputFilesPath, "issue16.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "issue16.go.en.json"),
				)
			})

			It("issue16.go.en.json", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "issue16.go.en.json"),
					GetFilePath(inputFilesPath, "issue16.go.en.json"),
				)
			})
		})

		Context("when the file specified has no strings at all", func() {
			var (
				OUTPUT_PATH string
			)

			BeforeEach(func() {
				var err error
				OUTPUT_PATH, err = ioutil.TempDir("", "i18n4go4go")
				Ω(err).ShouldNot(HaveOccurred())

				session := Runi18n("extract-strings", "-f", filepath.Join(inputFilesPath, "no_strings.go"), "-o", OUTPUT_PATH)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("does not generate any files", func() {
				println(OUTPUT_PATH)
				files, err := ioutil.ReadDir(OUTPUT_PATH)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(files).Should(BeEmpty())
			})
		})

		Context("GitHub issue #45: Extract Strings should extract strings string embedded inside a func, inside a func in a return", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "-f", filepath.Join(inputFilesPath, "issue45.go"))
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "issue45.go.en.json"),
				)
			})

			It("generates issue45.go.en.json correctly", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "issue45.go.en.json"),
					GetFilePath(inputFilesPath, "issue45.go.en.json"),
				)
			})
		})

	})
})
