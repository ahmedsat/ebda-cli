package kobo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ahmedsat/ebda-cli/config"
)

type KoboAsset interface {
	GetFormID() string
}

const (
	AssetsPath = "/api/v2/assets/"
)

func DoRequest(req *http.Request) (resp *http.Response, err error) {
	token := config.KoboAuthToken
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

func GetAssets[T KoboAsset]() (result []T, err error) {
	start := 0
	res, err := GetAssetsExt[T](0, start)
	if err != nil {
		return
	}
	result = res.Results
	for res.Next != "" {
		start += 0
		res, err = GetAssetsExt[T](0, start)
		if err != nil {
			return
		}
		result = append(result, res.Results...)
	}
	return
}

func GetAssetsExt[T KoboAsset](limit int, start int) (result AssetsResponse[T], err error) {
	var t T
	url, err := url.Parse(config.KoboBaseURL)
	if err != nil {
		return
	}
	url = url.JoinPath(AssetsPath, t.GetFormID(), "data")

	query := url.Query()
	if limit != 0 {
		query.Set("limit", fmt.Sprint(limit))
	}

	if start != 0 {
		query.Set("start", fmt.Sprint(start))
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

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	return
}
