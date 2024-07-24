package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type key int

const requestIDKey key = 0

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getLogFilePath(log_dir string) string {

	log_uuid, err := uuid.NewV7()
	if err != nil {
		log_uuid = uuid.New()
		log.Printf("error generating v7 uuid; will use regular uuid: %v", err)
	}
	logfile := fmt.Sprintf("log-%s.txt", log_uuid.String())
	logpath := fmt.Sprintf("%s/%s", log_dir, logfile)
	return logpath
}
