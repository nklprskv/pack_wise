package app

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

const maxLoggedBodyBytes = 4096

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if w.body.Len() < maxLoggedBodyBytes {
		remaining := maxLoggedBodyBytes - w.body.Len()
		if len(data) > remaining {
			w.body.Write(data[:remaining])
		} else {
			w.body.Write(data)
		}
	}

	return w.ResponseWriter.Write(data)
}

// loggingMiddleware records request and response payloads together with the
// status code and request duration for operational visibility.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		requestBody := readRequestBody(r)
		recorder := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		log.Printf(
			"%s %s status=%d duration=%s requestBody=%q responseBody=%q",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			time.Since(startedAt),
			requestBody,
			trimLoggedBody(recorder.body.String()),
		)
	})
}

func readRequestBody(r *http.Request) string {
	if r.Body == nil {
		return ""
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxLoggedBodyBytes+1))
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(nil))
		return "<failed to read request body>"
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	return trimLoggedBody(string(body))
}

func trimLoggedBody(body string) string {
	if len(body) <= maxLoggedBodyBytes {
		return body
	}

	return body[:maxLoggedBodyBytes] + "...(truncated)"
}
