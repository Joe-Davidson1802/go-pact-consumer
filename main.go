package main

import (
	"net/http"

	"github.com/joe-davidson1802/go-pact-consumer/handlers"
)

func main() {
	h := handlers.GetTimeHandler("http://localhost:8081")
	panic(http.ListenAndServe(":8080", h))
}
