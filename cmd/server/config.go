package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

const Address = "localhost:8080"

type Config struct {
	Address string `env:"ADDRESS"`
}

func NewConfig() Config {
	return Config{}
}

func (cnf *Config) InitFlags(params []string) error {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", Address, "address and port to run server")
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
