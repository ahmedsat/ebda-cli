package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/sheets"
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
	// 11tXfIz9o_cgD-czMTQRRLF9JEkmvNTF5QSmdY6lVQFs
	err := sheets.WriteRange(
		context.Background(),
		"11tXfIz9o_cgD-czMTQRRLF9JEkmvNTF5QSmdY6lVQFs",
		"sheet1!A1",
		[][]any{
			{"hello", "world"},
		},
	)
	fmt.Println(err)
}
