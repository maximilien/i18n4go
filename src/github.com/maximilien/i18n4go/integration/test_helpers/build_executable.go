package test_helpers

import (
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
)

var Gi18nExec string

func BuildExecutable() {
	var err error
	Gi18nExec, err = gexec.Build("./../../gi18n")
	Ω(err).ShouldNot(HaveOccurred())
}
