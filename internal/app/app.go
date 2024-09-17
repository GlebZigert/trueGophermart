package app

import (
	"github.com/GlebZigert/trueGophermart/internal/auth"
	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/dblayer"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/middleware"
	"github.com/GlebZigert/trueGophermart/internal/server"
)

func Run() (err error) {

	cfg := config.NewConfig()
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		return err
	}

	db, err := dblayer.NewDB(cfg.DatabaseDSN)
	if err != nil {
		return
	}

	auc := auth.NewAuth(cfg.SECRETKEY, cfg.TOKENEXP)

	mdl := middleware.NewMiddlewares(auc)

	h, err := server.NewServer(db, cfg, mdl)

	if err != nil {
		return
	}
	err = h.Start()

	return
}
