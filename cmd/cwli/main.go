package main

import (
	"log"

	"github.com/winebarrel/cwli"
	"github.com/winebarrel/cwli/cli"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	flags := cli.ParseFlags()
	runner := cwli.NewRunner(flags)
	runner.Run()
}
