//go:build !release

package commands

import (
	"flag"
	"fmt"
)

type TrainingOptions struct {
	GlobalOptions
}

func (t *TrainingOptions) AddFlags(cmd *flag.FlagSet) {
	t.GlobalOptions.AddFlags(cmd)
	// ? this is a placeholder
}

func (t *TrainingOptions) Validate() (err error) {
	err = t.GlobalOptions.Validate()
	if err != nil {
		return
	}

	// ? this is a placeholder

	return
}

func Training(opt TrainingOptions) (err error) {
	err = opt.Validate()
	if err != nil {
		return
	}

	fmt.Println("Training")

	return
}
