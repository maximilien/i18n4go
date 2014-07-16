package organization

import (
	"errors"
	"fmt"
)

type DeleteOrg struct {
}

type CommandMetadata struct {
	Name        string
	ShortName   string
	Description string
	Usage       string
}

type BoolFlag struct {
	Name  string
	Usage string
}

func (command *DeleteOrg) Metadata() CommandMetadata {
	return CommandMetadata{
		Name:        "delete-org",
		Description: "Delete an org",
		Usage:       "CF_NAME delete-org ORG [-f]",
		Flags:       fmt.Printf(BoolFlag{Name: "f", Usage: "Force deletion without confirmation"}),
	}
}

func (cmd *DeleteOrg) GetRequirements(err error) {
	err = errors.New("Incorrect Usage")
}

func (cmd *DeleteOrg) Run() {
	var orgName, username string

	fmt.Println("f")
	fmt.Println("org", orgName)

	fmt.Printf("Deleting org %s as %s...", orgName, username)
	fmt.Printf("Org %s does not exist.", orgName)
	return
}
