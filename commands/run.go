package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/cgo"
)

type Run struct{}

// Description implements [main.subcommand].
func (r *Run) Description() string {
	return "Run lua script"
}

// Name implements [main.subcommand].
func (r *Run) Name() string {
	return "run"
}

// Result implements [main.subcommand].
func (r *Run) Result() any {
	panic("unimplemented")
}

// Run implements [main.subcommand].
func (r *Run) Run(args []string) error {

	if len(args) < 1 {
		return errors.New("No enough arguments")
	}

	lua, err := cgo.NewLuaState()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer lua.Close()

	lua.OpenLibs()

	err = lua.DoFile(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)

	return nil
}

// Usage implements [main.subcommand].
func (r *Run) Usage() string {
	panic("unimplemented")
}
