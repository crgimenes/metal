package cmd

import (
	"fmt"
	"os"
)

func Eval(c string) {
	fmt.Printf("[%s]\n", c)

	if c == `EXIT` {
		os.Exit(0)
	}
}
