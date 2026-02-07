//go:build !release

package main

import "github.com/ahmedsat/ebda-cli/commands"

func init() {
	AddSubCommand(&commands.Missing{})

}
