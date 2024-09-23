package config

import (
	"flag"
	"log"
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

// Person — структура, описывающая человека.
type Config struct {
	RunAddr      string
	FlagLogLevel string
	DatabaseDSN  string
	TOKENEXP     int

	SECRETKEY string

	AccrualAddress string
}

// NewPerson возвращает новую структуру Person.
func NewConfig() *Config {
	cfg := Config{}
	cfg.ParseFlags()
	return &cfg
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.RunAddr, "a", "localhost:8081", "address and port to run server")
	flag.StringVar(&c.FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&c.AccrualAddress, "r", "http://localhost:8080", "accrual system address")

	flag.StringVar(&c.SECRETKEY, "SECRETKEY", "supersecretkey", "ключ")
	flag.IntVar(&c.TOKENEXP, "TOKENEXP", 3, "время жизни токена в часах")
	flag.Parse()

	log.Println("RunAddr", c.RunAddr)
	log.Println("DatabaseDSN", c.DatabaseDSN)
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		c.RunAddr = envRunAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		c.DatabaseDSN = envDatabaseDSN
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		c.DatabaseDSN = envDatabaseDSN
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		c.AccrualAddress = envAccrualSystemAddress
	}

	log.Println("RunAddr", c.RunAddr)
	log.Println("DatabaseDSN", c.DatabaseDSN)
}
