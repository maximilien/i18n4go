package common

import (
	"io/ioutil"
	"strings"
)

func ParseLanguages(languages string) []string {
	langArray := strings.Split(languages, ",")
	parsedLanguages := make([]string, len(langArray))
	for i, language := range langArray {
		parsedLanguages[i] = strings.Trim(language, " ")
	}
	return parsedLanguages
}

func CopyFileContents(src, dst string) error {
	byteArray, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, byteArray, 0644)
}
