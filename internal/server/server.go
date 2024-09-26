package server

import (
	"errors"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/middleware"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type IAuth interface {
	BuildJWTString(id int) (string, error)
	GetUserID(tokenString string) (int, error)
}

type Server struct {
	auch   IAuth
	DB     *gorm.DB
	cfg    *config.Config
	mdl    *middleware.Middleware
	logger logger.Logger
}

var errNoAuthMiddleware = errors.New("В миддлеварах не определен auth")

func NewServer(db *gorm.DB, cfg *config.Config, mdl *middleware.Middleware, logger logger.Logger) (*Server, error) {
	auch := mdl.GetAuch()
	if auch == nil {
		return nil, errNoAuthMiddleware
	}
	return &Server{auch, db, cfg, mdl, logger}, nil
}

func (srv *Server) Start() (err error) {
	srv.logger.Info("Running server", map[string]interface{}{
		"address": srv.cfg.RunAddr,
	})

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(srv.mdl.ErrHandler)
		r.Use(srv.mdl.Log)

		//в бизнес-логику нужна аутентификация
		r.Group(func(r chi.Router) {
			r.Use(srv.mdl.Auth)
			r.Get("/api/user/orders", srv.OrderGet)
			r.Post("/api/user/orders", srv.OrderPost)
			r.Get("/api/user/balance", srv.BalanceGet)
			r.Post("/api/user/balance/withdraw", srv.Withdraw)
			r.Get("/api/user/withdrawals", srv.WithdrawalsGet)

		})

		//в регистрацию-авторизацию не нужна аутентификация
		r.Post("/api/user/register", srv.Register)
		r.Post("/api/user/login", srv.Login)
	})

	err = http.ListenAndServe(srv.cfg.RunAddr, r)
	return
}
