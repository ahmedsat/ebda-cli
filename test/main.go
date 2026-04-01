package main

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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

func main() {

	_, err := frappe.Login()
	handelError(err)

	f, err := os.Open("workers.tsv")
	handelError(err)

	reader := csv.NewReader(f)
	reader.Comma = '\t'
	_, err = reader.Read() // consume the header
	handelError(err)

	lines, err := reader.ReadAll()
	handelError(err)

	runner := utils.NewSyncRunner(25, 100)
	ios := utils.SyncIoWriter{
		Writer: os.Stdout,
	}

	for _, line := range lines {

		runner.Run(func() error {

			fmt.Fprintln(os.Stderr, line[0])

			if len(line) != 14 {
				panic(fmt.Errorf("farm (%s) has %d lines, we expect 14", line[0], len(line)))
			}

			farm, err := frappe.Get1[types.Farm](line[1])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			Males1624, err := strconv.Atoi(line[8])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			Females1624, err := strconv.Atoi(line[9])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			Males2540, err := strconv.Atoi(line[10])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			Females2540, err := strconv.Atoi(line[11])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			Males4060, err := strconv.Atoi(line[12])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			Females4060, err := strconv.Atoi(line[13])
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			workers := []types.FarmWorker{
				{
					Count:  Males1624,
					Age:    "16-24",
					Gender: "Male",
				},
				{
					Count:  Females1624,
					Age:    "16-24",
					Gender: "Female",
				},
				{
					Count:  Males2540,
					Age:    "25-40",
					Gender: "Male",
				},
				{
					Count:  Females2540,
					Age:    "25-40",
					Gender: "Female",
				},
				{
					Count:  Males4060,
					Age:    "40-60",
					Gender: "Male",
				},
				{
					Count:  Females4060,
					Age:    "40-60",
					Gender: "Female",
				},
			}

			farm.Workers = workers
			err = farm.Update()
			if err != nil {
				fmt.Fprintf(&ios, "%s\t\"%s\"\n", line[0], err)
				return nil
			}

			return nil
		})

	}

}
