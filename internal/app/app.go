package app

import (
	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/dblayer"
	"github.com/GlebZigert/gophermart/internal/server"
)

func Run() error {
	config.ParseFlags()
	dblayer.Init()
	server.InitRouter()

	return nil
}
