package organization

import (
	"errors"
	"fmt"
)

type CreateOrg struct {
}

type CommandMetadata struct {
	Name        string
	ShortName   string
	Description string
	Usage       string
}

func (command CreateOrg) Metadata() CommandMetadata {
	return CommandMetadata{
		Name:        "create-org",
		ShortName:   "co",
		Description: "Create an org",
		Usage:       "CF_NAME create-org ORG",
	}
}

func (cmd CreateOrg) GetRequirements(err error) {
	err = errors.New("Incorrect Usage")
}

func (cmd CreateOrg) Run() {
	var name, username string
	fmt.Printf("Creating org %s as %s...", name, username)
	fmt.Printf("Org %s already exists", name)
	fmt.Printf("\nTIP: Use '%s' to target new org", name+" target -o "+name)
}
