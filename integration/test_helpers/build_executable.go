package test_helpers

import (
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
)

var gi18nExec string

func BuildExecutable() {
	var err error
	gi18nExec, err = gexec.Build("./../../gi18n")
	Ω(err).ShouldNot(HaveOccurred())
}
