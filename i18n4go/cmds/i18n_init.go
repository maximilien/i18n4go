package cmds

import (
	i18n "github.com/maximilien/i18n4go/i18n4go/i18n"
)

var T i18n.TranslateFunc

func init() {
	T = i18n.Init("", "./i18n4go/i18n/"+i18n.GetResourcesPath())
}
