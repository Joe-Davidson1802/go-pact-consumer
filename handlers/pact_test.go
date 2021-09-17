package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

const PACT_DIR = "./pacts"

var version = os.Getenv("VERSION")
var pactbroker = os.Getenv("PACT_BROKER_BASE_URL")
var pact *dsl.Pact

func TestConsumer(t *testing.T) {
	// Setup pact config
	pact = &dsl.Pact{
		Consumer: "time-front-end",
		Provider: "time-back-end",
		Host:     "0.0.0.0",
		PactDir:  PACT_DIR,
	}
	defer pact.Teardown()
	pact.Setup(false)

	// Create interaction
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

	// Write JSON file
	pact.WritePact()

	publish()
}

func publish() {
	// Publish pact to pactflow.io
	p := dsl.Publisher{
		LogLevel: "DEBUG",
	}
	contractName := fmt.Sprintf("%s/%s-%s.json", PACT_DIR, pact.Consumer, pact.Provider)
	err := p.Publish(types.PublishRequest{
		PactURLs:        []string{filepath.FromSlash(contractName)},
		PactBroker:      pactbroker,
		ConsumerVersion: version,
		Tags:            []string{"main"},
	})
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
}
