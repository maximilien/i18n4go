package checkup_test

import (
	"github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCheckup(t *testing.T) {
	BeforeSuite(test_helpers.BuildExecutable)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Checkup Suite")
}
