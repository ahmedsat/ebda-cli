package main

import (
	"time"

	"github.com/ahmedsat/ebda-cli/utils"
)

func main() {
	for i := range 100 {
		n, err := utils.NewProgressNotification("prog", "myProg", "prog", i)
		if err != nil {
			panic(err)
		}
		_, _, err = n.Run()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

}
