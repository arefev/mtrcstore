package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

const (
	address         string = "localhost:8080"
	logLevel        string = "info"
	databaseDSN     string = ""
	fileStoragePath string = ""
	secretKey       string = ""
	cryptoKey       string = ""
	configPath      string = ""
	trustedSubnet   string = ""
	grpcAddress     string = ""
	storeInterval   int    = 300
	restore         bool   = true
)

type Config struct {
	Address         string `env:"ADDRESS" json:"address"`
	LogLevel        string `env:"LOG_LEVEL" json:"log_level"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	SecretKey       string `env:"KEY" json:"secret_key"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath      string `env:"CONFIG" json:"-"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GRPCAddress     string `env:"GRPC_ADDRESSS" json:"grpc_address"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	Restore         bool   `env:"RESTORE" json:"restore"`
}

func NewConfig(params []string) (Config, error) {
	cnf := Config{
		Address:         address,
		LogLevel:        logLevel,
		DatabaseDSN:     databaseDSN,
		FileStoragePath: fileStoragePath,
		SecretKey:       secretKey,
		CryptoKey:       cryptoKey,
		ConfigPath:      configPath,
		StoreInterval:   storeInterval,
		Restore:         restore,
		TrustedSubnet:   trustedSubnet,
		GRPCAddress:     grpcAddress,
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
	f.StringVar(&cnf.Address, "a", cnf.Address, "address and port to run server")
	f.StringVar(&cnf.LogLevel, "l", cnf.LogLevel, "log level")
	f.StringVar(&cnf.FileStoragePath, "f", cnf.FileStoragePath, "file storage path interval")
	f.StringVar(&cnf.DatabaseDSN, "d", cnf.DatabaseDSN, "db connection string")
	f.StringVar(&cnf.SecretKey, "k", cnf.SecretKey, "secret key")
	f.StringVar(&cnf.CryptoKey, "crypto-key", cnf.CryptoKey, "path to file with private key")
	f.StringVar(&cnf.ConfigPath, "c", cnf.ConfigPath, "path to file with config")
	f.StringVar(&cnf.ConfigPath, "config", cnf.ConfigPath, "path to file with config")
	f.StringVar(&cnf.TrustedSubnet, "t", cnf.TrustedSubnet, "CIDR")
	f.StringVar(&cnf.GRPCAddress, "grpc-addr", cnf.GRPCAddress, "GRPC address")
	f.IntVar(&cnf.StoreInterval, "i", cnf.StoreInterval, "store interval")
	f.BoolVar(&cnf.Restore, "r", cnf.Restore, "need restore")
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
