package cmds

import (
	"fmt"

	"github.com/maximilien/i18n4go/i18n4go/i18n"
	"github.com/spf13/cobra"
)

var Version string
var BuildDate string
var GitRevision string

// NewVersionCommand implements 'i18n version' command
func NewVersionCommand(p *I18NParams) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("Show the version of the i18n client"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(i18n.T("Version:      {{.Arg0}}\n", map[string]interface{}{"Arg0": Version}))
			fmt.Printf(i18n.T("Build Date:   {{.Arg0}}\n", map[string]interface{}{"Arg0": BuildDate}))
			fmt.Printf(i18n.T("Git Revision: {{.Arg0}}\n", map[string]interface{}{"Arg0": GitRevision}))
			return nil
		},
	}
	return versionCmd
}
