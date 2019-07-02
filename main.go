package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// Command line argument
var (
	oConfigFile = ""
	oMaxCPU     = 0
)

func init() {
	log.SetOutput(os.Stdout)
	flag.StringVar(&oConfigFile, "c", "cwc.toml", "Config file")
	flag.IntVar(&oMaxCPU, "p", 0, "Max CPU")
	flag.Parse()

	if oMaxCPU <= 0 {
		oMaxCPU = runtime.NumCPU()
	}
	fmt.Printf("Use %d CPU\n", oMaxCPU)
}

func main() {
	config, err := LoadConfig(oConfigFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	wallet := NewWallet(&config.Wallet)

	if flag.NArg() > 0 {
		for _, password := range flag.Args() {
			if wallet.CheckPassword(password) {
				fmt.Fprintf(os.Stderr, "OK: %s\n", password)
				break
			}
		}
		os.Exit(0)
	}

	var wg sync.WaitGroup
	pipe := make(chan string, 1024)
	for i := 0; i < oMaxCPU; i++ {
		go func(id int) {
			wg.Add(1)
			defer wg.Done()

			for password := range pipe {
				if wallet.CheckPassword(password) {
					fmt.Fprintf(os.Stderr, "OK: %s\n", password)
					os.Exit(0)
				}
			}
		}(i)
	}

	t_start := time.Now()
	elapsed := func() float64 { return time.Now().Sub(t_start).Seconds() }

	count := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		password := scanner.Text()
		pipe <- password
		count++
	}

	close(pipe)
	wg.Wait()
	fmt.Printf("Total: %d, PPS: %.02f\n", count, float64(count)/elapsed())
}
