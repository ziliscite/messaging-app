package middleware

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"go.elastic.co/apm"
	"log"
	"net/http"
	"time"
)

func CustomLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, spanCtx := apm.StartSpan(r.Context(), "logger middleware", "middleware")
		defer span.End()

		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			duration := time.Since(start)

			// Get the content length
			contentLength := ww.BytesWritten()
			contentLengthStr := fmt.Sprintf("%dB", contentLength)

			// Format the duration
			var durationStr string
			if duration < time.Second {
				durationStr = fmt.Sprintf("%.4fms", duration.Seconds()*1000)
			} else {
				durationStr = duration.String()
			}

			log.Printf(
				"%q from %s - %d %s in %s\n",
				fmt.Sprintf("%s %s %s", r.Method, r.URL.String(), r.Proto),
				r.RemoteAddr, ww.Status(), contentLengthStr, durationStr,
			)
		}()

		next.ServeHTTP(ww, r.WithContext(spanCtx))
	})
}
