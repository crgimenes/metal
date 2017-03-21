package main

import (
	"fmt"
	"os"
	"text/scanner"

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

	var s scanner.Scanner
	s.Init(iFile)
	s.Filename = cfg.InputFile
	s.Mode = scanner.ScanIdents |
		scanner.ScanFloats |
		scanner.ScanChars |
		scanner.ScanStrings |
		scanner.ScanRawStrings |
		scanner.ScanComments
	var tok rune

	for tok != scanner.EOF {
		tok = s.Scan()
		if tok == scanner.Comment {
			println(s.TokenText())
		} else {
			fmt.Println("At position", s.Pos(), ":", s.TokenText())

		}
	}

}
