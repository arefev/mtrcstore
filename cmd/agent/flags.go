package main

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
)

type Config struct {
	Address string `env:"ADDRESS"`
	PollInterval int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
}

func ParseFlags() Config {
	cnf := Config{}

	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", "localhost:8080", "server address and port")
	f.IntVar(&cnf.PollInterval, "p", 2, "poll interval")
	f.IntVar(&cnf.ReportInterval, "r", 10, "report interval")
	f.Parse(os.Args[1:])

	// var cnf Config
	err := env.Parse(&cnf)
	if err != nil {
		log.Fatal(err)
	}

	return cnf
}
