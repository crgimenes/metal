package main

import (
	"github.com/crgimenes/goConfig"
)

type config struct {
	InputFile  string `json:"i" cfg:"i"`
	OutputFile string `json:"o" cfg:"o"`
}

func main() {

	cfg := config{}

	goConfig.PrefixEnv = "BMM"
	goConfig.Parse(&cfg)

}
