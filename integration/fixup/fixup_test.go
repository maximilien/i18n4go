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

		It("adds strings when user types 'y'", func() {
			// Write to the writer whatever we want.
			file, err := ioutil.ReadFile("./stdin.txt")
			Expect(err).ShouldNot(HaveOccurred())

			writer.Write(file)

			session.Wait()
			Expect(err).ShouldNot(HaveOccurred())

			output := string(session.Out.Contents())

			// strings wrapped in T() in the code that don't have corresponding keys in the translation files
			Ω(output).Should(ContainSubstring("Adding these strings to the translation file:"))
			Ω(output).Should(ContainSubstring("Heal the world"))

			// check english json files and see that new string is added
			file, err = ioutil.ReadFile(filepath.Join(".", "translations", "en_US.all.json"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(file).Should(ContainSubstring("\"Heal the world\""))

			// check other translation files to see if key exists
			chineseFile, err := ioutil.ReadFile(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(chineseFile).Should(ContainSubstring("\"Heal the world\""))

			// keys in the translations that don't have corresponding strings wrapped in T() in the code
			//Ω(output).Should(ContainSubstring("\"Make it a better place\" exists in en_US, but not in the code"))

			// keys in non-english translations that don't exist in the english translation
			//Ω(output).Should(ContainSubstring("\"For you and for me\" exists in zh_CN, but not in en_US"))

			// keys that exist in the english translation but are missing in non-english translations
			//Ω(output).Should(ContainSubstring("\"And the entire human race\" exists in en_US, but not in zh_CN"))

			Ω(session.ExitCode()).Should(Equal(0))

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
