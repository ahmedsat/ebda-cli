package commands

import (
	_ "embed"
	"strings"
)

//go:embed farms.tsv
var farms string
var farmsMap map[string][]string

func init() {
	farmsMap = map[string][]string{}
	for line := range strings.SplitSeq(farms, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		farmsMap[fields[0]] = fields[1:]
	}
}

type Farm struct {
}

// Description implements [main.subcommand].
func (f *Farm) Description() string {
	panic("unimplemented")
}

// Name implements [main.subcommand].
func (f *Farm) Name() string {
	panic("unimplemented")
}

// Result implements [main.subcommand].
func (f *Farm) Result() any {
	panic("unimplemented")
}

// Run implements [main.subcommand].
func (f *Farm) Run([]string) error {
	panic("unimplemented")
}

// Usage implements [main.subcommand].
func (f *Farm) Usage() string {
	panic("unimplemented")
}
