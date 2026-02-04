package frappe

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/utils"
)

type FrappeDoctype interface {
	DocTypeName() string
}

var client *http.Client

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client = &http.Client{
		Jar: jar,
	}
}

func Get[T FrappeDoctype](filters Filters, fields List) (result []T, err error) {
	return GetEx[T](filters, fields, false)
}

func GetEx[T FrappeDoctype](filters Filters, fields List, restricted bool) (result []T, err error) {
	var t T
	url, err := url.JoinPath(config.ErpBaseUrl, "/api/resource", t.DocTypeName())
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("limit", "0")

	if len(filters) != 0 {
		q.Add("filters", filters.String())
	}

	if len(fields) == 0 {
		fields = List{"*"}
	}
	q.Add("fields", fields.String())

	req.URL.RawQuery = q.Encode()

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

	decoder := json.NewDecoder(resp.Body)
	if restricted {
		decoder.DisallowUnknownFields()
	}

	var response struct {
		Data []T `json:"data"`
	}

	err = decoder.Decode(&response)
	if err != nil {
		return
	}

	result = response.Data

	return
}

func Get1[T FrappeDoctype](id string) (result T, err error) {
	url, err := url.JoinPath(config.ErpBaseUrl, "/api/resource", result.DocTypeName(), id)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

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

	decoder := json.NewDecoder(resp.Body)

	var response = struct {
		Data T `json:"data"`
	}{
		Data: result,
	}

	err = decoder.Decode(&response)
	if err != nil {
		return
	}

	result = response.Data
	return
}
