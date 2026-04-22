package main

import (
	"fmt"
	"os"
)

type HelpCommand struct {
}

// Description implements [subcommand].
func (h *HelpCommand) Description() string {
	return "Display help"
}

func (h *HelpCommand) Name() string { return "Help" }

func (h *HelpCommand) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Usage: %s", h.Usage())
	}

	_, ok := subcommands[args[0]]
	if !ok {
		return fmt.Errorf("Usage: %s", h.Usage())
	}

	fmt.Fprintf(os.Stderr, "Usage: %s %s\n", executable, h.Usage())
	return nil
}

func (h *HelpCommand) Usage() string {
	return fmt.Sprintf("%s subcommand [options]\n", executable)
}
