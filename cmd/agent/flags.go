package main

import (
	"flag"
	"os"
)

var flagServerAddr string
var flagPollInterval int
var flagReportInterval int

func parseFlags() {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&flagServerAddr, "a", "localhost:8080", "server address and port")
	f.IntVar(&flagPollInterval, "p", 2, "poll interval")
	f.IntVar(&flagReportInterval, "r", 10, "report interval")
	f.Parse(os.Args[1:])
}
