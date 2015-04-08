package test_helpers

import (
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
)

var I18n4goExec string

func BuildExecutable() {
	var err error
	I18n4goExec, err = gexec.Build("./../../i18n4go")
	Ω(err).ShouldNot(HaveOccurred())
}
