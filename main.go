package main

import (
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/commands"
	"github.com/ahmedsat/ebda-cli/commands/gui"
	"github.com/ahmedsat/ebda-cli/commands/training"
	"github.com/ahmedsat/ebda-cli/commands/web"
	"github.com/ahmedsat/ebda-cli/config"
)

var executable string

type subcommand interface {
	Name() (name string)
	Usage() (usage string)
	Run(args []string) (err error)
	Description() (desc string)
}

var subcommands = map[string]subcommand{}

func AddSubCommand(scs ...subcommand) {
	for _, sc := range scs {
		subcommands[sc.Name()] = sc
	}
}

func init() {
	AddSubCommand(
		&HelpCommand{},
		&commands.Totals{},
		&commands.FollowUpCommand{},
		&commands.Pgs{},
		&commands.Map{},
		&commands.Soil{},
		&commands.Info{},
		&commands.Farm{},
		&training.Training{},
		&gui.Gui{},
		&web.WebUi{},
		&commands.Update{},
	)
}

func usage(executable string) {
	fmt.Fprintf(os.Stderr, "Usage: %s subcommand [options]\n", executable)
	fmt.Fprintln(os.Stderr, "subcommands:")
	for _, subcommand := range subcommands {

		fmt.Fprintf(os.Stderr, "  %-10s : %s\n", subcommand.Name(), subcommand.Description())
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not subcommand provided")
		usage(os.Args[0])
		os.Exit(1)
	}

	err := config.Configure()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	executable = os.Args[0]
	subcommand := os.Args[1]

	sbc, ok := subcommands[subcommand]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", subcommand)
		usage(executable)
		os.Exit(1)
	}

	err = sbc.Run(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

}
