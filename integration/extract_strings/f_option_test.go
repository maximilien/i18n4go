package extract_strings_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("extract-strings -f fileName", func() {
	var (
		INPUT_FILES_PATH    = filepath.Join("f_option", "input_files")
		EXPECTED_FILES_PATH = filepath.Join("f_option", "expected_output")
		oldPwd              = os.Getenv("PWD")
	)

	BeforeEach(func() {
		oldPwd = os.Getenv("PWD")
	})

	AfterEach(func() {
		os.Setenv("PWD", oldPwd)
		removeAllFiles(
			getFilePath(INPUT_FILES_PATH, "app.go.en.json"),
			getFilePath(INPUT_FILES_PATH, "app.go.en.po"),
			getFilePath(INPUT_FILES_PATH, "app.go.extracted.json"),
		)
	})

	Context("compare generated and expected file", func() {
		BeforeEach(func() {
			session := Runi18n("-extract-strings", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "app.go"))
			Ω(session.ExitCode()).Should(Equal(0))
		})

		It("app.go.en.json", func() {
			compareExpectedToGeneratedTraslationJson(
				getFilePath(EXPECTED_FILES_PATH, "app.go.en.json"),
				getFilePath(INPUT_FILES_PATH, "app.go.en.json"),
			)
		})

		It("app.go.extracted.json", func() {
			compareExpectedToGeneratedExtendedJson(
				getFilePath(EXPECTED_FILES_PATH, "app.go.extracted.json"),
				getFilePath(INPUT_FILES_PATH, "app.go.extracted.json"),
			)
		})

		It("app.go.en.po", func() {
			compareExpectedToGeneratedPo(
				getFilePath(EXPECTED_FILES_PATH, "app.go.en.po"),
				getFilePath(INPUT_FILES_PATH, "app.go.en.po"),
			)
		})

	})
})

func compareExpectedToGeneratedPo(expectedFilePath string, generatedFilePath string) {
	expectedTranslation := ReadPo(expectedFilePath)
	generatedTranslation := ReadPo(generatedFilePath)

	Ω(reflect.DeepEqual(expectedTranslation, generatedTranslation)).Should(BeTrue())
}

func compareExpectedToGeneratedTraslationJson(expectedFilePath string, generatedFilePath string) {
	expectedTranslation := ReadJson(expectedFilePath)
	generatedTranslation := ReadJson(generatedFilePath)

	Ω(reflect.DeepEqual(expectedTranslation, generatedTranslation)).Should(BeTrue())
}

func compareExpectedToGeneratedExtendedJson(expectedFilePath string, generatedFilePath string) {
	expectedTranslation := ReadJsonExtended(expectedFilePath)
	generatedTranslation := ReadJsonExtended(generatedFilePath)

	Ω(reflect.DeepEqual(expectedTranslation, generatedTranslation)).Should(BeTrue())
}

func getFilePath(input_dir string, fileName string) string {
	return filepath.Join(os.Getenv("PWD"), input_dir, fileName)
}

func removeAllFiles(args ...string) {
	for _, arg := range args {
		os.Remove(arg)
	}
}

func Runi18n(args ...string) *Session {
	session := RunCommand(gi18nExec, args...)
	return session
}

func RunCommand(cmd string, args ...string) *Session {
	command := exec.Command(cmd, args...)
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	session.Wait()
	return session
}

func ReadPo(fileName string) map[string]string {
	file, _ := os.Open(fileName)
	r := bufio.NewReader(file)

	myMap := make(map[string]string)
	for rawLine, _, err := r.ReadLine(); err != io.EOF; rawLine, _, err = r.ReadLine() {
		if err != nil {
			Fail(fmt.Sprintf("Error: %v", err))
		}

		line := string(rawLine)
		if strings.HasPrefix(line, "msgid") {
			rawLine, _, err = r.ReadLine()
			if err != nil {
				Fail(fmt.Sprintf("Error: %v", err))
			}

			myMap[line] = string(rawLine)
		}
	}

	return myMap
}

func ReadJson(fileName string) map[string]string {
	fileByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		Fail("Cannot open json file:" + fileName)
	}

	var b interface{}

	if err := json.Unmarshal(fileByte, &b); err != nil {
		Fail(fmt.Sprintf("Cannot unmarshal: %v", err))
	}

	myMap := make(map[string]string)

	for _, value := range b.([]interface{}) {
		valueMap := value.(map[string]interface{})
		myMap[valueMap["id"].(string)] = valueMap["translation"].(string)
	}

	return myMap
}

func ReadJsonExtended(fileName string) map[string]map[string]string {
	fileByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		Fail("Cannot open json file:" + fileName)
	}

	var b interface{}

	if err := json.Unmarshal(fileByte, &b); err != nil {
		Fail(fmt.Sprintf("Cannot unmarshal: %v", err))
	}

	myMap := make(map[string]map[string]string)

	for _, value := range b.([]interface{}) {
		valueMap := value.(map[string]interface{})

		dataMap := make(map[string]string)

		for key, val := range valueMap {
			switch val.(type) {
			case string:
				dataMap[key] = val.(string)
			case float64:
				dataMap[key] = fmt.Sprintf("%v", int(val.(float64)))
			default:
				fmt.Println("We did something wrong", key)
			}
		}

		myMap[valueMap["value"].(string)] = dataMap
	}

	return myMap
}
