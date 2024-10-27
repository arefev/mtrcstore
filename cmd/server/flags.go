package main

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func ParseFlags() Config {
	cnf := Config{}
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", "localhost:8080", "address and port to run server")
	f.Parse(os.Args[1:])

	err := env.Parse(&cnf)
	if err != nil {
		log.Fatal(err)
	}

	return cnf
}
