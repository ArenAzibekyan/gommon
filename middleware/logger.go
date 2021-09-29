package middleware

import (
	"net/http"

	"github.com/ArenAzibekyan/gommon/logger"
	"github.com/sirupsen/logrus"
)

func AddLogger(log *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.NewContext(r.Context(), log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WriteLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.WithFields(logger.RequestFields(r)).Info("req")
		next.ServeHTTP(w, r)
	})
}
