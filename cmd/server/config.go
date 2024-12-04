package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

const (
	Address         string = "localhost:8080"
	LogLevel        string = "info"
	DatabaseDSN     string = ""
	StoreInterval   int    = 300
	FileStoragePath string = ""
	Restore         bool   = true
)

type Config struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
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
	f.IntVar(&cnf.StoreInterval, "i", StoreInterval, "store interval")
	f.StringVar(&cnf.FileStoragePath, "f", FileStoragePath, "file storage path interval")
	f.StringVar(&cnf.DatabaseDSN, "d", DatabaseDSN, "db connection string")
	f.BoolVar(&cnf.Restore, "r", Restore, "need restore")
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
