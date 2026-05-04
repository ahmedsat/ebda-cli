//go:build !release

package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/utils"
	"github.com/gen2brain/go-fitz"
	"gorm.io/gorm"
)

type SubmissionState struct {
	gorm.Model
	FarmersDone  bool
	SoilDone     bool
	BoundaryDone bool
}

func init() {
	config.MigrationsList = append(config.MigrationsList, &SubmissionState{})
}

type Missing struct {
}

// Usage implements [main.subcommand].
func (m *Missing) Usage() string {
	panic("unimplemented")
}

// Description implements [main.subcommand].
func (m *Missing) Description() string {
	return "Migrating data from kobo to frappe"
}

// Name implements [main.subcommand].
func (m *Missing) Name() string {
	return "missing"
}

// Run implements [main.subcommand].
func (m *Missing) Run(args []string) (err error) {

	file, err := os.OpenFile("rejected.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	flagSet := flag.NewFlagSet("missing", flag.ExitOnError)
	noFarmers := flagSet.Bool("no-farmers", false, "Do not migrate farmers")
	flagSet.Parse(args)
	args = flagSet.Args()

	fmt.Fprintln(os.Stderr, "getting data from kobo...")
	var data []kobo.Collect
	if len(args) > 0 {
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			asset, err := kobo.GetAssetByID[kobo.Collect](id)
			if err != nil {
				return err
			}
			data = append(data, asset)
		}
	} else {
		query := kobo.Query{
			kobo.ValidationStatusKey: nil,
		}

		data, err = kobo.GetAssets[kobo.Collect](query)
		if err != nil {
			return
		}
	}

	runner := utils.NewSyncRunner(10, 100)
	counter := 1
	for _, d := range data {
		runner.Run(func() (err error) {
			var submissionState SubmissionState
			defer func() {
				if err != nil {
					fmt.Printf("%d\t%s\t%s\n", d.ID, d.Farm, err)
					err = nil
				}
				config.DB.Save(&submissionState)
				counter++
			}()

			if d.CollectValidationSate.Label == "Approved" || d.CollectValidationSate.Label == "Not Approved" {
				return
			}

			err = config.DB.First(&submissionState, d.ID).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return
			}

			if err == gorm.ErrRecordNotFound {
				submissionState.ID = uint(d.ID)
				err = config.DB.Create(&submissionState).Error
				if err != nil {
					return
				}
			}

			if d.Code == "" {
				fmt.Fprintf(file, "submission %d has no code\n", d.ID)
				_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusNotApproved)
				if err != nil {
					return
				}
				return
			}

			farms, err := frappe.Get[types.Farm](frappe.Filters{frappe.NewFilter("farm_id", frappe.Eq, d.Code)}, nil, nil)
			if err != nil {
				fmt.Fprintf(file, "submission %d has no farm\n", d.ID)
				_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusNotApproved)
				if err != nil {
					return
				}
				return
			}

			if len(farms) == 0 {
				fmt.Fprintf(file, "submission %d has no farm\n", d.ID)
				_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusNotApproved)
				if err != nil {
					return
				}
				return
			}

			if len(farms) > 1 {
				fmt.Fprintf(file, "submission %d has multiple farms\n", d.ID)
				_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusNotApproved)
				if err != nil {
					return
				}
				return
			}

			d.Farm = farms[0].Name

			err = HandleSoil(&d, &submissionState)
			if err != nil {
				return
			}

			err = HandleBoundary(&d, &submissionState)
			if err != nil {
				return
			}

			err = HandleFarmers(&d, &submissionState, file, *noFarmers)
			if err != nil {
				return
			}

			if submissionState.FarmersDone && submissionState.SoilDone && submissionState.BoundaryDone {
				_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusApproved)
				if err != nil {
					return
				}
			}

			return
		})
	}

	err = runner.Wait()

	time.Sleep(time.Second * 10)

	return
}

