package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ahmedsat/erp-reports-cli/utils"
	"github.com/atotto/clipboard"
)

type Dates struct {
	FromStr    string `json:"from"`
	ToStr      string `json:"to"`
	DateLayout string `json:"date_layout"`

	From time.Time `json:"-"`
	To   time.Time `json:"-"`
}

func (d *Dates) AddFlags(cmd *flag.FlagSet) {
	cmd.StringVar(&d.FromStr, "from", "1-1-2022", "Date part of ISO")
	cmd.StringVar(&d.ToStr, "to", time.Now().Format("2-1-2006"), "Date part of ISO")
	cmd.StringVar(&d.DateLayout, "date-layout", "2-1-2006", "Date layout")
}

func (d *Dates) Validate() (err error) {
	if d.DateLayout == "" {
		return errors.New("-date-layout not set")
	}

	// check if the layout is valid
	// bug: this does not check the layout
	_, err = time.Parse(d.DateLayout, "2-1-2006")
	if err != nil {
		return
	}

	if d.FromStr == "" {
		return errors.New("-from not set")
	}

	if d.ToStr == "" {
		return errors.New("-to not set")
	}

	d.From, err = time.Parse(d.DateLayout, d.FromStr)
	if err != nil {
		return
	}

	d.To, err = time.Parse(d.DateLayout, d.ToStr)
	if err != nil {
		return
	}

	if d.From.After(d.To) {
		return errors.New("-from is after -to")
	}

	return
}

type Output struct {
	Format     string `json:"format"`
	Copy       bool   `json:"copy"`
	OutputFile string `json:"output"`

	TablePrinter utils.TablePrinter `json:"-"`
}

func (o *Output) AddFlags(cmd *flag.FlagSet) {
	cmd.StringVar(&o.Format, "format", "csv", "Format of output [csv, tsv]")
	cmd.BoolVar(&o.Copy, "copy", false, "Copy to clipboard")
	cmd.StringVar(&o.OutputFile, "output", "", "Output file")
}

func (o *Output) Validate() (err error) {
	switch o.Format {
	case "csv":
		o.TablePrinter = utils.TablePrinterCsv
	case "tsv":
		o.TablePrinter = utils.TablePrinterTsv
	default:
		o.TablePrinter = utils.TablePrinterCsv
		fmt.Printf("unknown format: %s defaulting to csv\n", o.Format)
	}

	return nil
}

func (o *Output) Print(data string) (err error) {

	if o.Copy {
		err = clipboard.WriteAll(data)
		if err != nil {
			return
		}
	}

	if o.OutputFile != "" {
		err = os.WriteFile(o.OutputFile, []byte(data), 0644)
		if err != nil {
			return
		}
	}

	if o.OutputFile == "" && !o.Copy {
		fmt.Println(data)
	}

	return
}

type GlobalOptions struct {
	Dates
	Output

	IncludeCanceled bool `json:"include_canceled"`
}

func (g *GlobalOptions) AddFlags(cmd *flag.FlagSet) {
	g.Dates.AddFlags(cmd)
	g.Output.AddFlags(cmd)

	cmd.BoolVar(&g.IncludeCanceled, "include-canceled", false, "Include canceled farms")
}

func (g *GlobalOptions) Validate() (err error) {

	err = g.Dates.Validate()
	if err != nil {
		return
	}

	err = g.Output.Validate()
	if err != nil {
		return
	}

	return nil
}
