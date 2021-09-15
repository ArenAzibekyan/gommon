package httpserv

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/ArenAzibekyan/gommon/logger"
	"github.com/sirupsen/logrus"
)

func Addr(host string, port uint16) string {
	if port == 0 {
		return host
	}
	return net.JoinHostPort(host, strconv.Itoa(int(port)))
}

type Config struct {
	Port uint16
}

type Server struct {
	srv *http.Server
	ctx context.Context
	log *logrus.Entry
}

func New(ctx context.Context, conf *Config, handler http.Handler, log *logrus.Entry) *Server {
	srv := &http.Server{
		Addr:    Addr("", conf.Port),
		Handler: handler,
	}
	if log == nil {
		log = logger.Standard()
	}
	log = log.WithField("addr", srv.Addr)
	return &Server{srv, ctx, log}
}

func (s *Server) Start() {
	go func() {
		<-s.ctx.Done()
		s.srv.Close()
	}()
	go func() {
		s.log.Info("Starting to listen")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal(err)
		}
	}()
}
