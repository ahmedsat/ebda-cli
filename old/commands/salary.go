package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

var availableFields = []string{"all", "total_salary", "payment_due", "base_gross", "loans", "salary_items", "benefits"}

type SalaryOptions struct {
	Output
	Fields ListFlagString `json:"fields"`
}

func (s *SalaryOptions) AddFlags(cmd *flag.FlagSet) {
	s.Output.AddFlags(cmd)
	cmd.Var(&s.Fields, "fields", "Fields to get")
}

func (s *SalaryOptions) Validate() (err error) {
	if len(s.Fields) <= 0 {
		return fmt.Errorf("no fields set in --fields\n available fields are \n%v", availableFields)
	}

	for _, f := range s.Fields {
		if f == availableFields[0] {
			s.Fields = availableFields[1:]
			break
		}
		if !slices.Contains(availableFields, f) {
			return fmt.Errorf("field %s is not available \n available fields are \n%v", f, availableFields)
		}
	}

	return s.Output.Validate()
}

const SalaryPath = "/salary-slip"

type Item struct {
	Name   string
	Amount float64
	IsLoss bool
}

type SalaryPage struct {
	TotalSalary float64
	PaymentDue  string
	BaseGross   float64
	Loans       []Item
	SalaryItems []Item
	Benefits    []Item
}

func Salary(opt SalaryOptions) (err error) {
	err = opt.Validate()
	if err != nil {
		return
	}

	res, err := utils.Get(SalaryPath)
	if err != nil {
		return
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	html := string(bytes)
	page, err := parseSalaryHTML(html)
	if err != nil {
		log.Fatal(err)
	}

	t := utils.TableBase{}
	for _, f := range opt.Fields {
		switch f {
		case "total_salary":
			t.AppendRow("Total Salary", fmt.Sprintf("%v", page.TotalSalary))
			t.AppendRow()
		case "payment_due":
			t.AppendRow("Payment Due", fmt.Sprintf("%v", page.PaymentDue))
			t.AppendRow()
		case "base_gross":
			t.AppendRow("Base Gross", fmt.Sprintf("%v", page.BaseGross))
			t.AppendRow()
		case "loans":
			t.AppendRow("Loans")
			t.AppendRow("Name", "Amount")
			for _, l := range page.Loans {
				t.AppendRow(l.Name, fmt.Sprintf("%v", l.Amount))
			}
			t.AppendRow()
		case "salary_items":
			t.AppendRow("Salary Items")
			t.AppendRow("Name", "Amount", "Is Loss")
			for _, l := range page.SalaryItems {
				t.AppendRow(l.Name, fmt.Sprintf("%v", l.Amount), fmt.Sprintf("%v", l.IsLoss))
			}
			t.AppendRow()
		case "benefits":

			t.AppendRow("Benefits")
			t.AppendRow("Name", "Amount")
			for _, l := range page.Benefits {
				t.AppendRow(l.Name, fmt.Sprintf("%v", l.Amount))
			}
		}
	}
	opt.Output.Print(opt.TablePrinter(&t))

	return
}

func parseSalaryHTML(html string) (*SalaryPage, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	page := &SalaryPage{}

	// Clean number: "EGP 8,709.505" → 8709.505
	number := func(s string) float64 {
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "EGP", "")
		s = strings.ReplaceAll(s, ",", "")
		s = strings.TrimSpace(s)
		f, _ := strconv.ParseFloat(s, 64)
		return f
	}

	// -------- Total Salary --------
	page.TotalSalary = number(doc.Find(".earn").First().Text())

	// -------- Payment Due --------
	doc.Find("div.card-body .h3").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(strings.ToLower(s.Text()), "september") ||
			regexp.MustCompile(`[A-Za-z]+\s+\d+`).MatchString(s.Text()) {
			page.PaymentDue = strings.TrimSpace(s.Text())
		}
	})

	// -------- Base Gross --------
	doc.Find(".card-body .h5").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Gross") {
			parts := strings.Split(s.Text(), ":")
			if len(parts) == 2 {
				page.BaseGross = number(parts[1])
			}
		}
	})

	// -------- Loans table --------
	doc.Find("table").First().Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		tds := tr.Find("td")
		if tds.Length() == 2 {
			page.Loans = append(page.Loans, Item{
				Name:   strings.TrimSpace(tds.Eq(0).Text()),
				Amount: number(tds.Eq(1).Text()),
			})
		}
	})

	// -------- Salary Components --------
	doc.Find("#salary-slip table tbody tr").Each(func(i int, tr *goquery.Selection) {
		tds := tr.Find("td")
		if tds.Length() == 2 {
			amountText := tds.Eq(1).Text()

			item := Item{
				Name:   strings.TrimSpace(tds.Eq(0).Text()),
				Amount: number(amountText),
				IsLoss: strings.Contains(tr.Find("span").AttrOr("class", ""), "loss"),
			}

			page.SalaryItems = append(page.SalaryItems, item)
		}
	})

	// -------- Employee Benefits --------
	doc.Find("div.card.hidden table tbody tr").Each(func(i int, tr *goquery.Selection) {
		tds := tr.Find("td")
		if tds.Length() == 2 {
			page.Benefits = append(page.Benefits, Item{
				Name:   strings.TrimSpace(tds.Eq(0).Text()),
				Amount: number(tds.Eq(1).Text()),
			})
		}
	})

	return page, nil
}
