package extract_strings_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

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

	It("should output three files: app.en.json, app.en.extracted.json, and app.en.po", func() {
		session := Runi18n("-extract-strings", "-v", "-f", filepath.Join(INPUT_FILES_PATH, "app.go"))

		Ω(session.ExitCode()).Should(Equal(0))

		expectedTranslation := ReadJson(getFilePath(EXPECTED_FILES_PATH, "app.go.en.json"))
		generatedTranslation := ReadJson(getFilePath(INPUT_FILES_PATH, "app.go.en.json"))

		Ω(reflect.DeepEqual(expectedTranslation, generatedTranslation)).Should(BeTrue())
	})
})

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
