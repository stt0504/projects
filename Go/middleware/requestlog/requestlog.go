//go:build !solution

package requestlog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func generateKey(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateKey(r.URL.Path)

			l.Info("request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("request_id", requestID),
			)

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			ctx := r.Context()
			ctx = context.WithValue(ctx, requestKey{}, requestID)

			start := time.Now()
			defer func() {
				duration := time.Since(start)
				if err := recover(); err != nil {
					l.Error("request panicked",
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("request_id", requestID),
						zap.Any("error", err),
					)
					panic(err)
				}

				l.Info("request finished",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("request_id", requestID),
					zap.Int("status_code", rw.statusCode),
					zap.Duration("duration", duration),
				)
			}()

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

type requestKey struct{}
