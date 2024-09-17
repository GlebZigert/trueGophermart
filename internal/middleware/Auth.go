package middleware

import (
	"context"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
	"go.uber.org/zap"
)

func (mdl *Middleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer packerr.AddErrToReqContext(r, &err)
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		authv := r.Header.Get("Authorization")
		logger.Log.Info("auth: ", zap.String("", authv))

		id, err := mdl.auch.GetUserID(authv)

		if err != nil {

			logger.Log.Error("Auth: ", zap.String("", err.Error()))

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
