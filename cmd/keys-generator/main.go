package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"log"
	"os"
)

const (
	bits           int         = 4096
	filePermission fs.FileMode = 0o644
	privateKeyName string      = "private.pem"
	publicKeyName  string      = "public.pem"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("key generator run - GenerateKey failed: %w", err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile(privateKeyName, privateKeyPEM, filePermission)
	if err != nil {
		return fmt.Errorf("key generator run - WriteFile with private key failed: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("key generator run - MarshalPKIXPublicKey failed: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile(publicKeyName, publicKeyPEM, filePermission)
	if err != nil {
		return fmt.Errorf("key generator run - WriteFile with public key failed: %w", err)
	}

	return nil
}
