package router

import (
	"log"
	"net/http"
	"time"
)

// loggingResponseWriter wraps the http.ResponseWriter to capture status code and response size.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK} // Default status to OK
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

// LogMiddleware creates a middleware that logs HTTP requests and responses.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r) // Serve the actual request

		duration := time.Since(start)

		log.Printf("[%s] %s %s %d %dbytes %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			lrw.statusCode,
			lrw.size,
			duration,
		)
	})
}

/*
mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	// Wrap the mux with the logging middleware
	loggedMux := LogMiddleware(mux)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
*/
