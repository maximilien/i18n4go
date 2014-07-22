package fixup_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("fixup", func() {
	var (
		fixturesPath string
		session      *Session
		writer       io.WriteCloser
		curDir       string
		jsonFiles    map[string][]byte
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
		// restore json files
		for path, bytes := range jsonFiles {
			err = ioutil.WriteFile(path, bytes, 0666)
			if err != nil {
				fmt.Println("Could not rewrite backup JSON files")
				panic(err.Error())
			}
		}

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

		// backup json files
		jsonFiles, err = storeTranslationFiles(".")
		if err != nil {
			fmt.Println("Could not back up the JSON files.")
			panic(err.Error())
		}

		//session = Runi18n("-c", "fixup")
		cmd := exec.Command(Gi18nExec, "-c", "fixup")
		writer, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println("Could not get the stdin pipe.")
			panic(err.Error())
		}
		session, err = Start(cmd, GinkgoWriter, GinkgoWriter)
		if err != nil {
			fmt.Println("Could not run fixup")
			panic(err.Error())
		}
	})

	Context("When there are no problems", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "allgood")
		})

		It("returns 0 and prints a reassuring message", func() {
			session.Wait()
			Ω(session.ExitCode()).Should(Equal(0))
			Ω(session).Should(Say("OK"))
		})
	})

	Context("When there are brand new strings in the code that don't exist in en_US", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "notsogood")
		})

		It("adds strings to all the locales", func() {
			session.Wait()
			output := string(session.Out.Contents())

			Ω(output).Should(ContainSubstring("Adding these strings to the translation file:"))
			Ω(output).Should(ContainSubstring("Heal the world"))

			file, err := ioutil.ReadFile(filepath.Join(".", "translations", "en_US.all.json"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(file).Should(ContainSubstring("\"Heal the world\""))

			chineseFile, err := ioutil.ReadFile(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(chineseFile).Should(ContainSubstring("\"Heal the world\""))

		})
	})

	Context("When there are old strings in the translations that don't exist in the code", func() {
		It("removes the strings from all the locales", func() {
			session.Wait()
			output := string(session.Out.Contents())

			Ω(output).Should(ContainSubstring("Removing these strings from the translation file:"))
			Ω(output).Should(ContainSubstring("Make it a better place"))

			file, err := ioutil.ReadFile(filepath.Join(".", "translations", "en_US.all.json"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(file).ShouldNot(ContainSubstring("\"Make it a better place\""))

			chineseFile, err := ioutil.ReadFile(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(chineseFile).ShouldNot(ContainSubstring("\"Make it a better place\""))
		})
	})

	Context("When a string has been updated in the code", func() {
		It("prompts the user again if they do not input a correct response", func() {

		})

		It("displayes all the possible translations that the new string could map to", func() {

		})

		Context("When the user says the translation was updated", func() {
			It("marks the foreign language translations as dirty", func() {

			})

			It("Updates the keys for all translation files", func() {

			})

			It("Updates the english translation", func() {

			})
		})

		Context("When the user says the translation was not updated", func() {
			It("adds the new translation", func() {

			})
		})

		Context("when a user quits the interactive prompt", func() {
			It("does not add any new translations", func() {

			})

			It("does not remove any translations", func() {

			})

			It("does not update any of the translations", func() {

			})
		})
	})
})

func storeTranslationFiles(dir string) (files map[string][]byte, err error) {
	files = make(map[string][]byte)
	contents, _ := ioutil.ReadDir(dir)

	for _, fileInfo := range contents {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()

			if strings.HasSuffix(name, ".all.json") {
				path := filepath.Join(dir, fileInfo.Name())
				files[path], err = ioutil.ReadFile(path)

				if err != nil {
					return nil, err
				}
			}
		} else {
			newFiles, err := storeTranslationFiles(filepath.Join(dir, fileInfo.Name()))
			if err != nil {
				return nil, err
			}

			for path, bytes := range newFiles {
				files[path] = bytes
			}
		}
	}

	return
}
