package frappe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/utils"
)

type FrappeDoctype interface {
	DocTypeName() string
	DocName() string
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

	if id == "" {
		err = errors.New("id is required")
		return
	}

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
	// utils.SaveHttpResponse(*resp)

	decoder := json.NewDecoder(resp.Body)

	// decoder.DisallowUnknownFields()

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

func Delete[T FrappeDoctype](id string) (err error) {

	var DeleteResponse = struct {
		Data string `json:"data"`
	}{}

	var t T
	url, err := url.JoinPath(config.ErpBaseUrl, "/api/resource", t.DocTypeName(), id)
	if err != nil {
		return
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&DeleteResponse)
	if err != nil {
		return
	}

	if resp.StatusCode != 202 {
		utils.SaveHttpResponse(*resp)
		err = errors.Join(fmt.Errorf("http error: %d", resp.StatusCode), errors.New(DeleteResponse.Data))
		return
	}

	return
}

var ErrDuplicated = errors.New("duplicated")

func Create[T FrappeDoctype](data T) (result T, err error) {

	url, err := url.JoinPath(config.ErpBaseUrl, "/api/resource", data.DocTypeName())
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(jsonData))

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		utils.SaveHttpResponse(*resp)
		err = fmt.Errorf("http error: %d", resp.StatusCode)
		if resp.StatusCode == 409 {
			err = ErrDuplicated
		}
		return
	}

	decoder := json.NewDecoder(resp.Body)

	var response = struct {
		Data T `json:"data"`
	}{
		Data: data,
	}

	err = decoder.Decode(&response)
	if err != nil {
		return
	}

	result = response.Data

	return
}

func UpdateDoc[T FrappeDoctype](doc T) (result T, err error) {

	if doc.DocName() == "" {
		err = errors.New("doc name is required")
		return
	}

	url, err := url.JoinPath(config.ErpBaseUrl, "/api/resource", doc.DocTypeName(), doc.DocName())
	if err != nil {
		return
	}
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	jsonData, err := json.Marshal(doc)
	if err != nil {
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(jsonData))

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

	data := struct {
		Data T `json:"data"`
	}{
		Data: doc,
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}

	result = data.Data

	return
}
