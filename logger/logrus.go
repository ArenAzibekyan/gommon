package logger

import (
	"os"
	"path/filepath"

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
	}
	log.SetReportCaller(conf.ReportCaller)
	err := setLevel(log, conf)
	if err != nil {
		return nil, err
	}
	err = setOutput(log, conf)
	if err != nil {
		return nil, err
	}
	if conf.NoLock {
		log.SetNoLock()
	}
	return logrus.NewEntry(log), nil
}

func setLevel(log *logrus.Logger, conf *Config) error {
	lvl, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	return nil
}

func setOutput(log *logrus.Logger, conf *Config) error {
	if conf.Output == "" {
		log.SetOutput(os.Stdout)
		return nil
	}
	f, err := openFile(conf.Output)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}

func openFile(path string) (*os.File, error) {
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
