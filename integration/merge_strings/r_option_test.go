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
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("merge-strings -d dirName -r", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "merge_strings")
		inputFilesPath = filepath.Join(fixturesPath, "r_option", "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "r_option", "expected_output")
	})

	Context("Using legacy commands", func() {
		Context("can combine multiple language files per directory", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "merge-strings", "-v", "-r", "-d", filepath.Join(inputFilesPath), "--source-language", "en")
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "all.en.json"),
					GetFilePath(inputFilesPath+"/sub", "all.en.json"),
				)
			})

			It("all.en.json contains translations from both files", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "all.en.json"),
					GetFilePath(inputFilesPath, "all.en.json"),
				)
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath+"/sub", "all.en.json"),
					GetFilePath(inputFilesPath+"/sub", "all.en.json"),
				)
			})
		})
	})

	Context("Using cobra commands", func() {
		Context("can combine multiple language files per directory", func() {
			BeforeEach(func() {
				session := Runi18n("merge-strings", "-v", "-r", "-d", filepath.Join(inputFilesPath), "--source-language", "en")
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(inputFilesPath, "all.en.json"),
					GetFilePath(inputFilesPath+"/sub", "all.en.json"),
				)
			})

			It("all.en.json contains translations from both files", func() {
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath, "all.en.json"),
					GetFilePath(inputFilesPath, "all.en.json"),
				)
				CompareExpectedToGeneratedTraslationJson(
					GetFilePath(expectedFilesPath+"/sub", "all.en.json"),
					GetFilePath(inputFilesPath+"/sub", "all.en.json"),
				)
			})
		})
	})

})
