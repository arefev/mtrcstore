package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

const (
	address         string = "localhost:8080"
	logLevel        string = "info"
	databaseDSN     string = ""
	fileStoragePath string = ""
	secretKey       string = ""
	cryptoKey       string = ""
	storeInterval   int    = 300
	restore         bool   = true
)

type Config struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	SecretKey       string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
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
	f.StringVar(&cnf.Address, "a", address, "address and port to run server")
	f.StringVar(&cnf.LogLevel, "l", logLevel, "log level")
	f.StringVar(&cnf.FileStoragePath, "f", fileStoragePath, "file storage path interval")
	f.StringVar(&cnf.DatabaseDSN, "d", databaseDSN, "db connection string")
	f.StringVar(&cnf.SecretKey, "k", secretKey, "secret key")
	f.StringVar(&cnf.CryptoKey, "crypto-key", cryptoKey, "path to file with private key")
	f.IntVar(&cnf.StoreInterval, "i", storeInterval, "store interval")
	f.BoolVar(&cnf.Restore, "r", restore, "need restore")
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
