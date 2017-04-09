package cmd

import (
	"fmt"
	"os"
)

func Eval(c string) {

	if c == "" {
		return
	}
	fmt.Printf("%s\n", c)

	if c == `EXIT` {
		os.Exit(0)
	}
}
