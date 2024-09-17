package middleware

type IAuth interface {
	BuildJWTString(id int) (string, error)
	GetUserID(tokenString string) (int, error)
}

type Middleware struct {
	auch IAuth
}

func (m *Middleware) GetAuch() IAuth {
	return m.auch
}

func NewMiddlewares(auth IAuth) *Middleware {
	return &Middleware{auth}
}
