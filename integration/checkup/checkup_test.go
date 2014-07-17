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

	JustBeforeEach(func() {
		err = os.Chdir(fixturesPath)
		if err != nil {
			fmt.Println("Could not change to fixtures directory")
			panic(err.Error())
		}

		session = Runi18n("-c", "checkup")
	})

	Context("When there are no problems", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "checkup", "allgood")
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
})
