package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "v0.0.1"

type cmdOptions struct {
	version bool
	locale  string
}

func main() {
	opts := cmdOptions{}
	rootCmd := &cobra.Command{
		Use:   "foo",
		Short: T("Test cli command"),
		Long:  "foo <arg>",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.version {
				fmt.Printf(T("Version {{.Version}}", map[string]interface{}{"Version": VERSION}))
			}
		},
	}

	rootCmd.Flags().BoolVarP(&opts.version, "version", "v", false, T("Show the current version"))
	rootCmd.Flags().StringVarP(&opts.locale, "locale", "l", "en_US", T("Change current locale"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
