package middlewares

import (
	schema "pokemon-service/schema"
	utility "pokemon-service/utility"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func LoggingRequest(handler http.HandlerFunc,l schema.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		l.InfoLogger.Println("Received request with:", req.URL.Path+" and Method:"+req.Method)
		var pokeResp schema.PokemonResponse
		reqBytes, err := io.ReadAll(req.Body)
		if err != nil {
			utility.FrameHttpResponse(400, "Invalid Json request", &pokeResp, start, w)
			return
		}
		l.InfoLogger.Println("Request:", string(reqBytes))
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
		handler.ServeHTTP(w, req)
	}
}
func Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(w, req)
	}
}
