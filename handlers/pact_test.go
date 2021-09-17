package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

func TestConsumer(t *testing.T) {
	// Setup pact config
	pact := &dsl.Pact{
		Consumer: "time-front-end",
		Provider: "time-back-end",
		Host:     "0.0.0.0",
	}
	defer pact.Teardown()
	pact.Setup(false)

	pact.AddInteraction().
		UponReceiving("a time request").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String("/"),
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Body: dsl.EachLike(
				dsl.MapMatcher{
					"zone": dsl.String("UTC"),
					"time": dsl.Timestamp(),
				}, 1),
		})

	// Run server against pact mock
	test := func() error {
		h := GetTimeHandler(fmt.Sprintf("http://0.0.0.0:%d", pact.Server.Port))
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		h.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Fatalf("expected 200 response got %d", rec.Code)
		}

		return nil
	}
	if err := pact.Verify(test); err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
	pact.WritePact()
}
