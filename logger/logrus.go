package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Standard() *logrus.Entry {
	conf := &Config{Level: "info"}
	log, _ := New(conf)
	return log
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

func output(path string) (io.Writer, error) {
	switch strings.ToLower(path) {
	case "", "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	}
	path, err := abs(path)
	if err != nil {
		return nil, err
	}
	w := &lumberjack.Logger{Filename: path}
	return w, nil
}

func abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	path = filepath.Join(filepath.Dir(ex), path)
	return path, nil
}
