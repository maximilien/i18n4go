package i18n

import "path/filepath"

var T TranslateFunc

func init() {
	T = Init("", filepath.Join("i18n4go", GetResourcesPath()), func(asset string) ([]byte, error) {
		return Asset(asset)
	})
}
