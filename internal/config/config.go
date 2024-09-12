package config

import (
	"flag"
	"fmt"
	"os"
)

type key int

const (
	UIDkey key = iota
	JWTkey key = iota
	NEWkey key = iota
	Errkey key = iota
	// ...
)

var (
	RunAddr      string
	FlagLogLevel string
	DatabaseDSN  string
	TOKENEXP     int

	SECRETKEY string
)

func ParseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&DatabaseDSN, "d", "", "database dsn")

	flag.StringVar(&SECRETKEY, "SECRETKEY", "supersecretkey", "ключ")
	flag.IntVar(&TOKENEXP, "TOKENEXP", 3, "время жизни токена в часах")
	flag.Parse()

	fmt.Println("RunAddr", RunAddr)
	fmt.Println("DatabaseDSN", DatabaseDSN)
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		RunAddr = envRunAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		DatabaseDSN = envDatabaseDSN
	}

	fmt.Println("RunAddr", RunAddr)
	fmt.Println("DatabaseDSN", DatabaseDSN)
}
