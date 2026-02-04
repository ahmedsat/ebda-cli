package erp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ahmedsat/erp-reports-cli/utils"
)

func CallMethod(method string, args map[string]any) (result []byte, err error) {

	url_, err := url.JoinPath(os.Getenv("ERP_BASE_URL"), "/api/method/", method)
	if err != nil {
		return
	}

	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url_, strings.NewReader(string(jsonArgs)))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := utils.DoRequest(req)
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
