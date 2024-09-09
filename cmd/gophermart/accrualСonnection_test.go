package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAccrualConnetion(t *testing.T) {
	t.Logf("its a test")

	response, err := http.Get("http://localhost:8080/api/orders/0")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Status Code: %d\r\n", response.StatusCode)
	for k, v := range response.Header {
		// заголовок может иметь несколько значений,
		// но для простоты запросим только первое
		fmt.Printf("%s: %v\r\n", k, v[0])
	}
}
