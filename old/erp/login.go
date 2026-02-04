package erp

import (
	"encoding/json"
	"errors"
	"os"
)

func Login() (result string, err error) {
	username := os.Getenv("ERP_USERNAME")
	if username == "" {
		return "", errors.New("ERP_USERNAME is not set")
	}
	password := os.Getenv("ERP_PASSWORD")
	if password == "" {
		return "", errors.New("ERP_PASSWORD is not set")
	}

	resultBytes, err := CallMethod("login", map[string]any{
		"usr": username,
		"pwd": password,
	})
	if err != nil {
		return
	}

	if resultBytes == nil {
		return "", errors.New("login failed")
	}

	type Result struct {
		FullName string `json:"full_name"`
	}
	var r Result
	json.Unmarshal(resultBytes, &r)

	result = r.FullName
	return
}
