package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

type FarmsOptions struct {
	GlobalOptions
	Fields                 ListFlagString `json:"fields"`
	FarmFields             []string
	FarmApplicationsFields []string
}

func (f *FarmsOptions) AddFlags(cmd *flag.FlagSet) {
	f.GlobalOptions.AddFlags(cmd)

	cmd.Var(&f.Fields, "fields", "Fields to get")
}

var validFarmsFields = []string{
	"all",
	"farm_id",
	"name",
	"arabic_name",
	"region",
	"total_farmers",
	"farm_area__feddan",
	"farm_application",
	"creation_date",
	"latitude",
	"longitude",
}

var validFarmApplicationsFields = []string{
	"name",
	"engineer_name",
	"user_name",
}

func tryCorrectField(field string) string {
	switch field {
	case "ar_name":
		return "arabic_name"
	case "a_engineer":
		return "a_engineer_name"
	case "engineer":
		return "engineer_name"
	case "code":
		return "farm_id"
	case "area":
		return "farm_area__feddan"
	case "total":
		return "total_farmers"
	case "application":
		return "farm_application"
	case "date":
		return "creation_date"
	case "lat":
		return "latitude"
	case "lng", "long":
		return "longitude"
	default:
		return field
	}
}

func (f *FarmsOptions) Validate() (err error) {
	err = f.GlobalOptions.Validate()
	if err != nil {
		return
	}

	if len(f.Fields) == 0 {
		return fmt.Errorf("no fields set in --fields\n available fields are %v", slices.Concat(validFarmsFields, validFarmApplicationsFields))
	}

	for i := range f.Fields {

		if f.Fields[i] == "all" {
			f.FarmFields = validFarmsFields[1:]
			f.FarmApplicationsFields = validFarmApplicationsFields
			f.Fields = append(f.FarmFields, func(in []string) (out []string) {
				for _, f := range in {
					out = append(out, strings.Join([]string{"a_", f}, ""))
				}
				return
			}(f.FarmApplicationsFields)...)
			break
		}

		f.Fields[i] = tryCorrectField(f.Fields[i])

		if after, ok := strings.CutPrefix(f.Fields[i], "a_"); ok {
			withoutPrefix := after
			// check if the field is valid
			if !slices.Contains(validFarmApplicationsFields, withoutPrefix) {
				return fmt.Errorf("invalid field a_%s", withoutPrefix)
			}
			f.FarmApplicationsFields = append(f.FarmApplicationsFields, withoutPrefix)
		} else {
			// check if the field is valid
			if !slices.Contains(validFarmsFields, f.Fields[i]) {
				return fmt.Errorf("invalid field %s", f.Fields[i])
			}
			f.FarmFields = append(f.FarmFields, f.Fields[i])
		}
	}

	if len(f.FarmFields) == 0 {
		return errors.New("no farm fields set")
	}

	return
}

func (f FarmApplication) GetField(field string) string {

	switch field {
	case "name":
		return f.Name
	case "engineer_name":
		return f.Engineer
	case "user_name":
		return f.UserName
	default:
		fmt.Fprintf(os.Stderr, "invalid field %s\n", field)
		return ""
	}
}

type Farm struct {
	Name            string  `json:"name"`
	ArabicName      string  `json:"arabic_name"`
	Region          string  `json:"region"`
	TotalFarmers    int     `json:"total_farmers"`
	Area            float64 `json:"farm_area__feddan"`
	Code            string  `json:"farm_id"`
	Application     string  `json:"farm_application"`
	CreationDateStr string  `json:"creation_date"`
	Latitude        string  `json:"latitude"`
	Longitude       string  `json:"longitude"`

	FarmApplication FarmApplication `json:"-"`

	// parsed
	CreationDate time.Time `json:"-"`
}

func (Farm) DocTypeName() string { return "Farm" }

func (f Farm) GetField(field string) string {

	if after, ok := strings.CutPrefix(field, "a_"); ok {
		return f.FarmApplication.GetField(after)
	}

	switch field {
	case "arabic_name":
		return f.ArabicName
	case "name":
		return f.Name
	case "region":
		return f.Region
	case "total_farmers":
		return fmt.Sprintf("%d", f.TotalFarmers)
	case "farm_area__feddan":
		return fmt.Sprintf("%.2f", f.Area)
	case "farm_id":
		return f.Code
	case "farm_application":
		return f.Application
	case "creation_date":
		return f.CreationDate.Format("01-02-2006")
	case "latitude":
		return f.Latitude
	case "longitude":
		return f.Longitude
	default:
		fmt.Fprintf(os.Stderr, "invalid field %s\n", field)
		return "UNKNOWN: " + field
	}
}

func Farms(opt FarmsOptions) (err error) {

	err = opt.Validate()
	if err != nil {
		return
	}

	// farms, err := erp.GetFarms[Farm](erp.FarmsOptions{
	// 	From: opt.From,
	// 	To:   opt.To,
	// 	Fields: utils.Ternary(
	// 		len(opt.FarmApplicationsFields) > 0,
	// 		append(opt.FarmFields, "farm_application"),
	// 		opt.FarmFields),
	// 	IncludeCanceled: opt.IncludeCanceled,
	// })

	filters := utils.Filters{
		utils.NewFilter("type", utils.Eq, "farm"),
		utils.NewFilter("creation_date", utils.Gte, opt.From.Format("2006-01-02")),
		// ? to is excluded, so we have to add 1 day
		utils.NewFilter("creation_date", utils.Lte, opt.To.AddDate(0, 0, 1).Format("2006-01-02")),
	}

	if !opt.IncludeCanceled {
		filters = append(filters, utils.NewFilter("farm_status", utils.Neq, "Cancelled"))
	}

	farms, err := erp.Get[Farm](filters, append(opt.FarmFields, "farm_application"))

	if err != nil {
		return
	}

	var farmApplications = []FarmApplication{}

	if len(opt.FarmApplicationsFields) > 0 {
		farmApplications, err = erp.Get[FarmApplication](utils.Filters{}, append(opt.FarmApplicationsFields, "name"))
		if err != nil {
			return
		}

		for i := range farms {
			farms[i].FarmApplication = utils.FindF(farmApplications, func(application FarmApplication) bool {
				return application.Name == farms[i].Application
			})
		}
	}

	t := utils.TableBase{}
	t.SetHeader(opt.Fields...)

	for _, farm := range farms {
		var row []string
		for _, field := range opt.Fields {
			if field == "creation_date" {
				farm.CreationDate, err = time.Parse("2006-01-02 15:04:05", farm.CreationDateStr)
				if err != nil {
					return
				}
			}
			row = append(row, farm.GetField(field))
		}
		t.AppendRow(row...)
	}

	data := opt.TablePrinter(&t)
	opt.Print(data)

	return
}
