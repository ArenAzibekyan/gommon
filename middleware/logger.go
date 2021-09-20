package middleware

import (
	"io"
	"net/http"
	"strings"

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

func WriteLog(withBody, withHeader bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			if withBody {
				log = log.WithField("body", bodyString(r))
			}
			if withHeader {
				log = log.WithFields(headerFields(r))
			}
			log.WithFields(logrus.Fields{
				"method": r.Method,
				"proto":  r.Proto,
				"from":   r.RemoteAddr,
				"url":    r.RequestURI,
			}).Info("req")
		})
	}
}

func bodyString(r *http.Request) string {
	b := new(strings.Builder)
	_, err := io.Copy(b, r.Body)
	if err != nil {
		return ""
	}
	return b.String()
}

func headerFields(r *http.Request) logrus.Fields {
	m := make(map[string]interface{}, len(r.Header))
	for k := range r.Header {
		m[k] = r.Header.Get(k)
	}
	return m
}
