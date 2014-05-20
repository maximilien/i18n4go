package common

import (
	"io/ioutil"
)

func ParseLanguages(languages string) ([]string, error) {
	return []string{"en", "fr_FR", "es"}, nil
}

func CopyFileContents(src, dst string) error {
	byteArray, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, byteArray, 0644)
}
