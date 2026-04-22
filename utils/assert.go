package utils

import (
	"fmt"
	"os"
	"runtime"
)

func Assert(condition bool, msgs ...string) {
	if condition {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("Assertion failed at %s:%d\n", file, line)
	} else {
		fmt.Printf("Assertion failed at ?:?\n")
	}

	for _, msg := range msgs {
		fmt.Println(msg)
	}

	os.Exit(1)
}
