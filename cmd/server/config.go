package main

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
)

const ADDRESS = "localhost:8080"

type Config struct {
	Address string `env:"ADDRESS"`
}

func NewConfig() Config {
	cnf := Config{}
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", ADDRESS, "address and port to run server")
	f.Parse(os.Args[1:])

	err := env.Parse(&cnf)
	if err != nil {
		log.Fatal(err)
	}

	return cnf
}
