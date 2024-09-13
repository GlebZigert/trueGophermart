package app

import (
	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/dblayer"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/server"
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
