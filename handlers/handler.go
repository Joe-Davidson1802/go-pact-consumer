package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joe-davidson1802/go-pact-consumer/models"
	"github.com/joe-davidson1802/go-pact-consumer/templates"
)

func GetTimeHandler(apiUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		response, err := http.DefaultClient.Get(apiUrl)
		if err != nil || response.StatusCode != 200 {
			fmt.Printf("error occured %v", err)
			w.WriteHeader(500)
			return
		}
		var times []models.TimeResponse
		err = json.NewDecoder(response.Body).Decode(&times)
		if err != nil {
			fmt.Printf("error occured %v", err)
			w.WriteHeader(500)
			return
		}
		templates.TimePage(times).Render(r.Context(), w)
	}
}
