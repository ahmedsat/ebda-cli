package kobo

import (
	"fmt"
	"strconv"
	"strings"
)

type CollectAttachments struct {
	DownloadUrl       string `json:"download_url"`
	Mimetype          string `json:"mimetype"`
	Filename          string `json:"filename"`
	MediaFileBasename string `json:"media_file_basename"`
	Uid               string `json:"uid"`
	IsDeleted         bool   `json:"is_deleted"`
	DownloadLargeUrl  string `json:"download_large_url"`
	DownloadMediumUrl string `json:"download_medium_url"`
	DownloadSmallUrl  string `json:"download_small_url"`
	QuestionXpath     string `json:"question_xpath"`
}

type CollectFarmers struct {
	IntegerTu0fc18 string `json:"farmers/integer_tu0fc18"`
	Name           string `json:"farmers/name"`
	Gender         string `json:"farmers/gender"`
	IdFace         string `json:"farmers/id_face"`
	IdBack         string `json:"farmers/id_back"`
	IdNumber       string `json:"farmers/id_number"`
	Phone          string `json:"farmers/phone"`
}

type CollectPoint struct {
	PointStr string `json:"points/point"`
}

func (cp CollectPoint) GeoInfo() (lat, lng, alt, acc float64, err error) {
	parts := strings.Split(cp.PointStr, " ")
	if len(parts) != 4 {
		err = fmt.Errorf("Wrong amount of parts, expect (4) got (%d)", len(parts))
		return
	}

	lat, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return
	}
	lng, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return
	}
	alt, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return
	}
	acc, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return
	}

	return
}

type CollectValidationSate struct {
	Timestamp int    `json:"timestamp"`
	Uid       string `json:"uid"`
	ByWhom    string `json:"by_whom"`
	Label     string `json:"label"`
}

type Collect struct {
	ID               int    `json:"_id"`
	Start            string `json:"start"`
	End              string `json:"end"`
	Today            string `json:"today"`
	Username         string `json:"username"`
	Deviceid         string `json:"deviceid"`
	Phonenumber      string `json:"phonenumber"`
	Farm             string `json:"farm"`
	Code             string `json:"text_ej0vt79"`
	MetaInstanceID   string `json:"meta/instanceID"`
	MetaRootUuid     string `json:"meta/rootUuid"`
	MetaDeprecatedID string `json:"meta/deprecatedID"`
	FormhubUuid      string `json:"formhub/uuid"`
	Version          string `json:"__version__"`
	XformIdString    string `json:"_xform_id_string"`
	Uuid             string `json:"_uuid"`
	Status           string `json:"_status"`
	SubmissionTime   string `json:"_submission_time"`
	SubmittedBy      string `json:"_submitted_by"`
	AreaOld          string `json:"area"`
	AreaNew          string `json:"group_rl8yk95/area"`
	Calculation      string `json:"group_rl8yk95/calculation"`
	CalculatedArea   string `json:"group_rl8yk95/calculated_area"`

	CollectValidationSate `json:"_validation_status"`

	Farmers     []CollectFarmers     `json:"farmers"`
	Attachments []CollectAttachments `json:"_attachments"`
	Points      []CollectPoint       `json:"points"`

	Notes       []any `json:"_tags"`
	Tags        []any `json:"_notes"`
	Geolocation []any `json:"_geolocation"`
}

func (c Collect) GetFormID() string { return "aE46T5DUSpYk6RCWu7uTiM" }