func HandleBoundary(collect *kobo.Collect, submissionState *SubmissionState) error {

	if submissionState.BoundaryDone {
		return nil
	}
	area := collect.AreaNew
	if area == "" {
		area = collect.AreaOld
	}
	if area == "" {
		submissionState.BoundaryDone = true
		return nil
	}

	pointsStr := strings.Split(area, ";")

	type point struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	var points []point
	for _, p := range pointsStr {
		p = strings.TrimSpace(p)
		parts := strings.Split(p, " ")
		if len(parts) != 4 {
			return fmt.Errorf("invalid point: %s", p)
		}

		lat, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return err
		}
		lng, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return err
		}

		var p = point{
			Lat: lat,
			Lng: lng,
		}
		points = append(points, p)
	}

	pointsJson, err := json.Marshal(points)
	if err != nil {
		return err
	}

	pointsJson = pointsJson[1 : len(pointsJson)-1]

	farms, err := frappe.Get[types.Farm](frappe.Filters{frappe.NewFilter("farm_id", frappe.Eq, collect.Code)}, nil, nil)
	if err != nil {
		return err
	}

	if len(farms) == 0 {
		return fmt.Errorf("farm not found: %s", collect.Farm)
	}

	if len(farms) > 1 {
		return fmt.Errorf("multiple farms found: %s", collect.Farm)
	}

	farm := farms[0]

	maps, err := frappe.Get[types.MapRecord](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, nil, nil)
	if err != nil {
		return err
	}

	for _, m := range maps {
		if m.Color != "#181818" && m.Color != "" {
			continue
		}

		err := frappe.Delete[types.MapRecord](m.Name)
		if err != nil {
			return err
		}
	}

	m := types.MapRecord{
		Farm:     collect.Farm,
		Jsoncode: string(pointsJson),
		Color:    "#181818",
	}
	_, err = frappe.Create(m)
	if err != nil {
		return err
	}

	submissionState.BoundaryDone = true
	return nil
}

func HandleSoil(collect *kobo.Collect, submissionState *SubmissionState) error {

	if submissionState.SoilDone {
		return nil
	}

	if len(collect.Points) == 0 {
		submissionState.SoilDone = true
		return nil
	}

	s := types.SoilAnalysis{
		Farm:               collect.Farm,
		FarmID:             collect.Code,
		CollectionDatetime: collect.Today,
		NamingSeries:       "Soil-kobo-.YY.-.MM.-",
	}
	features := []utils.Feature{}

	for _, s := range collect.Points {
		lat, lng, _, _, err := s.GeoInfo()
		if err != nil {
			return err
		}
		features = append(features, utils.NewPointFeature(s.PointStr, lat, lng))
	}

	geojson := utils.NewGeoJSON(features...)
	b, err := json.Marshal(geojson)
	if err != nil {
		return err
	}
	s.Location = string(b)

	s, err = frappe.Create(s)
	if err != nil {
		return err
	}

	submissionState.SoilDone = true

	return nil
}

