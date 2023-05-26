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

var _ = Describe("extract-strings -f fileName -o outputDir", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
		outputPath        string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "extract_strings", "f_option")
		inputFilesPath = filepath.Join(fixturesPath, "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "expected_output")
	})

	BeforeEach(func() {
		var err error
		outputPath, err = ioutil.TempDir("", "i18n4go4go")
		Ω(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(outputPath)
	})
	Context("Using legacy commands", func() {
		Context("-o outputDir --output-flat (default)", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "--po", "--meta", "-f", filepath.Join(inputFilesPath, "app.go"), "-o", outputPath)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("Walks input directory and compares each group of generated output to expected output", func() {

				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(expectedFilesPath, "app.go.en.json"),
					filepath.Join(outputPath, "app.go.en.json"),
				)

				CompareExpectedToGeneratedExtendedJson(
					filepath.Join(expectedFilesPath, "app.go.extracted.json"),
					filepath.Join(outputPath, "app.go.extracted.json"),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(expectedFilesPath, "app.go.en.po"),
					filepath.Join(outputPath, "app.go.en.po"),
				)
			})
		})

		Context("-o outputDir --output-match-package", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "extract-strings", "-v", "--po", "--meta", "-f", filepath.Join(inputFilesPath, "app.go"), "-o", outputPath, "--output-match-package")
				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("Walks input directory and compares each group of generated output to expected output and package subdirectories", func() {
				expectedFilesPath = filepath.Join(expectedFilesPath, "app")
				outputPath = filepath.Join(outputPath, "app")
				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(expectedFilesPath, "app.go.en.json"),
					filepath.Join(outputPath, "app.go.en.json"),
				)

				CompareExpectedToGeneratedExtendedJson(
					filepath.Join(expectedFilesPath, "app.go.extracted.json"),
					filepath.Join(outputPath, "app.go.extracted.json"),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(expectedFilesPath, "app.go.en.po"),
					filepath.Join(outputPath, "app.go.en.po"),
				)
			})
		})

	})

	Context("Using cobra commands", func() {
		Context("-o outputDir --output-flat (default)", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "--po", "--meta", "-f", filepath.Join(inputFilesPath, "app.go"), "-o", outputPath)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("Walks input directory and compares each group of generated output to expected output", func() {

				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(expectedFilesPath, "app.go.en.json"),
					filepath.Join(outputPath, "app.go.en.json"),
				)

				CompareExpectedToGeneratedExtendedJson(
					filepath.Join(expectedFilesPath, "app.go.extracted.json"),
					filepath.Join(outputPath, "app.go.extracted.json"),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(expectedFilesPath, "app.go.en.po"),
					filepath.Join(outputPath, "app.go.en.po"),
				)
			})
		})

		Context("-o outputDir --output-match-package", func() {
			BeforeEach(func() {
				session := Runi18n("extract-strings", "-v", "--po", "--meta", "-f", filepath.Join(inputFilesPath, "app.go"), "-o", outputPath, "--output-match-package")
				Ω(session.ExitCode()).Should(Equal(0))
			})

			It("Walks input directory and compares each group of generated output to expected output and package subdirectories", func() {
				expectedFilesPath = filepath.Join(expectedFilesPath, "app")
				outputPath = filepath.Join(outputPath, "app")
				CompareExpectedToGeneratedTraslationJson(
					filepath.Join(expectedFilesPath, "app.go.en.json"),
					filepath.Join(outputPath, "app.go.en.json"),
				)

				CompareExpectedToGeneratedExtendedJson(
					filepath.Join(expectedFilesPath, "app.go.extracted.json"),
					filepath.Join(outputPath, "app.go.extracted.json"),
				)

				CompareExpectedToGeneratedPo(
					filepath.Join(expectedFilesPath, "app.go.en.po"),
					filepath.Join(outputPath, "app.go.en.po"),
				)
			})
		})

	})
})
