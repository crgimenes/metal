package main

import (
	"fmt"
	"os"
	"strconv"
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

	fi := fileInfo{}

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

	ntok := 0
	for tok != scanner.EOF {
		tok = s.Scan()
		if tok == scanner.Comment {
			println(s.TokenText())
		} else {

			switch ntok {
			case 0:
				fi.Version = s.TokenText()
			case 1:
				fi.PackageName = s.TokenText()
			case 2:
				fi.ObjectName = s.TokenText()
			case 3:
				var h int
				h, err = strconv.Atoi(s.TokenText())
				if err != nil {
					println(err.Error())
					os.Exit(1)
				}
				fi.Height = h
			case 4:
				var w int
				w, err = strconv.Atoi(s.TokenText())
				if err != nil {
					println(err.Error())
					os.Exit(1)
				}
				fi.Width = w
			default:
				fmt.Println("At position", s.Pos(), ":", s.TokenText(), ntok)
			}

			ntok++
		}
	}

}
