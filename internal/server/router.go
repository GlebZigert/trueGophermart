package server

import (
	"net/http"

	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/logger"
	"github.com/GlebZigert/gophermart/internal/middleware"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func InitRouter() {
	logger.Log.Info("Running server", zap.String("address", config.RunAddr))
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.ErrHandler)
		r.Use(middleware.Log)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Get("/api/user/orders", OrdersGet)
		})

		r.Post("/api/user/register", Register)
	})

	err := http.ListenAndServe(config.RunAddr, r)
	if err != nil {
		logger.Log.Error("ListenAndServe", zap.String("err", err.Error()))
	}
}
