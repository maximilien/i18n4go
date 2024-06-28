package i18n

var T TranslateFunc

func init() {
	T = Init("", GetResourcesPath(), func(asset string) ([]byte, error) {
		return Asset(asset)
	})
}
