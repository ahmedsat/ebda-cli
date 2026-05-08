package kobo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/utils"
)

type KoboResponseCommon struct {
	Detail     string `json:"detail"`
	Status     string `json:"-"`
	StatusCode int    `json:"-"`
}

type KoboAsset interface {
	GetFormID() string
}

const (
	AssetsPath = "/api/v2/assets/"

	UpdateValidationStatePath = "/api/v2/assets/%s/data/%d/validation_status/" // (string,int)
	GetEditPath               = "/api/v2/assets/%s/data/%d/enketo/edit/"       // (string,int)
	ValidationStatusKey       = "_validation_status.uid"
)

func DoRequest(req *http.Request) (resp *http.Response, err error) {
	token := config.KoboAuthToken
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	return
}

func Get(path string) (resp *http.Response, err error) {
	url, err := url.JoinPath(config.KoboBaseURL, path)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return DoRequest(req)
}

type AssetsResponse[T KoboAsset] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

func GetAssets[T KoboAsset](query Query) (result []T, err error) {
	start := 0
	res, err := GetAssetsExt[T](query, 0, start)
	if err != nil {
		return
	}
	result = res.Results
	for res.Next != "" {
		start += len(res.Results)
		res, err = GetAssetsExt[T](query, 0, start)
		if err != nil {
			return
		}
		result = append(result, res.Results...)
	}
	return
}

func StreamAssets[T KoboAsset](query Query) (int, <-chan T, <-chan error) {
	out := make(chan T)
	errCh := make(chan error, 1) // buffered to avoid goroutine leak

	start := 0
	res, err := GetAssetsExt[T](query, 0, start)
	if err != nil {
		defer close(out)
		defer close(errCh)
		errCh <- err
		return 0, out, errCh
	}

	go func() {
		defer close(out)
		defer close(errCh)
		for {
			for _, item := range res.Results {
				out <- item
			}

			if res.Next == "" {
				return
			}

			start += len(res.Results)

			var err error
			res, err = GetAssetsExt[T](query, 0, start)
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	return res.Count, out, errCh
}

func GetAssetByID[T KoboAsset](id int) (result T, err error) {
	var t T
	url, err := url.Parse(config.KoboBaseURL)
	if err != nil {
		return
	}
	url = url.JoinPath(AssetsPath, t.GetFormID(), "data", fmt.Sprint(id))

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}
	resp, err := DoRequest(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("http error: %d", resp.StatusCode)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	// decoder.DisallowUnknownFields()

	err = decoder.Decode(&result)
	if err != nil {
		return
	}

	return
}

func GetAssetsExt[T KoboAsset](q Query, limit int, start int) (result AssetsResponse[T], err error) {

	// https://community.kobotoolbox.org/t/important-changes-to-api-v2-assets-uid-asset-data-result-limits/74610
	if limit > 1000 {
		limit = 1000
	}

	var t T
	url, err := url.Parse(config.KoboBaseURL)
	if err != nil {
		return
	}
	url = url.JoinPath(AssetsPath, t.GetFormID(), "data")

	query := url.Query()

	query.Set("format", "json")

	if limit != 0 {
		query.Set("limit", fmt.Sprint(limit))
	}

	if start != 0 {
		query.Set("start", fmt.Sprint(start))
	}

	if len(q) != 0 {
		query.Set("query", q.String())
		url.RawQuery = query.Encode()
	}

	url.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}
	resp, err := DoRequest(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("http error: %d", resp.StatusCode)
		return
	}

	decoder := json.NewDecoder(resp.Body)

	// utils.SaveHttpResponse(*resp)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&result)
	if err != nil {
		return
	}

	return
}

type ValidationSate string

const (
	ValidationStatusApproved    ValidationSate = "validation_status_approved"
	ValidationStatusNotApproved ValidationSate = "validation_status_not_approved"
	ValidationStatusOnHold      ValidationSate = "validation_status_on_hold"
)

type UpdateValidationStateResponse struct {
	KoboResponseCommon
	Timestamp int    `json:"timestamp"`
	Uid       string `json:"uid"`
	ByWhom    string `json:"by_whom"`
	Label     string `json:"label"`
}

func UpdateValidationState[T KoboAsset](id int, state ValidationSate) (response UpdateValidationStateResponse, err error) {
	var t T
	url, err := url.Parse(config.KoboBaseURL)
	if err != nil {
		return
	}

	url = url.JoinPath(fmt.Sprintf(UpdateValidationStatePath, t.GetFormID(), id))

	req, err := http.NewRequest(
		"PATCH",
		url.String(),
		strings.NewReader(fmt.Sprintf("{\"validation_status.uid\":\"%s\"}", state)),
	)
	if err != nil {
		return
	}

	resp, err := DoRequest(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		err = fmt.Errorf("http error: %d", resp.StatusCode)
		utils.SaveHttpResponse(*resp)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return
	}

	response.Status = resp.Status
	response.StatusCode = resp.StatusCode

	return
}

type GetUpdateURLResponse struct {
	KoboResponseCommon
	Url     string `json:"url"`
	Version string `json:"version"`
}

func GetUpdateURL[T KoboAsset](id int) (response GetUpdateURLResponse, err error) {
	var t T
	url, err := url.Parse(config.KoboBaseURL)
	if err != nil {
		return
	}
	url = url.JoinPath(fmt.Sprintf(GetEditPath, t.GetFormID(), id))

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}
	resp, err := DoRequest(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 &&
		resp.StatusCode != 404 &&
		resp.StatusCode != 403 {
		err = fmt.Errorf("http error: %d", resp.StatusCode)
		utils.SaveHttpResponse(*resp)
		return
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return
	}

	response.Status = resp.Status
	response.StatusCode = resp.StatusCode

	return

}

func Download(url string) (r io.ReadCloser, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := DoRequest(req)
	if err != nil {
		return
	}
	return resp.Body, nil
}

func DownloadAttach(att Attachment) (localPath string, err error) {

	localPath = att.Filename
	// check if file exists
	if utils.FileExists(localPath) {
		return
	}

	// local dir
	dir := filepath.Dir(localPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	f, err := os.Create(localPath)

	req, err := http.NewRequest("GET", att.DownloadUrl, nil)
	if err != nil {
		return
	}
	resp, err := DoRequest(req)
	if err != nil {
		return
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return
	}
	return

}
