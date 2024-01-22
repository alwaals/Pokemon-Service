package middlewares

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"pokemon-service/schema"
	"testing"
)

func TestLoggingResponse(t *testing.T) {
	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	logger := schema.Logger{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := LoggingRequest(nextHandler, logger)

	// create a mock request to use
	vars := map[string]string{
		"Name": "PK101",
	}

	req, err := http.NewRequest("GET", "/pokemon-service/{Name}", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, vars)

	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}
