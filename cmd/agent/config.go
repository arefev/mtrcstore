package main

import (
	"flag"
	"fmt"

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

func NewConfig() Config {
	return Config{}
}

func (cnf *Config) InitFlags(params []string) error {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", Address, "server address and port")
	f.IntVar(&cnf.PollInterval, "p", PollInterval, "poll interval")
	f.IntVar(&cnf.ReportInterval, "r", ReportInterval, "report interval")
	if err := f.Parse(params); err != nil {
		return fmt.Errorf("InitFlags: parse flags fail: %w", err)
	}

	return nil
}

func (cnf *Config) InitEnvs() error {
	if err := env.Parse(cnf); err != nil {
		return fmt.Errorf("InitEnvs: parse envs fail: %w", err)
	}

	return nil
}


