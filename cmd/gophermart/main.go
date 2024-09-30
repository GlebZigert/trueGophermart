package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/app"
)

func main() {
	err := app.Run()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

}
