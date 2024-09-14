package server

import (
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/middleware"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func InitRouter(h *handler) {
	logger.Log.Info("Running server", zap.String("address", config.RunAddr))
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.ErrHandler)
		r.Use(middleware.Log)

		//в бизнес-логику нужна аутентификация
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Get("/api/user/orders", h.OrdersGet)
		})

		//в регистрацию-авторизацию не нужна аутентификация
		r.Post("/api/user/register", h.Register)
		r.Post("/api/user/login", h.Login)
	})

	err := http.ListenAndServe(config.RunAddr, r)
	if err != nil {
		logger.Log.Error("ListenAndServe", zap.String("err", err.Error()))
	}
}
