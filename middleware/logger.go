package middleware

import (
	"net/http"
	"time"

	"github.com/ArenAzibekyan/gommon/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AddLogger(log *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.WithField("reqID", uuid.NewString())
			ctx := logger.NewContext(r.Context(), log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WriteLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		log := logger.FromContext(r.Context())
		log.WithFields(logger.RequestFields(r)).Info("Req")
		next.ServeHTTP(w, r)
		log.WithField("durat", time.Since(t).String()).Info("Took")
	})
}
