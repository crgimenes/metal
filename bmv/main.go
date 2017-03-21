package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/crgimenes/goConfig"
)

type config struct {
	InputFile  string `json:"i" cfg:"i"`
	OutputFile string `json:"o" cfg:"o"`
}

type fileInfo struct {
	Version     string
	PackageName string
	ObjectName  string
	Height      int
	Width       int
}

func main() {

	cfg := config{}

	goConfig.PrefixEnv = "BMV"
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

	s := bufio.NewScanner(iFile)

	finfo := fileInfo{}
	line := 0
	var l string
	for s.Scan() {
		l = strings.TrimSpace(s.Text())

		if l == "" {
			continue
		}

		switch line {
		case 0:
			finfo.Version = l
		case 1:
			finfo.PackageName = l
		case 2:
			finfo.ObjectName = l
		case 3:
			p := strings.Split(l, " ")
			if len(p) != 2 {
				println("File format error, line 3, expected height and width")
				os.Exit(1)
			}
			finfo.Height, err = strconv.Atoi(p[0])
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
			finfo.Width, err = strconv.Atoi(p[1])
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
		default:
			println(l)
		}

		line++
	}
}
