package middleware

import (
	"context"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

func (mdl *Middleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer packerr.AddErrToReqContext(r, &err)
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		authv := r.Header.Get("Authorization")
		mdl.logger.Info("auth: ", map[string]interface{}{
			"auth": authv,
		})

		id, err := mdl.auch.GetUID(authv)

		if err != nil {

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			w.Write([]byte{})
			return

		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, config.UIDkey, id)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)

	})
}
