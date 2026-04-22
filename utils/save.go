package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
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

func RenderHtmlTemplateAndPrintPDF(name, tpl string, data any) error {
	t := template.Must(template.New(name).Parse(tpl))

	// render HTML to buffer
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	// create temp file
	tmp, err := os.CreateTemp("", name+"-*.html")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // auto delete

	// write HTML
	if _, err := tmp.Write(buf.Bytes()); err != nil {
		return err
	}
	tmp.Close()

	// output PDF
	cmd := exec.Command(
		"brave-browser",
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--print-to-pdf="+name+".pdf",
		"--no-pdf-header-footer",
		"--enable-local-file-access",
		"file://"+tmp.Name(),
	)

	return cmd.Run()
}
