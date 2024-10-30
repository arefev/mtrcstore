package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

const (
	Address = "localhost:8080"
	PollInterval = 2
	ReportInterval = 10
)

type Config struct {
	Address string `env:"ADDRESS"`
	PollInterval int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
}

func NewConfig() (Config, error) {
	cnf := Config{}

	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", Address, "server address and port")
	f.IntVar(&cnf.PollInterval, "p", PollInterval, "poll interval")
	f.IntVar(&cnf.ReportInterval, "r", ReportInterval, "report interval")
	f.Parse(os.Args[1:])

	if err := env.Parse(cnf); err != nil {
		return Config{}, fmt.Errorf("NewConfig: parse envs fail: %w", err)
	}

	return cnf, nil
}
