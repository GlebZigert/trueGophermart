package app

import (
	"context"

	"github.com/GlebZigert/trueGophermart/internal/auth"
	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/dblayer"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/middleware"
	"github.com/GlebZigert/trueGophermart/internal/server"
)

func Run() (err error) {

	cfg := config.NewConfig()
	ctx := context.TODO()
	logger := logger.NewLogrusLogger(cfg.FlagLogLevel, ctx)

	db, err := dblayer.NewDB(cfg.DatabaseDSN)
	if err != nil {
		return
	}

	auc := auth.NewAuth(cfg.SECRETKEY, cfg.TOKENEXP)

	mdl := middleware.NewMiddlewares(auc, logger)

	h, err := server.NewServer(db, cfg, mdl, logger)

	if err != nil {
		return
	}
	err = h.Start()

	return
}
