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

package checkup_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("checkup", func() {
	var (
		fixturesPath string
		session      *Session
		curDir       string
		err          error
	)

	BeforeEach(func() {
		curDir, err = os.Getwd()
		if err != nil {
			fmt.Println("Could not get working directory")
			panic(err.Error())
		}
	})

	AfterEach(func() {
		err = os.Chdir(curDir)
		if err != nil {
			fmt.Println("Could not change back to working directory")
			panic(err.Error())
		}
	})

	Context("When there are no problems", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "checkup", "allgood")
			err = os.Chdir(fixturesPath)
			Ω(err).ToNot(HaveOccurred(), "Could not change to fixtures directory")

			session = Runi18n("-c", "checkup", "-v")
		})

		It("returns 0", func() {
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("prints a reassuring message", func() {
			Ω(session).Should(Say("OK"))
		})
	})

	Context("when the i18n package is fully qualified", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "checkup", "qualified")
			err = os.Chdir(fixturesPath)
			Ω(err).ToNot(HaveOccurred(), "Could not change to fixtures directory")

			session = Runi18n("-c", "checkup", "-v", "-q", "i18n")
		})

		It("returns 0", func() {
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("prints a reassuring message", func() {
			session = Runi18n("-c", "checkup", "-v", "-q", "i18n")
			Ω(session).Should(Say("OK"))
		})
	})

	Context("When the translation files is in format all.<lang>.json", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "checkup", "fileformat")
			err = os.Chdir(fixturesPath)
			Ω(err).ToNot(HaveOccurred(), "Could not change to fixtures directory")

			session = Runi18n("-c", "checkup", "-v")
		})

		It("returns 0", func() {
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("prints a reassuring message", func() {
			Ω(session).Should(Say("OK"))
		})
	})

	Context("When there are problems", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "checkup", "notsogood")
			err = os.Chdir(fixturesPath)
			Ω(err).ToNot(HaveOccurred(), "Could not change to fixtures directory")

			session = Runi18n("-c", "checkup", "-v")
		})

		It("shows all inconsistent strings and returns 1", func() {
			output := string(session.Out.Contents())

			// strings wrapped in T() in the code that don't have corresponding keys in the translation files
			Ω(output).Should(ContainSubstring("\"Heal the world\" exists in the code, but not in en_US"))

			// keys in the translations that don't have corresponding strings wrapped in T() in the code
			Ω(output).Should(ContainSubstring("\"Make it a better place\" exists in en_US, but not in the code"))

			// keys in non-english translations that don't exist in the english translation
			Ω(output).Should(ContainSubstring("\"For you and for me\" exists in zh_CN, but not in en_US"))

			// keys that exist in the english translation but are missing in non-english translations
			Ω(output).Should(ContainSubstring("\"And the entire human race\" exists in en_US, but not in zh_CN"))

			Ω(session.ExitCode()).Should(Equal(1))
		})
	})

	Context("When translation IDs are (re)assigned to variables", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "checkup", "variable")
			err = os.Chdir(fixturesPath)
			Ω(err).ToNot(HaveOccurred(), "Could not change to fixtures directory")

			session = Runi18n("-c", "checkup", "-v")
		})

		It("returns 0", func() {
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("prints a reassuring message", func() {
			Ω(session).Should(Say("OK"))
		})
	})
})
