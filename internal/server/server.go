package server

import (
	"fmt"
	"net/http"

	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/go-chi/chi"
)

func InitRouter() {
	fmt.Println("starting on address ", config.RunAddr)
	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get req")
		w.Write([]byte("Hello World!"))
	})

	err := http.ListenAndServe(config.RunAddr, r)
	if err != nil {
		fmt.Println(err.Error())
	}
}
