package main

import (
	"fmt"
	"go/format"
	"os"
	"strconv"
	"text/scanner"

	"crg.eti.br/go/config"
)

type Config struct {
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

const alert = "/* Automatically generated, do not change manually. */"

func main() {

	cfg := Config{}

	config.PrefixEnv = "BMV"
	err := config.Parse(&cfg)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if cfg.InputFile == "" {
		config.Usage()
		os.Exit(1)
	}

	err = parse(cfg.InputFile)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

}

func parse(fileName string) (err error) {

	var out string
	var iFile *os.File
	var hCount int
	out = alert + "\n\n"
	//var oFile *os.File

	iFile, err = os.Open(fileName)
	if err != nil {
		return
	}
	defer func() {
		if e := iFile.Close(); e != nil {
			err = e
			return
		}
	}()

	fi := fileInfo{}

	var s scanner.Scanner
	s.Init(iFile)
	s.Filename = fileName
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
			out += s.TokenText() + "\n"
		} else {

			switch ntok {
			case 0:
				fi.Version = s.TokenText()
			case 1:
				var h int
				h, err = strconv.Atoi(s.TokenText())
				if err != nil {
					return
				}
				fi.Height = h
			case 2:
				var w int
				w, err = strconv.Atoi(s.TokenText())
				if err != nil {
					return
				}
				fi.Width = w
			case 3:
				fi.PackageName = s.TokenText()
				out += "package " + fi.PackageName + "\n"
			case 4:
				fi.ObjectName = s.TokenText()

				out += fmt.Sprintf(`type %v struct {
Height int
Width int
Bitmap [][]byte
}
func (f *%v)Load(){
f.Height = %d
f.Width = %d
f.Bitmap = [][]byte{
`, fi.ObjectName, fi.ObjectName, fi.Height, fi.Width)
			default:
				//fmt.Println("At position", s.Pos(), ":", s.TokenText(), ntok)
				bs := s.TokenText()
				if bs == "" {
					continue
				}
				if hCount == 0 {
					out += "{\n"
				}

				if len(bs) > fi.Width {
					err = fmt.Errorf("Error at %v", s.Pos())
					println(err.Error())
					return
				}

				var r int64
				r, err = strconv.ParseInt(bs, 2, 64)
				if err != nil {
					println(err.Error())
					return
				}
				out += fmt.Sprintf("0x%02X, // %v\n", r, bs)

				hCount++
				if hCount >= fi.Height {
					hCount = 0
					out += "},\n"
				}
			}

			ntok++
		}

	}
	out += "}}\n"
	//println(out)

	arr := []byte(out)
	arr, err = format.Source(arr)
	if err != nil {
		return
	}
	println(string(arr))

	return
}
