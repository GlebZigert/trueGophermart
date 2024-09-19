package middleware

import "github.com/GlebZigert/trueGophermart/internal/logger"

type IAuth interface {
	BuildJWTString(id int) (string, error)
	GetUserID(tokenString string) (int, error)
}

type Middleware struct {
	auch   IAuth
	logger logger.Logger
}

func (m *Middleware) GetAuch() IAuth {
	return m.auch
}

func NewMiddlewares(auth IAuth, logger logger.Logger) *Middleware {
	return &Middleware{auth, logger}
}
