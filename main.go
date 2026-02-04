package main

import (
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/commands"
	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
)

var executable string

type subcommand interface {
	Name() string
	Usage() string
	Run([]string) error
	Result() any
}

var subcommands map[string]subcommand

func init() {
	subcommands = make(map[string]subcommand)
	subcommands["help"] = &HelpCommand{}
	subcommands["follow-up"] = &commands.FollowUpCommand{}
	subcommands["pgs"] = &commands.Pgs{}
}

func usage(executable string) {
	fmt.Fprintf(os.Stderr, "Usage: %s subcommand [options]\n", executable)
	fmt.Fprintln(os.Stderr, "subcommands:")
	for subcommand := range subcommands {
		fmt.Fprintf(os.Stderr, "  %s\n", subcommand)
	}
}

func main() {

	err := config.Configure()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	res, err := frappe.Login()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Logged result %+v\n", res)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not subcommand provided")
		usage(os.Args[0])
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

	result := sbc.Result()
	if result != nil {
		fmt.Printf("%+v\n", result)
	}

}
