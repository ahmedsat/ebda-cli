package frappe

import (
	"github.com/ahmedsat/ebda-cli/config"
)

type LoginResult struct {
	Message  string `json:"message"`
	HomePage string `json:"home_page"`
	FullName string `json:"full_name"`
}

func Login() (result LoginResult, err error) {
	username := config.ErpUsername
	password := config.ErpPassword

	result, err = CallMethodT[LoginResult]("login", map[string]any{
		"usr": username,
		"pwd": password,
	})
	return
}
