package server

import (
	"net/http"

	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/go-chi/chi"
)

func InitRouter() {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(config.RunAddr, r)
}
