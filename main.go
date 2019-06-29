package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Command line argument
var (
	oConfigFile = ""
)

func init() {
	log.SetOutput(os.Stdout)
	flag.StringVar(&oConfigFile, "c", "cwc.toml", "Config file")
	flag.Parse()
}

func main() {
	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s password\n", os.Args[0])
		os.Exit(1)
	}
	password := flag.Args()[0]

	config, err := LoadConfig(oConfigFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	wallet := NewWallet(&config.Wallet)
	fmt.Printf("%s: %t\n", password, wallet.CheckPassword(password))
}
