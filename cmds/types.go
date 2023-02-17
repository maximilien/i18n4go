package cmds

import (
	"io"
)

// I18NParams for creating commands.
// Useful for inserting mocks for testing
// and to have common parameters across commands.
type I18NParams struct {
	Output io.Writer
}
