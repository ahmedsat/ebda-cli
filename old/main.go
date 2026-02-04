package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ahmedsat/erp-reports-cli/commands"
	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

var subcommands = []string{"totals", "farms", "farm-applications", "salary", "follow-up"}

func usage() {
	fmt.Printf("Usage: %s subcommand [options]\n", os.Args[0])
	fmt.Println("subcommands:")
	for _, subcommand := range subcommands {
		fmt.Printf("  %s\n", subcommand)
	}
}

func main() {

	r, err := erp.Login()
	if err != nil {
		utils.HandelErr(err)
	}

	fmt.Fprintln(os.Stderr, "Logged in as: ", r)

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "totals":
		cmd := flag.NewFlagSet("totals", flag.ExitOnError)
		opt := commands.TotalsOptions{}
		opt.AddFlags(cmd)

		utils.HandelErr(cmd.Parse(os.Args[2:]))

		utils.HandelErr(commands.Totals(opt))

	case "farms":
		cmd := flag.NewFlagSet("farms", flag.ExitOnError)
		opt := commands.FarmsOptions{}
		opt.AddFlags(cmd)

		utils.HandelErr(cmd.Parse(os.Args[2:]))

		utils.HandelErr(commands.Farms(opt))

	case "farm-applications":
		utils.HandelErr(commands.FarmApplications(os.Args[2:]))
	case "training":
		cmd := flag.NewFlagSet("training", flag.ExitOnError)
		opt := commands.TrainingOptions{}
		opt.AddFlags(cmd)

		utils.HandelErr(cmd.Parse(os.Args[2:]))

		utils.HandelErr(commands.Training(opt))

	case "salary":
		cmd := flag.NewFlagSet("salary", flag.ExitOnError)
		opt := commands.SalaryOptions{}
		opt.AddFlags(cmd)

		utils.HandelErr(cmd.Parse(os.Args[2:]))

		utils.HandelErr(commands.Salary(opt))

	case "map":
		utils.HandelErr(commands.Map(os.Args[2:]))

	case "pgs":
		utils.HandelErr(commands.Pgs(os.Args[2:]))

	case "records":
		utils.HandelErr(commands.Records(os.Args[2:]))
	case "follow-up":
		utils.HandelErr(commands.FollowUp(os.Args[2:]))

	case "help":
		usage()

	default:
		usage()
	}
}
