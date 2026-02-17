package commands

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/utils"
	"github.com/atotto/clipboard"
)

//go:embed farms.tsv
var farms string
var farmsMap map[string][]string

func init() {
	farmsMap = map[string][]string{}
	for line := range strings.SplitSeq(farms, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		farmsMap[fields[0]] = fields[1:]
	}
}

type Farm struct{}

// Description implements [main.subcommand].
func (f *Farm) Description() string {
	return "information about farms"
}

// Name implements [main.subcommand].
func (f *Farm) Name() string {
	return "farms"
}

// Result implements [main.subcommand].
func (f *Farm) Result() any {
	return nil
}

// Run implements [main.subcommand].
func (f *Farm) Run(args []string) (result any, err error) {

	runner := utils.NewSyncRunner(100, 0)

	firstDay := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)

	fs := flag.NewFlagSet("farms", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	pgs := fs.Bool("pgs", false, "Extract PGS data from kobo")
	followUp := fs.Bool("follow-up", false, "Follow up")
	farmers := fs.Bool("farmers", false, "Farmers")
	start := fs.String("start", firstDay.Format("2-1-2006"), "Start date")
	fs.Parse(args)

	args = fs.Args()
	filters := frappe.Filters{
		frappe.NewFilter("type", frappe.Eq, "Farm"),
		frappe.NewFilter("farm_status", frappe.Neq, "Cancelled"),
	}

	if len(args) > 0 {
		filters = append(filters, frappe.NewFilter("farm_id", frappe.In, frappe.FiltersValueList(args...)))
	}

	fmt.Fprintln(os.Stderr, "getting data from frappe...")
	farms, err := frappe.Get[types.Farm](filters, frappe.List{
		"farm_id",
		"name",
	})

	var out utils.SyncIoWriter
	out.Writer = os.Stdout

	if *copy {
		sb := strings.Builder{}
		out.Writer = &sb
		defer func() {
			clipboard.WriteAll(sb.String())
		}()
	}

	fmt.Fprintf(&out, "%s", strings.Join([]string{
		"farm code", "farm name", "arabic_name", "region", "eng", "group",
	}, "\t"))

	if *farmers {
		fmt.Fprintf(&out, "\tcount of farmers\thas id\thas phone")
	}

	var followUps []types.FarmFollowUp
	if *followUp {
		var d time.Time
		d, err = time.Parse("2-1-2006", *start)
		if err != nil {
			return
		}
		fmt.Fprintln(os.Stderr, "getting follow-up data from frappe...")
		fmt.Fprint(&out, "\tcount of visits\trate of visits")
		followUps, err = frappe.Get[types.FarmFollowUp](frappe.Filters{
			frappe.NewFilter("visit_date", frappe.Gte, d.Format("2006-01-02")),
		}, frappe.List{"name", "farm", "farm_code"})
		if err != nil {
			return
		}
	}

	var pgsSubmissions []kobo.PGSNew
	if *pgs {
		fmt.Fprintln(os.Stderr, "getting pgs data from kobo...")
		fmt.Fprintf(&out, "\tPGS Count\tApproved\tRejected\tPending")
		pgsSubmissions, err = kobo.GetAssets[kobo.PGSNew](kobo.Query{
			"at_house/visit_date": kobo.Query{
				"$gte": firstDay.Format("2006-01-02"),
			},
		})
		if err != nil {
			return
		}
	}
	fmt.Fprintf(&out, "\n")

	fmt.Fprintln(os.Stderr, "start evaluating farms")

	for i, farm := range farms {
		runner.Run(func() (err error) {
			if _, ok := farmsMap[farm.FarmId]; !ok {
				return
			}
			sb := strings.Builder{}
			defer func() {
				fmt.Fprintln(&out, sb.String())
				n, _ := utils.NewProgressNotification(
					"Farms Report",
					fmt.Sprintf("Getting Data from frappe [%d:%d]", i+1, len(farms)),
					"farms-report", int((i+1)*100/len(farms)),
				)
				n.Run()
			}()

			fmt.Fprintf(&sb, "%s\t%s", farm.FarmId, strings.Join(farmsMap[farm.FarmId], "\t"))

			if *farmers {

				farm, err := frappe.Get1[types.Farm](farm.Name)
				if err != nil {
					return err
				}

				var countOfFarmers, hasId, hasPhone int
				for _, farmer := range farm.Farmers {
					countOfFarmers++
					if farmer.NationalIdNumber != "" {
						hasId++
					}
					if farmer.Phone != "" {
						hasPhone++
					}
				}

				fmt.Fprintf(&sb, "\t%d\t%d\t%d", countOfFarmers, hasId, hasPhone)
			}

			if *followUp {
				rate := 0.0
				count := 0
				for _, followUp := range followUps {
					if followUp.Farm == farm.Name {
						err = followUp.Rate()
						if err != nil {
							return
						}
						count++
						rate += followUp.RatePercent
					}
				}
				if count != 0 {
					fmt.Fprintf(&sb, "\t%d\t%.2f", count, rate/float64(count))
				} else {
					fmt.Fprintf(&sb, "\t0\t0")
				}
			}

			if *pgs {
				count := 0
				approved := 0
				rejected := 0
				pending := 0
				for _, s := range pgsSubmissions {
					if s.FarmCode == farm.FarmId {
						count++
						switch s.ValidationStatus.Label {
						case "Approved":
							approved++
						case "Not Approved":
							rejected++
						default:
							pending++
						}
					}
				}
				if count != 0 {
					fmt.Fprintf(&sb, "\t%d\t%d\t%d\t%d", count, approved, rejected, pending)
				} else {
					fmt.Fprintf(&sb, "\t0\t0\t0\t0")
				}
			}

			return
		})
	}

	return nil, runner.Wait()
}

// Usage implements [main.subcommand].
func (f *Farm) Usage() string {
	panic("unimplemented")
}
