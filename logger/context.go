package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var ctxKey struct{}

func NewContext(ctx context.Context, log *logrus.Entry) context.Context {
	if log == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKey, log)
}

func FromContext(ctx context.Context) *logrus.Entry {
	log, ok := ctx.Value(ctxKey).(*logrus.Entry)
	if ok {
		return log
	}
	return nil
}
