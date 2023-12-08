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

package merge_strings_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/maximilien/i18n4go/i18n4go/common"
	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("merge-strings -d dirName", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "merge_strings")
		inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "d_option", "expected_output")
	})

	Context("Using legacy commands", func() {
		Context("can combine multiple language files", func() {
			Context("merging en files in input_files path", func() {
				BeforeEach(func() {
					session := Runi18n("-c", "merge-strings", "-v", "-d", filepath.Join(inputFilesPath), "--source-language", "en")
					Ω(session.ExitCode()).Should(Equal(0))
				})

				AfterEach(func() {
					RemoveAllFiles(
						GetFilePath(inputFilesPath, "all.en.json"),
					)
				})

				It("creates an all.en.json that contains translations from both files", func() {
					CompareExpectedToGeneratedTraslationJson(
						GetFilePath(expectedFilesPath, "all.en.json"),
						GetFilePath(inputFilesPath, "all.en.json"),
					)
				})

				It("creates an all.en.json for which the translation strings order are stable", func() {
					expectedFilePath := GetFilePath(expectedFilesPath, "all.en.json")
					actualFilePath := GetFilePath(inputFilesPath, "all.en.json")

					expectedBytes, err := ioutil.ReadFile(expectedFilePath)
					Ω(err).Should(BeNil())
					Ω(expectedBytes).ShouldNot(BeNil())

					actualBytes, err := ioutil.ReadFile(actualFilePath)
					Ω(err).Should(BeNil())
					Ω(actualBytes).ShouldNot(BeNil())

					expectedTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(expectedBytes, &expectedTranslations)
					Ω(err).Should(BeNil())

					actualTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(actualBytes, &actualTranslations)
					Ω(err).Should(BeNil())

					Ω(len(actualTranslations)).Should(Equal(len(expectedTranslations)))
					for idx := range actualTranslations {
						Ω(actualTranslations[idx].ID).Should(Equal(expectedTranslations[idx].ID))
						Ω(actualTranslations[idx].Translation).Should(Equal(expectedTranslations[idx].Translation))
					}

				})
			})

			Context("merging en files in input_files/reordered path", func() {
				BeforeEach(func() {
					session := Runi18n("-c", "merge-strings", "-v", "-d", filepath.Join(inputFilesPath, "reordered"), "--source-language", "en")
					Ω(session.ExitCode()).Should(Equal(0))
				})

				AfterEach(func() {
					RemoveAllFiles(
						GetFilePath(filepath.Join(inputFilesPath, "reordered"), "all.en.json"),
					)
				})

				It("creates an all.en.json keeping the stable order", func() {
					expectedFilePath := GetFilePath(expectedFilesPath, "all.en.json")
					actualFilePath := GetFilePath(filepath.Join(inputFilesPath, "reordered"), "all.en.json")

					expectedBytes, err := ioutil.ReadFile(expectedFilePath)
					Ω(err).Should(BeNil())
					Ω(expectedBytes).ShouldNot(BeNil())

					actualBytes, err := ioutil.ReadFile(actualFilePath)
					Ω(err).Should(BeNil())
					Ω(actualBytes).ShouldNot(BeNil())

					expectedTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(expectedBytes, &expectedTranslations)
					Ω(err).Should(BeNil())

					actualTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(actualBytes, &actualTranslations)
					Ω(err).Should(BeNil())

					Ω(len(actualTranslations)).Should(Equal(len(expectedTranslations)))
					for idx := range actualTranslations {
						Ω(actualTranslations[idx].ID).Should(Equal(expectedTranslations[idx].ID))
						Ω(actualTranslations[idx].Translation).Should(Equal(expectedTranslations[idx].Translation))
					}
				})
			})
		})
	})

	Context("Using cobra commands", func() {
		Context("can combine multiple language files", func() {
			Context("merging en files in input_files path", func() {
				BeforeEach(func() {
					session := Runi18n("merge-strings", "-v", "-d", filepath.Join(inputFilesPath), "--source-language", "en")
					Ω(session.ExitCode()).Should(Equal(0))
				})

				AfterEach(func() {
					RemoveAllFiles(
						GetFilePath(inputFilesPath, "all.en.json"),
					)
				})

				It("creates an all.en.json that contains translations from both files", func() {
					CompareExpectedToGeneratedTraslationJson(
						GetFilePath(expectedFilesPath, "all.en.json"),
						GetFilePath(inputFilesPath, "all.en.json"),
					)
				})

				It("creates an all.en.json for which the translation strings order are stable", func() {
					expectedFilePath := GetFilePath(expectedFilesPath, "all.en.json")
					actualFilePath := GetFilePath(inputFilesPath, "all.en.json")

					expectedBytes, err := ioutil.ReadFile(expectedFilePath)
					Ω(err).Should(BeNil())
					Ω(expectedBytes).ShouldNot(BeNil())

					actualBytes, err := ioutil.ReadFile(actualFilePath)
					Ω(err).Should(BeNil())
					Ω(actualBytes).ShouldNot(BeNil())

					expectedTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(expectedBytes, &expectedTranslations)
					Ω(err).Should(BeNil())

					actualTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(actualBytes, &actualTranslations)
					Ω(err).Should(BeNil())

					Ω(len(actualTranslations)).Should(Equal(len(expectedTranslations)))
					for idx := range actualTranslations {
						Ω(actualTranslations[idx].ID).Should(Equal(expectedTranslations[idx].ID))
						Ω(actualTranslations[idx].Translation).Should(Equal(expectedTranslations[idx].Translation))
					}
				})
			})

			Context("merging en files in input_files/reordered path", func() {
				BeforeEach(func() {
					session := Runi18n("merge-strings", "-v", "-d", filepath.Join(inputFilesPath, "reordered"), "--source-language", "en")
					Ω(session.ExitCode()).Should(Equal(0))
				})

				AfterEach(func() {
					RemoveAllFiles(
						GetFilePath(filepath.Join(inputFilesPath, "reordered"), "all.en.json"),
					)
				})

				It("creates an all.en.json keeping the stable order", func() {
					expectedFilePath := GetFilePath(expectedFilesPath, "all.en.json")
					actualFilePath := GetFilePath(filepath.Join(inputFilesPath, "reordered"), "all.en.json")

					expectedBytes, err := ioutil.ReadFile(expectedFilePath)
					Ω(err).Should(BeNil())
					Ω(expectedBytes).ShouldNot(BeNil())

					actualBytes, err := ioutil.ReadFile(actualFilePath)
					Ω(err).Should(BeNil())
					Ω(actualBytes).ShouldNot(BeNil())

					expectedTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(expectedBytes, &expectedTranslations)
					Ω(err).Should(BeNil())

					actualTranslations := []common.I18nStringInfo{}
					err = json.Unmarshal(actualBytes, &actualTranslations)
					Ω(err).Should(BeNil())

					Ω(len(actualTranslations)).Should(Equal(len(expectedTranslations)))
					for idx := range actualTranslations {
						Ω(actualTranslations[idx].ID).Should(Equal(expectedTranslations[idx].ID))
						Ω(actualTranslations[idx].Translation).Should(Equal(expectedTranslations[idx].Translation))
					}
				})
			})
		})
	})
})
