package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/atotto/clipboard"
)

type Pgs struct {
	copy        bool
	Submissions []kobo.PGSNew
}

// Description implements [main.subcommand].
func (p *Pgs) Description() string {
	return "Extract PGS data from kobo"
}

// Name implements [main.subcommand].
func (p *Pgs) Name() string {
	return "pgs"
}

// Result implements [main.subcommand].
func (p *Pgs) Result() any {
	sb := strings.Builder{}
	fmt.Fprintln(&sb, strings.Join([]string{
		"Code",
		"Visit Date",
		"Eng Name",
		"Label",
	}, "\t"))
	for _, s := range p.Submissions {
		fmt.Fprintln(&sb, strings.Join([]string{
			s.FarmCode,
			s.VisitDate,
			s.EngName,
			s.Label,
		}, "\t"))
	}

	if p.copy {
		clipboard.WriteAll(sb.String())
		return "copied to clipboard"
	}

	return sb.String()
}

// Run implements [main.subcommand].
func (p *Pgs) Run(args []string) (r any, err error) {
	fs := flag.NewFlagSet("pgs", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)
	p.copy = *copy

	p.Submissions, err = kobo.GetAssets[kobo.PGSNew](nil)
	r = p.Result()
	return
}

// Usage implements [main.subcommand].
func (p *Pgs) Usage() string {
	panic("unimplemented")
}
