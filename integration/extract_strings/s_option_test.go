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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/maximilien/i18n4go/common"
	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract-strings -s filePath", func() {
	var (
		outputPath            string
		fixturesPath          string
		inputFilesPath        string
		expectedFilesPath     string
		matchingGroupFilePath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		outputPath, err = ioutil.TempDir("", "i18n4go4go")
		Ω(err).ToNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "extract_strings")
		inputFilesPath = filepath.Join(fixturesPath, "s_option", "input_files", "app", "app.go")
		expectedFilesPath = filepath.Join(fixturesPath, "s_option", "expected_output")
		matchingGroupFilePath = filepath.Join(fixturesPath, "s_option", "input_files", "matching_group.json")
	})

	AfterEach(func() {
		os.RemoveAll(outputPath)
	})

	Context("Using legacy commands", func() {
		Context("When i18n4go4go is run with the -s flag", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "-f", inputFilesPath, "-o", outputPath, "-s", matchingGroupFilePath)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("use the regexFile to parse substring", func() {
				expectedFilePath := filepath.Join(expectedFilesPath, "app.go.en.json")
				actualFilePath := filepath.Join(outputPath, "app.go.en.json")

				expectedBytes, err := ioutil.ReadFile(expectedFilePath)
				Ω(err).Should(BeNil())
				Ω(expectedBytes).ShouldNot(BeNil())

				actualBytes, err := ioutil.ReadFile(actualFilePath)
				Ω(err).Should(BeNil())
				Ω(actualBytes).ShouldNot(BeNil())

				var translations []common.I18nStringInfo
				err = json.Unmarshal(actualBytes, &translations)
				Ω(err).Should(BeNil())

				Ω(translations).To(ContainElement(common.I18nStringInfo{
					ID:          "a string",
					Translation: "a string",
					Modified:    false,
				}))

				Ω(translations).To(ContainElement(common.I18nStringInfo{
					ID:          "show the app details",
					Translation: "show the app details",
					Modified:    false,
				}))
			})
		})

	})

	Context("Using cobra commands", func() {
		Context("When i18n4go4go is run with the -s flag", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "-f", inputFilesPath, "-o", outputPath, "-s", matchingGroupFilePath)

				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("use the regexFile to parse substring", func() {
				expectedFilePath := filepath.Join(expectedFilesPath, "app.go.en.json")
				actualFilePath := filepath.Join(outputPath, "app.go.en.json")

				expectedBytes, err := ioutil.ReadFile(expectedFilePath)
				Ω(err).Should(BeNil())
				Ω(expectedBytes).ShouldNot(BeNil())

				actualBytes, err := ioutil.ReadFile(actualFilePath)
				Ω(err).Should(BeNil())
				Ω(actualBytes).ShouldNot(BeNil())

				var translations []common.I18nStringInfo
				err = json.Unmarshal(actualBytes, &translations)
				Ω(err).Should(BeNil())

				Ω(translations).To(ContainElement(common.I18nStringInfo{
					ID:          "a string",
					Translation: "a string",
					Modified:    false,
				}))

				Ω(translations).To(ContainElement(common.I18nStringInfo{
					ID:          "show the app details",
					Translation: "show the app details",
					Modified:    false,
				}))
			})
		})

	})
})
