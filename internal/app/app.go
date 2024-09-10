package app

import (
	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/server"
)

func Run() error {
	config.ParseFlags()
	server.InitRouter()

	return nil
}
