package main

import (
	"bufio"
	"os"

	"github.com/crgimenes/goConfig"
)

type config struct {
	InputFile  string `json:"i" cfg:"i"`
	OutputFile string `json:"o" cfg:"o"`
}

func main() {

	cfg := config{}

	goConfig.PrefixEnv = "BMM"
	err := goConfig.Parse(&cfg)
	if err != nil {
		os.Exit(1)
	}

	if cfg.InputFile == "" {
		goConfig.Usage()
		os.Exit(1)
	}

	var iFile *os.File
	//var oFile *os.File

	iFile, err = os.Open(cfg.InputFile)
	if err != nil {
		os.Exit(1)
	}

	r := bufio.NewReader(iFile)

	var c rune
	for {
		c, _, err = r.ReadRune()
		if err != nil {
			os.Exit(1)
		}
		print(c)
	}

}
