package fixup_test

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/maximilien/i18n4go/common"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("fixup", func() {
	var (
		fixturesPath string
		cmd          *exec.Cmd
		stdinPipe    io.WriteCloser
		stdoutPipe   io.ReadCloser
		stdoutReader *bufio.Reader
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
		cmd = exec.Command(I18n4goExec, "-c", "fixup")

		stdinPipe, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println("Could not get the stdin pipe.")
			panic(err.Error())
		}

		stdoutPipe, err = cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Could not get the stdout pipe.")
			panic(err.Error())
		}
		stdoutReader = bufio.NewReader(stdoutPipe)

		_, err = cmd.StderrPipe()
		if err != nil {
			fmt.Println("Could not get the stderr pipe.")
			panic(err.Error())
		}

		err = cmd.Start()
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
			Ω(getNextOutputLine(stdoutReader)).Should(Equal("OK"))

			exitCode := cmd.Wait()
			Ω(exitCode).Should(BeNil())
		})

		Context("for adding", func() {
			It("does not print empty add messages", func() {
				Ω(getNextOutputLine(stdoutReader)).ShouldNot(ContainSubstring("Adding these strings"))

				exitCode := cmd.Wait()
				Ω(exitCode).Should(BeNil())
			})
		})

		Context("for deleting", func() {
			It("does not print empty delete messages", func() {
				Ω(getNextOutputLine(stdoutReader)).ShouldNot(ContainSubstring("Removing these strings"))

				exitCode := cmd.Wait()
				Ω(exitCode).Should(BeNil())
			})
		})

		Context("For updating", func() {
			It("does not print empty new or update check message", func() {
				Ω(getNextOutputLine(stdoutReader)).ShouldNot(ContainSubstring("new or updated string"))

				exitCode := cmd.Wait()
				Ω(exitCode).Should(BeNil())
			})
		})
	})

	Context("When there are brand new strings in the code that don't exist in en_US", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "notsogood", "add")
		})

		It("adds strings to all the locales", func() {
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Adding these strings"))
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Heal the world"))

			exitCode := cmd.Wait()
			Ω(exitCode).Should(BeNil())

			file, err := ioutil.ReadFile(filepath.Join(".", "translations", "en_US.all.json"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(file).Should(ContainSubstring("\"Heal the world\""))

			chineseFile, err := ioutil.ReadFile(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(chineseFile).Should(ContainSubstring("\"Heal the world\""))
		})
	})

	Context("When there are old strings in the translations that don't exist in the code", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "notsogood", "delete")
		})

		It("removes the strings from all the locales", func() {
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Removing these strings"))
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Heal the world"))

			exitCode := cmd.Wait()
			Ω(exitCode).Should(BeNil())

			file, err := ioutil.ReadFile(filepath.Join(".", "translations", "en_US.all.json"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(file).ShouldNot(ContainSubstring("\"Heal the world\""))

			chineseFile, err := ioutil.ReadFile(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(chineseFile).ShouldNot(ContainSubstring("\"Heal the world\""))
		})
	})

	Context("When a string has been updated or added in the code", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "notsogood", "update")
		})

		It("cancels the interactive update when a user types exit", func() {
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Is the string \"I like apples.\" a new or updated string? [new/upd]"))

			stdinPipe.Write([]byte("exit\n"))
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Canceling fixup"))

			exitCode := cmd.Wait()
			Ω(exitCode).Should(BeNil())
		})

		It("prompts the user again if they do not input a correct response", func() {
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Is the string \"I like apples.\" a new or updated string? [new/upd]"))

			stdinPipe.Write([]byte("nope\n"))

			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Invalid response"))
			Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Is the string \"I like apples.\" a new or updated string? [new/upd]"))

			stdinPipe.Write([]byte("exit\n"))

			exitCode := cmd.Wait()
			Ω(exitCode).Should(BeNil())
		})

		Context("When the user says the translation was updated", func() {
			JustBeforeEach(func() {
				Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Is the string \"I like apples.\" a new or updated string? [new/upd]"))
				stdinPipe.Write([]byte("upd\n"))
				stdinPipe.Write([]byte("1\n"))
			})

			It("Updates the keys for all translation files", func() {
				cmd.Wait()

				translations, err := common.LoadI18nStringInfos(filepath.Join(".", "translations", "en_US.all.json"))
				Ω(err).ShouldNot(HaveOccurred())
				mappedTranslations, err := common.CreateI18nStringInfoMap(translations)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(mappedTranslations["I like bananas."]).Should(Equal(common.I18nStringInfo{}))
				Ω(mappedTranslations["I like apples."]).ShouldNot(Equal(common.I18nStringInfo{}))

				translations, err = common.LoadI18nStringInfos(filepath.Join(".", "translations", "zh_CN.all.json"))
				Ω(err).ShouldNot(HaveOccurred())
				mappedTranslations, err = common.CreateI18nStringInfoMap(translations)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(mappedTranslations["I like bananas."]).Should(Equal(common.I18nStringInfo{}))
				Ω(mappedTranslations["I like apples."]).ShouldNot(Equal(common.I18nStringInfo{}))
			})

			It("Updates all the translation", func() {
				cmd.Wait()

				translations, err := common.LoadI18nStringInfos(filepath.Join(".", "translations", "en_US.all.json"))
				Ω(err).ShouldNot(HaveOccurred())
				mappedTranslations, err := common.CreateI18nStringInfoMap(translations)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(mappedTranslations["I like apples."].Translation).Should(Equal("I like apples."))
			})

			It("marks the foreign language translations as updated", func() {
				cmd.Wait()

				translations, err := common.LoadI18nStringInfos(filepath.Join(".", "translations", "zh_CN.all.json"))
				Ω(err).ShouldNot(HaveOccurred())
				mappedTranslations, err := common.CreateI18nStringInfoMap(translations)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(mappedTranslations["I like apples."].Modified).Should(BeTrue())
				Ω(mappedTranslations["I like apples."].Translation).ShouldNot(Equal("I like apples."))
			})
		})

		Context("When the user says the translation is new", func() {
			var (
				apple = common.I18nStringInfo{ID: "I like apples.", Translation: "I like apples.", Modified: false}
			)

			JustBeforeEach(func() {
				Ω(getNextOutputLine(stdoutReader)).Should(ContainSubstring("Is the string \"I like apples.\" a new or updated string? [new/upd]"))
				stdinPipe.Write([]byte("new\n"))
			})

			It("adds the new translation and deletes the old translation from all translation files", func() {
				cmd.Wait()

				translations, err := common.LoadI18nStringInfos(filepath.Join(".", "translations", "en_US.all.json"))
				Ω(err).ShouldNot(HaveOccurred())
				mappedTranslations, err := common.CreateI18nStringInfoMap(translations)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(mappedTranslations["I like bananas."]).Should(Equal(common.I18nStringInfo{}))
				Ω(mappedTranslations["I like apples."]).Should(Equal(apple))

				translations, err = common.LoadI18nStringInfos(filepath.Join(".", "translations", "zh_CN.all.json"))
				Ω(err).ShouldNot(HaveOccurred())
				mappedTranslations, err = common.CreateI18nStringInfoMap(translations)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(mappedTranslations["I like bananas."]).Should(Equal(common.I18nStringInfo{}))
				Ω(mappedTranslations["I like apples."]).Should(Equal(apple))
			})
		})
	})

	Context("When a foreign language is missing an english translation", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "notsogood", "missing_foreign_key")
		})

		It("adds the extra translation", func() {
			cmd.Wait()

			translations, err := common.LoadI18nStringInfos(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())
			mappedTranslations, err := common.CreateI18nStringInfoMap(translations)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(mappedTranslations["I'm the extra key"]).Should(Equal(
				common.I18nStringInfo{ID: "I'm the extra key", Translation: "I'm the extra key", Modified: false},
			))
		})
	})

	Context("When a foreign language has an extra key", func() {
		BeforeEach(func() {
			fixturesPath = filepath.Join("..", "..", "test_fixtures", "fixup", "notsogood", "extra_foreign_key")
		})

		It("removes the extra translation", func() {
			cmd.Wait()

			translations, err := common.LoadI18nStringInfos(filepath.Join(".", "translations", "zh_CN.all.json"))
			Ω(err).ShouldNot(HaveOccurred())
			mappedTranslations, err := common.CreateI18nStringInfoMap(translations)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(mappedTranslations["I'm the extra key"]).Should(Equal(common.I18nStringInfo{}))
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

func getNextOutputLine(reader *bufio.Reader) string {
	line, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}

	return string(line)
}
