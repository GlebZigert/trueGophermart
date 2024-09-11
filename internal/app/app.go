package app

import (
	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/logger"
	"github.com/GlebZigert/gophermart/internal/server"
)

func Run() error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	config.ParseFlags()
	server.InitRouter()

	return nil
}
