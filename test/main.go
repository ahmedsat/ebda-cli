package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/utils"
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

var f = "auto-sheet"

func main() {
	switch f {
	case "auto-sheet":
		autoSheet()
	case "follow-up":
		followUp()
	default:
		panic("unknown command")
	}
}

func autoSheet() {
	// // 11tXfIz9o_cgD-czMTQRRLF9JEkmvNTF5QSmdY6lVQFs
	// err := sheets.Append(
	// 	context.Background(),
	// 	"11tXfIz9o_cgD-czMTQRRLF9JEkmvNTF5QSmdY6lVQFs",
	// 	"sheet1!A1",
	// 	[][]any{
	// 		{"hello", "world"},
	// 	},
	// )
	// handelError(err)
}

func followUp() {
	maps, err := frappe.Get[types.MapRecord](nil, nil, nil)
	handelError(err)

	fmt.Println(len(maps))
	maps = utils.Filter(maps, func(m types.MapRecord) bool { return m.Color == "#181818" })
	fmt.Println(len(maps))

	bytes, err := types.RecordsToKML(maps)

	err = os.WriteFile("maps.kml", bytes, 0666)
	handelError(err)
}
