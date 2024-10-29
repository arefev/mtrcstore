package main

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

const Address = "localhost:8080"

type Config struct {
	Address string `env:"ADDRESS"`
}

func NewConfig() (Config, error) {
	cnf := Config{}
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", Address, "address and port to run server")
	f.Parse(os.Args[1:])

	if err := env.Parse(&cnf); err != nil {
		return Config{}, err
	}

	return cnf, nil
}
