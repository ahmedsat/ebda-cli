//go:build !release

package commands

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
	"gorm.io/gorm"
)

const dateLayout = "2/Jan/2006"

var EngToRegion = map[string]string{
	"ابراهيم شاذلى":        "aswan",
	"حاتم محمد":            "aswan",
	"محمد أحمد محمد":       "aswan",
	"يوسف صبرى":            "aswan",
	"أحمد عبدالله":         "aswan",
	"ابراهيم عبدالمعطى":    "luxor",
	"وائل عبدالعظيم":       "luxor",
	"عبدالرحمن محمد":       "luxor",
	"فارس ابراهيم عزيز":    "luxor",
	"محمد حسب النبى":       "luxor",
	"مصطفى حسن":            "luxor",
	"على جمال":             "qena",
	"مصطفى عبدالباسط":      "qena",
	"حسان عاطف":            "qena",
	"مصطفى حمام":           "sohag",
	"أحمد محمد عمر":        "sohag",
	"خالد ناصر":            "sohag",
	"على أحمد":             "sohag",
	"حسن السيد":            "sohag",
	"عمر الخطاب":           "sohag",
	"محمد جمال":            "assuit",
	"ياسر سيد":             "assuit",
	"محمد جمال محمد":       "assuit",
	"وليد عيد أحمد":        "assuit",
	"وسيم ابراهيم":         "assuit",
	"أحمد مجدى":            "new-valley",
	"حازم عبدالرحمن محمود": "new-valley",
	"أحمد محمد معوض":       "new-valley",
	"أحمد عزت":             "new-valley",
	"محمد رفعت":            "new-valley",
	"باهر نبيل":            "minya",
	"أنور بهاء":            "minya",
	"مرسى رضا":             "minya",
	"محمد أحمد":            "minya",
	"ضاحى محروص":           "minya",
	"ونجد ماجد":            "minya",
	"حازم مصطفى":           "minya",
	"محمد قطب":             "beni-suef",
	"عبدالله رضوان":        "beni-suef",
	"هشام عمران":           "beni-suef",
	"محمد خليفه":           "beni-suef",
	"اسلام عبدالناصر":      "beni-suef",
	"عبدالرحمن قرنى":       "beni-suef",
	"محمد سلامه":           "faiyum",
	"اسلام محمود رمضان":    "faiyum",
	"محمد عبدالحميد عطية":  "faiyum",
	"محمود عبدالتواب":      "faiyum",
	"زياد خالد":            "faiyum",
	"مصطفى محسن":           "faiyum",
	"محمود محسن":           "faiyum",
	"اسلام نادر":           "el-giza",
	"جورج ذكى":             "el-giza",
	"أحمد عبدالتواب":       "el-giza",
	"عبدالله سليمان":       "el-giza",
	"سيد أحمد حسنى":        "el-giza",
	"محمد عبدالله فوزى":    "el-wahat-el-baharia",
	"عيسى الشرنوبى":        "el-wahat-el-baharia",
	"اسامة محمد الصغير":    "el-wahat-el-baharia",
	"وليد سامى":            "el-wahat-el-baharia",
	"محمد ربيع":            "el-wahat-el-baharia",
	"محمد مخيمر":           "el-wahat-el-baharia",
	"محمد على عبدالرازق":   "qalyubia",
	"احمد سعيد":            "qalyubia",
	"على عطالله":           "sharqia",
	"ممدوح اكرم":           "sharqia",
	"محمود السيد فتحى":     "sharqia",
	"محمود عبد الفتاح":     "sharqia",
	"فارس تامر":            "sharqia",
	"محمد سماحه":           "ismailya",
	"احمد علي الفتاح":      "ismailya",
	"محمد ياسر الزفزافى":   "menoufia",
	"مصطفى على رجب":        "menoufia",
	"طارق مشحوت":           "menoufia",
	"طاهر حسن":             "menoufia",
	"على فرج":              "kafr-el-shiekh",
	"محمد عبدالرحيم":       "kafr-el-shiekh",
	"محمود شعبان":          "kafr-el-shiekh",
	"اسلام جمال":           "kafr-el-shiekh",
	"عبدالعظيم محمد":       "al-behiera",
	"وليد عرام":            "al-behiera",
	"بشار الصافى":          "al-behiera",
	"حسام ايمن":            "al-behiera",
	"محمد بدوى":            "al-behiera",
	"أشرف النجار":          "al-behiera",
	"ابراهيم المهدى":       "al-behiera",
	"محمد احمد الفقى":      "gharbia",
	"أحمد ياسر المنسى":     "gharbia",
	"عبدالحميد الجنجيهى":   "gharbia",
	"محمد عادل":            "gharbia",
	"عبدالله السيد":        "dakahlia",
	"أحمد عصام":            "damietta",
	"يوسف منصور":           "damietta",
	"محمد صبحى":            "alexandria-desert",
	"صالح ماضى":            "alexandria-desert",
	"محمد منير":            "alexandria-desert",
	"محمد إسماعيل جاتوه":   "siwa",
	"مصطفي عمر بن إدريس":   "siwa",
	"عبد الرحمن مفتاح":     "marsa-matruh",
	"احمد عيسي":            "marsa-matruh",
	"يوسف جلال":            "marsa-matruh",
	"إبراهيم يسري":         "marsa-matruh",
}

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

