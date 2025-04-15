package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

const (
	address        = "localhost:8080"
	secretKey      = ""
	cryptoKey      = ""
	configPath     = ""
	grpcAddress    = ""
	pollInterval   = 2
	reportInterval = 10
	rateLimit      = 3
)

type Config struct {
	Address        string `env:"ADDRESS" json:"address"`
	SecretKey      string `env:"KEY" json:"secret_key"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath     string `env:"CONFIG" json:"-"`
	GRPCAddress    string `env:"GRPC_ADDRESSS" json:"grpc_address"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	RateLimit      int    `env:"RATE_LIMIT" json:"rate_limit"`
}

func NewConfig(params []string) (Config, error) {
	cnf := Config{
		Address:        address,
		SecretKey:      secretKey,
		CryptoKey:      cryptoKey,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		RateLimit:      rateLimit,
		GRPCAddress:    grpcAddress,
	}

	if err := cnf.initConfig(params); err != nil {
		return Config{}, err
	}

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
	f.StringVar(&cnf.Address, "a", cnf.Address, "server address and port")
	f.StringVar(&cnf.SecretKey, "k", cnf.SecretKey, "secret key")
	f.StringVar(&cnf.CryptoKey, "crypto-key", cnf.CryptoKey, "path to file with public key")
	f.StringVar(&cnf.ConfigPath, "c", cnf.ConfigPath, "path to file with config")
	f.StringVar(&cnf.ConfigPath, "config", cnf.ConfigPath, "path to file with config")
	f.StringVar(&cnf.GRPCAddress, "grpc-addr", cnf.GRPCAddress, "GRPC address")
	f.IntVar(&cnf.PollInterval, "p", cnf.PollInterval, "poll interval")
	f.IntVar(&cnf.ReportInterval, "r", cnf.ReportInterval, "report interval")
	f.IntVar(&cnf.RateLimit, "l", cnf.RateLimit, "rate limit")
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

func (cnf *Config) initConfig(params []string) error {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		if err := cnf.initFlags(params); err != nil {
			return fmt.Errorf("findConfig: initFlags fail: %w", err)
		}

		configPath = cnf.ConfigPath
	}

	if configPath == "" {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("findConfig: read config file fail: %w", err)
	}

	err = json.Unmarshal(data, &cnf)
	if err != nil {
		return fmt.Errorf("findConfig: json unmarshal fail: %w", err)
	}

	return nil
}
