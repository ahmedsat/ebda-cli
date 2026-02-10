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

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/utils"
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

var userRegionMap = map[string]string{}

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

// Result implements [main.subcommand].
func (m *Missing) Result() any {
	return nil
}

// Run implements [main.subcommand].
func (m *Missing) Run(args []string) error {

	file, err := os.OpenFile("rejected.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	flagSet := flag.NewFlagSet("missing", flag.ExitOnError)
	fix := flagSet.Bool("fix", false, "Fix not approved submissions")
	flagSet.Parse(args)

	fmt.Fprintln(os.Stderr, "getting data from kobo...")
	data, err := kobo.GetAssets[kobo.Collect]()
	if err != nil {
		return err
	}

	if *fix {
		fmt.Fprintln(os.Stderr, "fixing not approved submissions...")
		for _, d := range data {
			if d.CollectValidationSate.Label == "Not Approved" {
				fmt.Println(d.ID)
				fmt.Println(d.Code)
				res, err := kobo.GetUpdateURL[kobo.Collect](d.ID)
				if err != nil {
					return err
				}

				fmt.Printf("%+v", res)
				// todo: ask for confirmation if confirmed update validation states
				fmt.Scanln()
			}
		}
		return nil
	}

	for i, d := range data {
		fmt.Fprintf(os.Stderr, "\rProgress {%d} [%d:%d] (%.2f%%)", d.ID, i+1, len(data), float64(i+1)/float64(len(data))*100)

		if d.CollectValidationSate.Label == "Approved" || d.CollectValidationSate.Label == "Not Approved" {
			continue
		}
		var submissionState SubmissionState
		err = config.DB.First(&submissionState, d.ID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == gorm.ErrRecordNotFound {
			submissionState.ID = uint(d.ID)
			err = config.DB.Create(&submissionState).Error
			if err != nil {
				return err
			}
		}

		if d.Farm == "" {
			fmt.Fprintf(file, "submission %d has no farm\n", d.ID)
			_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusNotApproved)
			if err != nil {
				return err
			}
			continue
		}

		err = HandleSoil(&d, &submissionState)
		if err != nil {
			return err
		}

		err = HandleBoundary(&d, &submissionState)
		if err != nil {
			return err
		}

		err = HandleFarmers(&d, &submissionState, file)
		if err != nil {
			return err
		}

		if submissionState.FarmersDone && submissionState.SoilDone && submissionState.BoundaryDone {
			_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusApproved)
			if err != nil {
				return err
			}
		}

		config.DB.Save(&submissionState)
	}
	fmt.Fprintln(os.Stderr, "")

	return nil
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
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	}

	var points []point
	for _, p := range pointsStr {
		p = strings.TrimSpace(p)
		parts := strings.Split(p, " ")
		if len(parts) != 4 {
			return fmt.Errorf("invalid point: %s", p)
		}
		var p = point{
			Lat: parts[0],
			Lng: parts[1],
		}
		points = append(points, p)
	}

	pointsJson, err := json.Marshal(points)
	if err != nil {
		return err
	}

	pointsJson = pointsJson[1 : len(pointsJson)-1]

	maps, err := frappe.Get[types.MapRecord](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, collect.Farm)}, nil)
	if err != nil {
		return err
	}

	for _, m := range maps {
		if m.Color != "#181818" {
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

func HandleFarmers(collect *kobo.Collect, submissionState *SubmissionState, log io.Writer) error {

	if submissionState.FarmersDone {
		return nil
	}

	if len(collect.Farmers) == 0 {
		submissionState.FarmersDone = true
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
		back := ""
		for _, att := range collect.Attachments {
			if att.QuestionXpath == fmt.Sprintf("farmers[%d]/id_face", i+1) {
				face = att.DownloadUrl
			}
			if att.QuestionXpath == fmt.Sprintf("farmers[%d]/id_back", i+1) {
				back = att.DownloadUrl
			}
		}

		reader, err := kobo.Download(face)
		faceImage, _, err := image.Decode(reader)
		if err != nil {
			if errors.Is(err, image.ErrFormat) {
				fmt.Fprintf(log, "unsupported image format: %s\n", face)
			}
			return err
		}
		reader.Close()

		var combinedImage = faceImage

		if back != "" {

			reader, err = kobo.Download(back)
			if err != nil {
				return err
			}

			backImage, _, err := image.Decode(reader)
			if err != nil {
				if errors.Is(err, image.ErrFormat) {
					fmt.Fprintf(log, "unsupported image format: %s\n", back)
				}
				return err
			}
			reader.Close()

			combinedImage = StackVertical(faceImage, backImage)
		}
		buf := bytes.Buffer{}
		err = png.Encode(&buf, combinedImage)
		if err != nil {
			return err
		}

		res, err := frappe.UploadFile(buf.Bytes(), fmt.Sprintf("face-%d-%d.jpg", collect.ID, num), func(sent, total int64) {
			percent := float64(sent) / float64(total) * 100
			fmt.Fprintf(os.Stderr, "\rUploading: %.1f%%", percent)
		})
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr)

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
		FarmerName: newFarmer.Name,
		Phone:      newFarmer.Phone,
		Gender:     newFarmer.Gender,
	}

	_, err := frappe.Create(new)
	if err != nil {
		fmt.Fprintf(log, "%d - filed to create farmer %d: %s error: %s\n", id, num, newFarmer.Name, err)
		kobo.UpdateValidationState[kobo.Collect](id, kobo.ValidationStatusNotApproved)
		return err
	}

	trainings, err := frappe.Get[types.EbdaTraining](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, []string{"name", "farm"})
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

	followUps, err := frappe.Get[types.FarmFollowUp](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, frappe.List{"name", "farm"})
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
