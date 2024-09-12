package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/gophermart/internal/packerr"
	"github.com/GlebZigert/gophermart/internal/users"
)

func Register(w http.ResponseWriter, req *http.Request) {
	var err error
	defer packerr.AddErrToReqContext(req, &err)

	var user users.User

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return //err
	}

	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return // err
	}

	//если пришла правильная посылка - возвращаю 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
