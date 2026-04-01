package commands

import (
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/utils"
)

type Debug struct{}

// Result implements [main.subcommand].
func (d *Debug) Result() any {
	panic("unimplemented")
}

// Description implements [main.subcommand].
func (d *Debug) Description() string {
	return "Debug commands"
}

// Name implements [main.subcommand].
func (d *Debug) Name() string {
	return "debug"
}

// Run implements [main.subcommand].
func (d *Debug) Run(args []string) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	debugCommand := args[0]
	switch args[0] {
	case "kobo-v":
		if len(args) < 2 {
			return nil, fmt.Errorf("not enough arguments")
		}
		return KoboVersion(args[1])

	default:
		fmt.Fprintf(os.Stderr, "unavailable commands: %s\n", debugCommand)
	}

	return nil, nil
}

func KoboVersion(form_id string) (string, error) {
	// /api/v2/assets/{uid_asset}/versions/
	resp, err := kobo.Get("/api/v2/assets/" + form_id + "/versions/")
	if err != nil {
		return "", err
	}

	utils.SaveHttpResponse(*resp)

	return "", nil
}

// Usage implements [main.subcommand].
func (d *Debug) Usage() string {
	panic("unimplemented")
}
