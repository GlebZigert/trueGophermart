package middleware

import (
	"context"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
)

type MyHandlerFunc func(w http.ResponseWriter, r *http.Request)

// Implement the http.Handler interface.
func (fn MyHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r) // Call handler function.
}

func (mdl *Middleware) ErrHandler(f http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		//помещаем в контекст реквеста указатель на ошибку
		ctx := context.WithValue(r.Context(), config.Errkey, &err)
		r = r.WithContext(ctx)
		f.ServeHTTP(w, r)

		if err != nil {
			mdl.logger.Error("auth: ", map[string]interface{}{
				"err": err.Error(),
			})
		}
	})
}
