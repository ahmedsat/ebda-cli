package main

import (
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/commands"
	"github.com/ahmedsat/ebda-cli/config"
)

// Ramdan Kamel
// EG/1262

func main() {
	err := config.Configure()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	u := commands.Update{}
	if err := u.Configure(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = u.Maps()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
