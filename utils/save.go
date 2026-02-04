package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func SaveFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func SaveHttpResponse(resp http.Response) {

	contentType := resp.Header.Get("Content-Type")

	extension := "txt"

	switch contentType {
	case "application/json":
		extension = "json"
	case "application/xml":
		extension = "xml"
	case "text/html":
		extension = "html"
	case "text/plain":
		extension = "txt"
	}

	os.Mkdir("logs", 0755)
	filePath := fmt.Sprintf("logs/http-%d-%s.%s", resp.StatusCode, time.Now().Format("2006-01-02"), extension)
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "request url: %s\n", resp.Request.URL)
	fmt.Fprintf(os.Stderr, "response saved to %s\n", filePath)
	if resp.StatusCode >= 300 && resp.StatusCode < 200 {
		panic(fmt.Sprintf("http error: %d", resp.StatusCode))
	}

	fmt.Fprintln(os.Stderr, resp.Status)

}
