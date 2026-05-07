package utils

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed trans.json
var transBytes []byte
var TransMap map[string]string

func Trans(key string) string {
	_, ok := TransMap[key]
	if !ok {
		return key
	}
	return TransMap[key]
}

func init() {
	err := json.Unmarshal(transBytes, &TransMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "enable to load translation: %s\n", err)
		os.Exit(1)
	}
}
