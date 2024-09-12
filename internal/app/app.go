package app

import (
	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/dblayer"
	"github.com/GlebZigert/gophermart/internal/logger"
	"github.com/GlebZigert/gophermart/internal/server"
	"go.uber.org/zap"
)

func Run() error {

	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	config.ParseFlags()

	err := dblayer.Init()
	if err != nil {
		logger.Log.Error("dblayer.Init: ", zap.String("", err.Error()))
		return err
	}
	server.InitRouter()

	return nil
}
