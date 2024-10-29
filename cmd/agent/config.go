package main

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
)

const (
	ADDRESS = "localhost:8080"
	POLL_INTERVAL = 2
	REPORT_INTERVAL = 10
)

type Config struct {
	Address string `env:"ADDRESS"`
	PollInterval int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
}

func NewConfig() Config {
	cnf := Config{}

	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", ADDRESS, "server address and port")
	f.IntVar(&cnf.PollInterval, "p", POLL_INTERVAL, "poll interval")
	f.IntVar(&cnf.ReportInterval, "r", REPORT_INTERVAL, "report interval")
	f.Parse(os.Args[1:])

	err := env.Parse(&cnf)
	if err != nil {
		log.Fatal(err)
	}

	return cnf
}
