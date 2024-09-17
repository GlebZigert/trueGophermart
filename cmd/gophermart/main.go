package main

import (
	"log"

	"github.com/GlebZigert/trueGophermart/internal/app"
)

func main() {
	err := app.Run()
	log.Fatalln(err.Error())
}
