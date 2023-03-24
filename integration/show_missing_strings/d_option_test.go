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

package show_missing_strings_test

import (
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("show-missing-strings -d dirName", func() {
	var (
		fixturesPath   string
		inputFilesPath string
		session        *Session
	)

	BeforeEach(func() {
		_, err := os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())

		fixturesPath = filepath.Join("..", "..", "test_fixtures", "show_missing_strings")
		inputFilesPath = filepath.Join(fixturesPath, "d_option", "input_files")
	})

	Context("When all the translated strings are in the json resource", func() {
		BeforeEach(func() {
			languageFilePath := filepath.Join(inputFilesPath, "no_missing_strings", "app.go.en.json")
			codeDirPath := filepath.Join(inputFilesPath, "no_missing_strings", "code")
			session = Runi18n("-c", "show-missing-strings", "-d", codeDirPath, "--i18n-strings-filename", languageFilePath)

			Eventually(session.ExitCode()).Should(Equal(0))
		})

		It("Should output nothing", func() {
			Ω(session).Should(Say(""))
		})
	})

	Context("When there are strings missing from the json resource", func() {
		BeforeEach(func() {
			languageFilePath := filepath.Join(inputFilesPath, "missing_strings", "app.go.en.json")
			codeDirPath := filepath.Join(inputFilesPath, "missing_strings", "code")
			session = Runi18n("-c", "show-missing-strings", "-d", codeDirPath, "--i18n-strings-filename", languageFilePath)

			Eventually(session.ExitCode()).Should(Equal(1))
		})

		It("Should output something", func() {
			Ω(session).Should(Say("Missing"))
		})
	})

	Context("When there are extra strings in the resource file", func() {
		BeforeEach(func() {
			languageFilePath := filepath.Join(inputFilesPath, "extra_strings", "app.go.en.json")
			codeDirPath := filepath.Join(inputFilesPath, "extra_strings", "code")
			session = Runi18n("-c", "show-missing-strings", "-d", codeDirPath, "--i18n-strings-filename", languageFilePath)

			Eventually(session.ExitCode()).Should(Equal(1))
		})

		It("Should output something", func() {
			Ω(session).Should(Say("Additional"))
		})
	})
})
