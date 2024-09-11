package server

import (
	"fmt"
	"net/http"

	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/logger"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func InitRouter() {
	logger.Log.Info("Running server", zap.String("address", config.RunAddr))
	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get req")
		w.Write([]byte("Hello World!"))
	})

	err := http.ListenAndServe(config.RunAddr, r)
	if err != nil {
		logger.Log.Error("ListenAndServe", zap.String("err", err.Error()))
	}
}