func HandleFarmers(collect *kobo.Collect, submissionState *SubmissionState, log io.Writer, no_farmer bool) error {

	if submissionState.FarmersDone {
		return nil
	}

	if len(collect.Farmers) == 0 {
		submissionState.FarmersDone = true
		return nil
	}

	if no_farmer {
		return nil
	}

	farm, err := frappe.Get1[types.Farm](collect.Farm)
	if err != nil {
		return err
	}

	for i, farmer := range collect.Farmers {

		num, err := strconv.ParseInt(farmer.Number, 10, 64)
		if err != nil {
			return err
		}

		if int(num) > len(farm.Farmers) {
			fmt.Fprintf(log, "%d - farmer number is too big: %d - max: %d\n", submissionState.ID, num, len(farm.Farmers)+1)
			kobo.UpdateValidationState[kobo.Collect](collect.ID, kobo.ValidationStatusNotApproved)
			return nil
		}

		if farmer.Name != farm.Farmers[num-1].FarmerName {
			err = UpdateFarmerName(&farm, num, farmer, log, collect.ID)
			if err != nil {
				return nil
			}
		}

		farm.Farmers[num-1].Phone = farmer.Phone
		farm.Farmers[num-1].NationalIdNumber = farmer.IdNumber

		face := ""
		facePdf := false
		back := ""
		backPdf := false
		for _, att := range collect.Attachments {
			if att.QuestionXpath == fmt.Sprintf("farmers[%d]/id_face", i+1) {
				face = att.DownloadUrl
				switch att.Mimetype {
				case "application/pdf":
					facePdf = true
				case "image/jpeg", "image/png", "image/gif": // do nothing
				default:
					return fmt.Errorf("unknown mimetype: %s", att.Mimetype)
				}
			}
			if att.QuestionXpath == fmt.Sprintf("farmers[%d]/id_back", i+1) {
				back = att.DownloadUrl
				switch att.Mimetype {
				case "application/pdf":
					backPdf = true
				case "image/jpeg", "image/png", "image/gif": // do nothing
				default:
					return fmt.Errorf("unknown mimetype: %s", att.Mimetype)
				}
			}
		}

		var combinedImage image.Image

		reader, err := kobo.Download(face)
		if facePdf {
			combinedImage, err = PdfToImage(reader)
			if err != nil {
				return err
			}
		} else {
			faceImage, _, err := image.Decode(reader)
			if err != nil {
				if errors.Is(err, image.ErrFormat) {
					fmt.Fprintf(log, "unsupported image format: %s\n", face)
				}
				return err
			}
			reader.Close()
			combinedImage = faceImage
		}

		if back != "" {
			reader, err = kobo.Download(back)
			if err != nil {
				return err
			}
			var backImage image.Image
			if backPdf {
				backImage, err = PdfToImage(reader)
				if err != nil {
					return err
				}
			} else {
				backImage, _, err = image.Decode(reader)
				if err != nil {
					if errors.Is(err, image.ErrFormat) {
						fmt.Fprintf(log, "unsupported image format: %s\n", back)
					}
					return err
				}

			}
			reader.Close()

			combinedImage = StackVertical(combinedImage, backImage)
		}

		buf := bytes.Buffer{}
		err = png.Encode(&buf, combinedImage)
		if err != nil {
			return err
		}

		res, err := frappe.UploadFile(
			buf.Bytes(),
			fmt.Sprintf("face-%d-%d.png", collect.ID, num),
			func(sent, total int64) {
				fmt.Fprintf(os.Stderr, "\r[%d:%d]", sent, total)
			})
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr)

		if res.FileUrl == "" {
			return fmt.Errorf("empty file url")
		}

		farm.Farmers[num-1].FarmerNationalIdImage = res.FileUrl
	}

	_, err = frappe.UpdateDoc(farm)
	if err != nil {
		return err
	}

	submissionState.FarmersDone = true

	return nil
}

