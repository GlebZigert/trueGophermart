package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
		body   string
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter

	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	r.responseData.body = string(b)
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func (mdl *Middleware) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer packerr.AddErrToReqContext(r, &err)
		t1 := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		id, ok := r.Context().Value(config.UIDkey).(int)
		if ok {
			mdl.logger.Info("auth: ", map[string]interface{}{
				"id": id,
			})

		}

		//
		/*
			jwt, ok := r.Context().Value(config.JWTkey).(string)
			if ok {

				lw.Header().Add("Authorization", string(jwt))
			}
		*/
		h.ServeHTTP(&lw, r)

		mdl.logger.Info("got incoming HTTP request: ", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"dt":     time.Since(t1).String(),
			"size":   strconv.Itoa(responseData.size),
			"UID":    id,
			"status": strconv.Itoa(responseData.status),
			"body":   responseData.body,
		})

	})
}
