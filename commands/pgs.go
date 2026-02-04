package commands

import (
	"fmt"
	"strings"

	"github.com/ahmedsat/ebda-cli/kobo"
)

type Pgs struct {
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
	return sb.String()
}

// Run implements [main.subcommand].
func (p *Pgs) Run([]string) (err error) {
	p.Submissions, err = kobo.GetAssets[kobo.PGSNew]()
	return
}

// Usage implements [main.subcommand].
func (p *Pgs) Usage() string {
	panic("unimplemented")
}
