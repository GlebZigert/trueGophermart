package config

import (
	"flag"
	"os"
)

var (
	RunAddr string
)

func ParseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		RunAddr = envRunAddr
	}
}
