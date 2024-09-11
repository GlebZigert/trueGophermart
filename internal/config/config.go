package config

import (
	"flag"
	"os"
)

var (
	RunAddr      string
	FlagLogLevel string
)

func ParseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		RunAddr = envRunAddr
	}
}
