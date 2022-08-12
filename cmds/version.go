package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string
var BuildDate string
var GitRevision string

// NewVersionCommand implements 'i18n version' command
func NewVersionCommand(p *cobra.Params) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of the i18n client",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version:      %s\n", Version)
			fmt.Printf("Build Date:   %s\n", BuildDate)
			fmt.Printf("Git Revision: %s\n", GitRevision)
			return nil
		},
	}

  return versionCmd
}
