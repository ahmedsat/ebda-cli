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

// Name implements [main.subcommand].
func (p *Pgs) Name() string {
	panic("unimplemented")
}

// Result implements [main.subcommand].
func (p *Pgs) Result() any {
	sb := strings.Builder{}
	for _, s := range p.Submissions {
		fmt.Fprintf(&sb, "%s\t%s\t%s\t%s\n", s.FormID, s.VisitDate, s.EngName, s.Label)
	}

	if p.copy {
		clipboard.WriteAll(sb.String())
		return "copied to clipboard"
	}

	return sb.String()
}

// Run implements [main.subcommand].
func (p *Pgs) Run(args []string) (err error) {
	fs := flag.NewFlagSet("pgs", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)
	p.copy = *copy

	p.Submissions, err = kobo.GetAssets[kobo.PGSNew]()
	return
}

// Usage implements [main.subcommand].
func (p *Pgs) Usage() string {
	panic("unimplemented")
}