// Run implements [main.subcommand].
func (v *VisitsPlan) Run(args []string) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	subcommand := args[0]

	switch subcommand {
	case "new":
		if len(args) < 2 {
			return nil, fmt.Errorf("not enough arguments")
		}
		return nil, v.New(args[1:])
	case "export":
		return nil, v.Export()
	default:
		return nil, fmt.Errorf("unavailable commands: %s", subcommand)
	}
}

func (v *VisitsPlan) Export() error {
	v.sb.Reset()

	var plans []Plan

	err := config.DB.Find(&plans).Error
	if err != nil {
		return err
	}

	fmt.Fprintf(&v.sb, "%s\t%s\t%s\t%s\t%s\n",
		"Engineer Name", "Region", "Date", "Count Of Visits", "Count Of PGS Audits")

	for _, plan := range plans {

		fmt.Fprintf(&v.sb, "%s\t%s\t%s\t%d\t%d\n",
			plan.EngineerName, plan.Region, plan.Date.Format("2006-01-02"), plan.CountOfVisits, plan.CountOfPGSAudits)
	}

	return nil
}

func (v *VisitsPlan) New(files []string) error {

	for _, filename := range files {

		base := filepath.Base(filename)
		if base == "" {
			return fmt.Errorf("empty file name")
		}

		region, ok := strings.CutSuffix(base, ".tsv")
		if !ok {
			return fmt.Errorf("file %s is not a csv file", filename)
		}

		fmt.Fprintf(os.Stderr, "Processing %s\n", filename)

		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = '\t'

		header, err := reader.Read()
		if err != nil {
			return err
		}
		engineers := []string{}

		for i := 1; i < len(header); i += 2 {
			name := strings.TrimSpace(header[i])
			if name == "" {
				return fmt.Errorf("empty engineer name at column %d", i)
			}

			if EngToRegion[name] != region {
				return fmt.Errorf("engineer %s is not in region %s", name, region)
			}

			engineers = append(engineers, name)
		}

		var currentDate time.Time

		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// minimum sanity check
			if len(record) < 1 {
				continue
			}

			dateStr := strings.TrimSpace(record[0])

			if dateStr != "" {
				d, err := time.Parse(dateLayout, dateStr)
				if err != nil {
					return fmt.Errorf("invalid date %q", dateStr)
				}
				currentDate = d
			} else {
				if currentDate.IsZero() {
					return fmt.Errorf("empty date with no previous date")
				}
			}

			expected := 1 + len(engineers)*2
			if len(record) != expected {
				return fmt.Errorf(
					"invalid record length %d, expected %d",
					len(record), expected,
				)
			}

			for i, engineer := range engineers {
				countStr := strings.TrimSpace(record[1+i*2])
				action := strings.TrimSpace(record[2+i*2])

				if countStr == "" {
					continue
				}

				count, err := strconv.Atoi(countStr)
				if err != nil {
					return fmt.Errorf("invalid count for %s on %s", engineer, currentDate)
				}

				plan := Plan{
					EngineerName: engineer,
					Date:         currentDate,
					Region:       EngToRegion[engineer],
				}

				switch {
				case strings.Contains(action, "تفتيش"):
					plan.CountOfPGSAudits = count
				case strings.Contains(action, "زيارة"):
					plan.CountOfVisits = count
				default:
					// تدريب / جمعة / etc → ignored numerically
					continue
				}

				if err := config.DB.Create(&plan).Error; err != nil {
					return err
				}
			}
		}

	}
	return nil
}

// Usage implements [main.subcommand].
func (v *VisitsPlan) Usage() string {
	panic("unimplemented")
}
