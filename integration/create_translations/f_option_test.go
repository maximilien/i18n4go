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

package create_translations_test

import (
	"path/filepath"

	"os"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("create-translations -f fileName --languages \"[lang,?]+\"", func() {
	var (
		fixturesPath      string
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "create_translations", "f_option")
		inputFilesPath = filepath.Join(fixturesPath, "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "expected_output")
	})

	Context("when valid input file is provided", func() {
		Context("and a single language is specified", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.fr.json"),
				)
			})

			It("creates the language file", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.fr.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.fr.json"))
			})
		})

		Context("and multiple languages are specified", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr,de\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(0))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.fr.json"),
					GetFilePath(expectedFilesPath, "quota.go.de.json"),
				)
			})

			It("creates a language file for each language specified", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.fr.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.fr.json"))

				fileInfo, err = os.Stat(GetFilePath(expectedFilesPath, "quota.go.de.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.de.json"))
			})
		})

	})

	Context("when invalid input file is provided", func() {
		Context("and file does not exist", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.de.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "zh_TW")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.fr.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})

		Context("and file is empty", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "create-translations", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.ja.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "ja")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.ja.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})
	})
})
