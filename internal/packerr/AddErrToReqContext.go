package packerr

import (
	"errors"
	"net/http"

	"github.com/GlebZigert/gophermart/internal/config"
)

// Метод для добавления ошибки в контекст реквеста
func AddErrToReqContext(r *http.Request, err *error) {
	if err == nil {
		return
	}

	ctxerr, ok := r.Context().Value(config.Errkey).(*error)
	if ok {

		*ctxerr = errors.Join(*ctxerr, *err)
	}
}

//Метод для чтения ошибки из контекста реквеста
