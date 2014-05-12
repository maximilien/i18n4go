package test_helpers

import (
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
)

var gi18nExec string

func BuildExecutable() {
	var err error
	gi18nExec, err = gexec.Build("./../../main")
	Ω(err).ShouldNot(HaveOccurred())
}
