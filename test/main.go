package main

import (
	_ "embed"
	"fmt"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/kobo"
)

func init() {

	err := config.Configure()
	if err != nil {
		panic(err)
	}
}

func handelError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	_, err := frappe.Login()
	handelError(err)

	r, err := kobo.GetAssetsExt[kobo.PGSNew](nil, 0, 1000)
	handelError(err)
	fmt.Println(len(r.Results))
}
