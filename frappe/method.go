package frappe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/utils"
)

func CallMethodT[T any](method string, args map[string]any) (result T, err error) {
	resultBytes, err := CallMethod(method, args)
	if err != nil {
		return
	}
	err = json.Unmarshal(resultBytes, &result)
	return
}

func CallMethod(method string, args map[string]any) (result []byte, err error) {
	url, err := url.JoinPath(config.ErpBaseUrl, "/api/method/", method)
	if err != nil {
		return
	}

	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonArgs)))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		utils.SaveHttpResponse(*resp)
		err = errors.Join(fmt.Errorf("http error: %d", resp.StatusCode), errors.New("failed to get response"))
		return
	}

	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

type ProgressReader struct {
	r        io.Reader
	total    int64
	read     int64
	onUpdate func(read, total int64)
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	if n > 0 {
		p.read += int64(n)
		if p.onUpdate != nil {
			p.onUpdate(p.read, p.total)
		}
	}
	return n, err
}

type FileUploadResponse struct {
	Name                  string `json:"name"`
	Owner                 string `json:"owner"`
	Creation              string `json:"creation"`
	Modified              string `json:"modified"`
	ModifiedBy            string `json:"modified_by"`
	Docstatus             int    `json:"docstatus"`
	Idx                   int    `json:"idx"`
	FileName              string `json:"file_name"`
	IsPrivate             int    `json:"is_private"`
	FileType              string `json:"file_type"`
	IsHomeFolder          int    `json:"is_home_folder"`
	IsAttachmentsFolder   int    `json:"is_attachments_folder"`
	FileSize              int    `json:"file_size"`
	FileUrl               string `json:"file_url"`
	Folder                string `json:"folder"`
	IsFolder              int    `json:"is_folder"`
	ContentHash           string `json:"content_hash"`
	UploadedToDropbox     int    `json:"uploaded_to_dropbox"`
	UploadedToGoogleDrive int    `json:"uploaded_to_google_drive"`
	Doctype               string `json:"doctype"`
}

func UploadFile(
	data []byte,
	fileName string,
	onProgress func(sent, total int64),
) (res FileUploadResponse, err error) {

	url_, err := url.JoinPath(config.ErpBaseUrl, "/api/method/upload_file")
	if err != nil {
		return
	}

	// Buffer for multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// File field (Frappe expects "file")
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return
	}

	if _, err = part.Write(data); err != nil {
		return
	}

	// Optional: explicit file name field (Frappe supports it)
	if err = writer.WriteField("file_name", fileName); err != nil {
		return
	}

	if err = writer.Close(); err != nil {
		return
	}

	// Wrap body with progress reader
	pr := &ProgressReader{
		r:     body,
		total: int64(body.Len()),
		onUpdate: func(read, total int64) {
			onProgress(read, total)
		},
	}

	req, err := http.NewRequest("POST", url_, pr)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.ContentLength = int64(body.Len())

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("upload failed: %d", resp.StatusCode)
		return
	}

	// utils.SaveHttpResponse(*resp)
	// Decode response
	resStruct := struct {
		Message FileUploadResponse `json:"message"`
	}{}

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&resStruct)
	if err != nil {
		return
	}

	res = resStruct.Message

	return
}

func TestUrl(url_ string) (resp *http.Response, err error) {

	url, err := url.JoinPath(config.ErpBaseUrl, url_)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err = client.Do(req)

	if err != nil {
		return
	}

	return
}
