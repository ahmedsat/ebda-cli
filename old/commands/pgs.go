package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Pgs(args []string) (err error) {
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, "PGS")
		}
	}()
	fmt.Println("PGS")

	formUrl := "https://kf.kobotoolbox.org/api/v2/assets/aX4NJWgge6tooXjfSYXhrq" + args[0]
	req, err := http.NewRequest("GET", formUrl, nil)
	if err != nil {
		return
	}

	token := os.Getenv("KOBO_AUTH_TOKEN")
	req.Header.Set("Authorization", "Token "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("http error: %d", res.StatusCode)
		return
	}

	io.Copy(os.Stderr, res.Body)

	return
}
