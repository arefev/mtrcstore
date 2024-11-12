package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

const Address = "localhost:8080"
const LogLevel = "info"

type Config struct {
	Address string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
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
	f.StringVar(&cnf.Address, "a", Address, "address and port to run server")
	f.StringVar(&cnf.LogLevel, "l", LogLevel, "log level")
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
