//go:build !release

package main

import "github.com/ahmedsat/ebda-cli/commands"

func init() {
	subcommands["missing"] = &commands.Missing{}

}
