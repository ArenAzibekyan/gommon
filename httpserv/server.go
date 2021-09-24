package httpserv

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
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

func new(conf *Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    Addr("", conf.Port),
		Handler: handler,
	}
}

func start(serv *http.Server, log *logrus.Entry) {
	if log == nil {
		log = logger.Standard()
	}
	log = log.WithField("addr", serv.Addr)
	log.Info("Starting HTTP server")
	if err := serv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func stop(ctx context.Context, serv *http.Server) {
	<-ctx.Done()
	serv.Close()
}

func StartContext(ctx context.Context, conf *Config, handler http.Handler, log *logrus.Entry) {
	serv := new(conf, handler)
	go stop(ctx, serv)
	start(serv, log)
}

func Start(conf *Config, handler http.Handler, log *logrus.Entry) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	StartContext(ctx, conf, handler, log)
}
