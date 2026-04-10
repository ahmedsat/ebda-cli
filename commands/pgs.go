package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/atotto/clipboard"
)

type Pgs struct {
}

// Description implements [main.subcommand].
func (p *Pgs) Description() string {
	return "Extract PGS data from kobo"
}

// Name implements [main.subcommand].
func (p *Pgs) Name() string {
	return "pgs"
}

// Run implements [main.subcommand].
func (p *Pgs) Run(args []string) (err error) {
	fs := flag.NewFlagSet("pgs", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)

	fmt.Fprintln(os.Stderr, "getting data")
	Submissions, err := kobo.GetAssets[kobo.PGSNew](nil)
	sb := strings.Builder{}
	fmt.Fprintln(&sb, strings.Join([]string{
		"Code",
		"Visit Date",
		"Eng Name",
		"Label",
	}, "\t"))
	for _, s := range Submissions {
		fmt.Fprintln(&sb, strings.Join([]string{
			s.AtHouseFarmId,
			s.AtHouseVisitDate,
			s.EngineerDataEngineerName,
			s.Label,
		}, "\t"))
	}

	if *copy {
		clipboard.WriteAll(sb.String())
		fmt.Println("copied to clipboard")
		return
	}

	fmt.Print(sb.String())
	return
}

// Usage implements [main.subcommand].
func (p *Pgs) Usage() string {
	panic("unimplemented")
}
