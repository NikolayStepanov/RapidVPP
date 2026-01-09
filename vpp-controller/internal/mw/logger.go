package mw

import (
	"net/http"
	"time"

	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.uber.org/zap"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logger.Debug("http",
			zap.String("component", "handler"),
			zap.String("method", r.Method),
			zap.String("url", r.RequestURI),
			zap.String("remote address", r.RemoteAddr),
			zap.Duration("duration", duration),
		)
	})
}
