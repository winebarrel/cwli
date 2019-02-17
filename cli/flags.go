package cli

import (
	"flag"
	"fmt"
	"os"
)

var version string

type Flags struct {
	Query   string
	Showptr bool
}

func ParseFlags() (flags *Flags) {
	flags = &Flags{}
	var printVersion bool

	flag.StringVar(&flags.Query, "query", "", "Query string")
	flag.BoolVar(&flags.Showptr, "showptr", false, "Show @ptr in query result")
	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.Parse()

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return
}
