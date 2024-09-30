package middleware

import "github.com/GlebZigert/trueGophermart/internal/logger"

type IAuth interface {
	BuildJWTString(id int) (string, error)
	GetUID(tokenString string) (int, error)
}

type Middleware struct {
	auch   IAuth
	logger logger.Logger
}

func (mdl *Middleware) GetAuch() IAuth {
	return mdl.auch
}

func NewMiddlewares(auth IAuth, logger logger.Logger) *Middleware {
	return &Middleware{auth, logger}
}
