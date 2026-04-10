//go:build !release

package commands

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/ahmedsat/ebda-cli/config"
	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	EngineerName     string
	Region           string
	Date             time.Time
	CountOfVisits    int
	CountOfPGSAudits int
}

func init() {
	config.MigrationsList = append(config.MigrationsList, &Plan{})
}

type VisitsPlan struct {
	sb strings.Builder
}

// Result implements [main.subcommand].
func (v *VisitsPlan) Result() any {
	panic("unimplemented")
}

// Description implements [main.subcommand].
func (v *VisitsPlan) Description() string {
	return "Visits plan"
}

// Name implements [main.subcommand].
func (v *VisitsPlan) Name() string {
	return "visits-plan"
}

func FormatName(s string) string {
	words := strings.Split(s, "-")
	for i, w := range words {
		r := []rune(w)
		if len(r) == 0 {
			continue
		}
		r[0] = unicode.ToUpper(r[0])
		words[i] = string(r)
	}

	return strings.Join(words, " ")
}

// Run implements [main.subcommand].
func (v *VisitsPlan) Run(args []string) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}

	plansDir := args[0]

	v.sb.WriteString("date\tengName\tregion\tcountOfVisits\tcountOfPGSAudits\n")

	err := filepath.WalkDir(plansDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".tsv") {
			return nil
		}

		region := strings.TrimSuffix(d.Name(), ".tsv")
		region = FormatName(region)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = '\t'

		engineers := []string{}

		names, err := reader.Read()
		if err != nil {
			return err
		}

		for _, name := range names[1:] {
			if name == "" {
				continue
			}

			name := strings.TrimPrefix(name, "م/")
			name = strings.TrimSpace(name)

			engineers = append(engineers, name)
		}

		// remove Header
		_, err = reader.Read()
		if err != nil {
			return err
		}

		records, err := reader.ReadAll()
		if err != nil {
			return err
		}

		lastDate := ""

		for _, record := range records {
			date := record[0]
			if date == "" {
				date = lastDate
			}

			if date == "" {
				return errors.New("empty date")
			}

			if len(record) != 2*len(engineers)+1 {
				return errors.New("invalid record")
			}

			for i, engineer := range engineers {
				countOfVisits := ""
				countOfPGSAudits := ""

				if strings.Contains(record[2*i+2], "زيارة") {
					countOfVisits = record[2*i+1]
				}

				if strings.Contains(record[2*i+2], "PGS") {
					countOfPGSAudits = record[2*i+1]
				}

				if countOfVisits == "" && countOfPGSAudits == "" {
					continue
				}

				fmt.Fprintf(&v.sb, "%s\t%s\t%s\t%s\t%s\n", date, engineer, region, countOfVisits, countOfPGSAudits)
				lastDate = date
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return v.sb.String(), nil
}

// Usage implements [main.subcommand].
func (v *VisitsPlan) Usage() string {
	panic("unimplemented")
}
