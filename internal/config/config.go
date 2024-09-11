package config

import (
	"flag"
	"os"
)

var (
	RunAddr     string
	DatabaseDSN string
)

func ParseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&DatabaseDSN, "d", "", "database dsn")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		RunAddr = envRunAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); DatabaseDSN != "" {
		DatabaseDSN = envDatabaseDSN
	}
}