func UpdateFarmerName(farm *types.Farm, num int64, newFarmer kobo.CollectFarmers, log io.Writer, id int) error {
	if newFarmer.Name == farm.Owner {
		fmt.Fprintf(log, "%d - farmer name is the same as owner %d: %s\n", id, num, newFarmer.Name)
		kobo.UpdateValidationState[kobo.Collect](id, kobo.ValidationStatusNotApproved)
		return fmt.Errorf("%d - farmer name is the same as owner %d: %s\n", id, num, newFarmer.Name)
	}

	if len(strings.Split(newFarmer.Name, " ")) < 3 {
		fmt.Fprintf(log, "%d - farmer name is to short %d: %s\n", id, num, newFarmer.Name)
		kobo.UpdateValidationState[kobo.Collect](id, kobo.ValidationStatusNotApproved)
		return fmt.Errorf("%d - farmer name is to short %d: %s\n", id, num, newFarmer.Name)
	}

	for _, farmer := range farm.Farmers {
		if farmer.FarmerName == newFarmer.Name {
			fmt.Fprintf(log, "%d - farmer name already exists %d: %s\n", id, num, newFarmer.Name)
			kobo.UpdateValidationState[kobo.Collect](id, kobo.ValidationStatusNotApproved)
			return fmt.Errorf("%d - farmer name already exists %d: %s\n", id, num, newFarmer.Name)
		}
	}

	new := types.Farmer{
		FarmerName: strings.TrimSpace(newFarmer.Name),
		Phone:      strings.TrimSpace(newFarmer.Phone),
		Gender:     strings.TrimSpace(newFarmer.Gender),
	}

	_, err := frappe.Create(new)
	if err != nil {
		fmt.Fprintf(log, "%d - filed to create farmer %d: %s error: %s\n", id, num, newFarmer.Name, err)
		kobo.UpdateValidationState[kobo.Collect](id, kobo.ValidationStatusNotApproved)
		return err
	}

	trainings, err := frappe.Get[types.EbdaTraining](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, []string{"name", "farm"}, nil)
	if err != nil {
		return err
	}

	for _, training := range trainings {
		training, err = frappe.Get1[types.EbdaTraining](training.Name)
		if err != nil {
			return err
		}
		for i := range training.Farmers {
			if training.Farmers[i].FarmerName == farm.Farmers[num-1].FarmerName {
				training.Farmers[i].FarmerName = newFarmer.Name
				break
			}
		}
		_, err = frappe.UpdateDoc(training)
		if err != nil {
			return err
		}
	}

	followUps, err := frappe.Get[types.FarmFollowUp](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, frappe.List{"name", "farm"}, nil)
	if err != nil {
		return err
	}

	for _, followUp := range followUps {
		followUp, err = frappe.Get1[types.FarmFollowUp](followUp.Name)
		if err != nil {
			return err
		}
		for i := range followUp.FarmersNames {
			if followUp.FarmersNames[i].Farmer == farm.Farmers[num-1].FarmerName {
				followUp.FarmersNames[i].Farmer = newFarmer.Name
				break
			}
		}
		_, err = frappe.UpdateDoc(followUp)
		if err != nil {
			return err
		}
	}

	for i := range farm.Farmers {
		if farm.Farmers[i].FarmerName == farm.Farmers[num-1].FarmerName {
			farm.Farmers[i].FarmerName = newFarmer.Name
			break
		}
	}

	_, err = frappe.UpdateDoc(farm)
	if err != nil {
		return err
	}

	return nil

}

func StackVertical(img1, img2 image.Image) image.Image {

	if img1 == nil && img2 == nil {
		return nil
	}

	if img1 == nil {
		return img2
	}

	if img2 == nil {
		return img1
	}

	w1, h1 := img1.Bounds().Dx(), img1.Bounds().Dy()
	w2, h2 := img2.Bounds().Dx(), img2.Bounds().Dy()

	// final width is the max width of both images
	width := max(w2, w1)

	// final height is sum of both heights
	height := h1 + h2

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// draw first image (top)
	draw.Draw(dst, image.Rect(0, 0, w1, h1), img1, img1.Bounds().Min, draw.Src)

	// draw second image (bottom)
	draw.Draw(
		dst,
		image.Rect(0, h1, w2, h1+h2),
		img2,
		img2.Bounds().Min,
		draw.Src,
	)

	return dst
}

func PdfToImage(reader io.Reader) (result image.Image, err error) {
	doc, err := fitz.NewFromReader(reader)
	if err != nil {
		return
	}
	defer doc.Close()

	for i := range doc.NumPage() {
		var img *image.RGBA
		img, err = doc.Image(i)
		if err != nil {
			return
		}
		result = StackVertical(result, img)
	}

	return
}
