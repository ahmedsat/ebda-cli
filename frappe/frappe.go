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
	"strings"
	"sync"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/utils"
)

type FrappeDoctype interface {
	DocTypeName() string
	DocName() string
}

var client *http.Client
var authMu sync.Mutex
var authenticated bool

func Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("nil request")
	}

	err := makeRequestReplayable(req)
	if err != nil {
		return nil, err
	}

	if isLoginRequest(req) {
		return rawDo(req)
	}

	err = ensureAuthenticated()
	if err != nil {
		return nil, err
	}

	resp, err := rawDo(req)
	if err != nil {
		return nil, err
	}

	if !shouldRetryAuthentication(resp.StatusCode) {
		return resp, nil
	}

	resp.Body.Close()
	resetAuthentication()

	err = ensureAuthenticated()
	if err != nil {
		return nil, err
	}

	retryReq, err := cloneRequest(req)
	if err != nil {
		return nil, err
	}

	return rawDo(retryReq)
}

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client = &http.Client{
		Jar: jar,
	}
}

func rawDo(req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

func ensureAuthenticated() error {
	authMu.Lock()
	defer authMu.Unlock()

	if authenticated {
		return nil
	}

	_, err := loginLocked()
	if err != nil {
		return err
	}

	authenticated = true
	return nil
}

func resetAuthentication() {
	authMu.Lock()
	defer authMu.Unlock()
	authenticated = false
}

func isLoginRequest(req *http.Request) bool {
	return strings.TrimSuffix(req.URL.Path, "/") == "/api/method/login"
}

func shouldRetryAuthentication(statusCode int) bool {
	return statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden
}

func makeRequestReplayable(req *http.Request) error {
	if req.Body == nil || req.GetBody != nil {
		return nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	req.ContentLength = int64(len(body))

	return nil
}

func cloneRequest(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())
	if req.GetBody == nil {
		return clone, nil
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	clone.Body = body
	return clone, nil
}

func Get[T FrappeDoctype](filters Filters, fields, expand List) (result []T, err error) {
	return GetEx[T](filters, fields, expand, false)
}

func GetEx[T FrappeDoctype](filters Filters, fields, expand List, restricted bool) (result []T, err error) {
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

	if len(expand) != 0 {
		q.Add("expand", expand.String())
	}

	req.URL.RawQuery = q.Encode()

	resp, err := Do(req)
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

	resp, err := Do(req)
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

	resp, err := Do(req)
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

	resp, err := Do(req)
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

	resp, err := Do(req)
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
