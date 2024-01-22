package middlewares

import (
	schema "pokemon-service/schema"
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

type Middleware func(http.HandlerFunc, schema.Logger) http.HandlerFunc

type ResponseWriterWrapper struct {
    w          http.ResponseWriter
    body       bytes.Buffer
    statusCode int
}

// NewResponseWriterWrapper static function creates a wrapper for the http.ResponseWriter
func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
    var buf bytes.Buffer
    var statusCode int = 200
    return &ResponseWriterWrapper{
        w:          w,
        body:       buf,
        statusCode: statusCode,
    }
}

func (rww *ResponseWriterWrapper) Write(buf []byte) (int, error) {
    rww.body.Write(buf)
    return (rww.w).Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function
func (rww *ResponseWriterWrapper) Header() http.Header {
    return (rww.w).Header()

}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
    (rww.statusCode) = statusCode
    (rww.w).WriteHeader(statusCode)
}

func (rww *ResponseWriterWrapper) String() string {
    var buf bytes.Buffer
    for k, v := range (rww.w).Header() {
        buf.WriteString(fmt.Sprintf("%s: %v", k, v))
    }
    buf.WriteString(fmt.Sprintf(" Status Code: %d", (rww.statusCode)))
    buf.WriteString(rww.body.String())
    return buf.String()
}
func LoggingResponse(handler http.HandlerFunc, l schema.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				l.InfoLogger.Println(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()
		wrapped := ResponseWriterWrapper{w: w}
		handler.ServeHTTP(&wrapped, r)
		l.InfoLogger.Println(
			wrapped.statusCode,
			wrapped.String(),
		)
	}
}
func Chain(handler http.HandlerFunc, logger schema.Logger, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		handler = m(handler, logger)
	}
	return handler
}
