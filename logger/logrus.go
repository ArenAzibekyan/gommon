package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func Standard() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}

type Config struct {
	JSONFormatter bool
	ReportCaller  bool
	Level         string
	Output        string
	NoLock        bool
}

func New(conf *Config) (*logrus.Entry, error) {
	log := logrus.New()
	if conf.JSONFormatter {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
	log.SetReportCaller(conf.ReportCaller)
	l, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(l)
	w, err := output(conf.Output)
	if err != nil {
		return nil, err
	}
	log.SetOutput(w)
	if conf.NoLock {
		log.SetNoLock()
	}
	return logrus.NewEntry(log), nil
}

func output(s string) (io.Writer, error) {
	switch strings.ToLower(s) {
	case "", "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	}
	return file(s)
}

func file(path string) (*os.File, error) {
	if !filepath.IsAbs(path) {
		exec, err := os.Executable()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(filepath.Dir(exec), path)
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}
