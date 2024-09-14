package app

import (
	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/dblayer"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/server"
)

func Run() error {

	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	config.ParseFlags()

	db := dblayer.Init()
	h := server.New(db)

	server.InitRouter(h)

	return nil
}
