package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitRouter() {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe("localhost:8081", r)
}
