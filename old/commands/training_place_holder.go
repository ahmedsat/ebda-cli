//go:build release

package commands

import (
	"errors"
	"flag"
)

type TrainingOptions struct{}

func (t *TrainingOptions) AddFlags(cmd *flag.FlagSet) {}

func (t *TrainingOptions) Validate() (err error) {
	return errors.New("not implemented")
}

func Training(opt TrainingOptions) (err error) {
	return errors.New("not implemented")
}
