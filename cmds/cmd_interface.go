package cmds

import (
	"github.com/maximilien/i18n4cf/common"
)

type CommandInterface interface {
	common.PrinterInterface
	Options() common.Options
	Run() error
}
