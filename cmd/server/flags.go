package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlags() {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	f.Parse(os.Args[1:])
}
