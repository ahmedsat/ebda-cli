package frappe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/utils"
)

type LoginResult struct {
	Message  string `json:"message"`
	HomePage string `json:"home_page"`
	FullName string `json:"full_name"`
}

func Login() (result LoginResult, err error) {
	authMu.Lock()
	defer authMu.Unlock()

	result, err = loginLocked()
	if err != nil {
		return result, err
	}

	authenticated = true
	return
}

func loginLocked() (result LoginResult, err error) {
	url, err := url.JoinPath(config.ErpBaseUrl, "/api/method/login")
	if err != nil {
		return result, err
	}

	payload, err := json.Marshal(map[string]any{
		"usr": config.ErpUsername,
		"pwd": config.ErpPassword,
	})
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := rawDo(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.SaveHttpResponse(*resp)
		return result, errors.Join(
			fmt.Errorf("http error: %d", resp.StatusCode),
			errors.New("failed to login"),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
