package kobo

import (
	"fmt"

	"github.com/ahmedsat/ebda-cli/utils"
)

type Attachment struct {
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

type Common struct {
	ID               int64        `json:"_id"`
	Start            string       `json:"start"`
	End              string       `json:"end"`
	StartGeopoint    string       `json:"start-geopoint"`
	Today            string       `json:"today"`
	Username         string       `json:"username"`
	Deviceid         string       `json:"deviceid"`
	Phonenumber      string       `json:"phonenumber"`
	FormhubUuid      string       `json:"formhub/uuid"`
	Version          string       `json:"__version__"`
	MetaInstanceID   string       `json:"meta/instanceID"`
	XformIdString    string       `json:"_xform_id_string"`
	Uuid             string       `json:"_uuid"`
	MetaRootUuid     string       `json:"meta/rootUuid"`
	Attachments      []Attachment `json:"_attachments"`
	Status           string       `json:"_status"`
	Geolocation      []float64    `json:"_geolocation"`
	SubmissionTime   string       `json:"_submission_time"`
	SubmittedBy      string       `json:"_submitted_by"`
	Tags             []any        `json:"_tags"`
	Notes            []any        `json:"_notes"`
	ValidationStatus `json:"_validation_status"`
	AttachmentMap    map[string]Attachment `json:"-"`
}

func (c Common) FillAttachmentMap() {
	c.AttachmentMap = make(map[string]Attachment)
	for i := range c.Attachments {
		_, ok := c.AttachmentMap[c.Attachments[i].MediaFileBasename]
		utils.Assert(!ok, fmt.Sprintf("duplicate attachment: %s in %d", c.Attachments[i].MediaFileBasename, c.ID))
		c.AttachmentMap[c.Attachments[i].MediaFileBasename] = c.Attachments[i]
	}
}

type AssetsList struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

func (AssetsList) GetFormID() string { panic("assets list is not a typical asset") }
