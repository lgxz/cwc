package main

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Wallet WalletOptions
}

type WalletOptions struct {
	Coin             string
	Salt             string
	DeriveIterations uint32
	Crypted_key      string
	Ckey             string
	Pubkey           string
}

func LoadConfig(configFile string) (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("Config file does not exist.")
	} else if err != nil {
		return nil, err
	}

	var conf Config
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
