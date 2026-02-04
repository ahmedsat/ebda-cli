package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(data any) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}
