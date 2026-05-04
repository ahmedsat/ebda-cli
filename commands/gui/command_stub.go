//go:build release || !cgo

package gui

import "fmt"

type Gui struct{}

func (g *Gui) Description() string {
	return "Experimental desktop GUI (unavailable in release/static builds)"
}

func (g *Gui) Name() string {
	return "gui"
}

func (g *Gui) Usage() string {
	return "gui"
}

func (g *Gui) Run(args []string) error {
	return fmt.Errorf("gui is unavailable in release builds or when CGO is disabled")
}
