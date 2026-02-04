package frappe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ahmedsat/ebda-cli/config"
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
		err = errors.Join(fmt.Errorf("http error: %d", resp.StatusCode), errors.New("failed to get response"))
		return
	}

	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
