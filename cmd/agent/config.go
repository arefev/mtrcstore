package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

const (
	Address        = "localhost:8080"
	SecretKey      = ""
	PollInterval   = 2
	ReportInterval = 10
	RateLimit      = 3
)

type Config struct {
	Address        string `env:"ADDRESS"`
	SecretKey      string `env:"KEY"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewConfig(params []string) (Config, error) {
	cnf := Config{}
	if err := cnf.initFlags(params); err != nil {
		return Config{}, err
	}

	if err := cnf.initEnvs(); err != nil {
		return Config{}, err
	}

	return cnf, nil
}

func (cnf *Config) initFlags(params []string) error {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", Address, "server address and port")
	f.StringVar(&cnf.SecretKey, "k", SecretKey, "secret key")
	f.IntVar(&cnf.PollInterval, "p", PollInterval, "poll interval")
	f.IntVar(&cnf.ReportInterval, "r", ReportInterval, "report interval")
	f.IntVar(&cnf.RateLimit, "l", RateLimit, "rate limit")
	if err := f.Parse(params); err != nil {
		return fmt.Errorf("InitFlags: parse flags fail: %w", err)
	}

	return nil
}

func (cnf *Config) initEnvs() error {
	if err := env.Parse(cnf); err != nil {
		return fmt.Errorf("InitEnvs: parse envs fail: %w", err)
	}

	return nil
}
