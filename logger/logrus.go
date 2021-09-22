package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	JSONFormatter bool
	Level         string
	Output        string
	NoLock        bool
}

func New(conf *Config) (*logrus.Entry, error) {
	f := formatter(conf.JSONFormatter)
	l, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	w, err := output(conf.Output)
	if err != nil {
		return nil, err
	}
	log := new(f, l, w)
	if conf.NoLock {
		log.SetNoLock()
	}
	return logrus.NewEntry(log), nil
}

func Standard() *logrus.Entry {
	log := new(formatter(false), logrus.DebugLevel, os.Stdout)
	return logrus.NewEntry(log)
}

func new(f logrus.Formatter, l logrus.Level, w io.Writer) *logrus.Logger {
	return &logrus.Logger{
		Formatter: f,
		Level:     l,
		Out:       w,
		Hooks:     make(logrus.LevelHooks),
		ExitFunc:  os.Exit,
	}
}

func formatter(json bool) logrus.Formatter {
	if json {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{FullTimestamp: true}
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
