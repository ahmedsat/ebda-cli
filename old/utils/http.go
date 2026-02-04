package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

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
	HandelErr(err)
	_, err = io.Copy(f, resp.Body)
	HandelErr(err)

	fmt.Fprintf(os.Stderr, "request url: %s\n", resp.Request.URL)
	fmt.Fprintf(os.Stderr, "response saved to %s\n", filePath)
	if resp.StatusCode >= 300 && resp.StatusCode < 200 {
		HandelErr(fmt.Errorf("http error: %d", resp.StatusCode))
	}

	fmt.Fprintln(os.Stderr, resp.Status)

}

func DoRequest(req *http.Request) (resp *http.Response, err error) {
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	return
}

func Get(path string) (resp *http.Response, err error) {

	url, err := url.JoinPath(os.Getenv("ERP_BASE_URL"), path)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return DoRequest(req)
}
