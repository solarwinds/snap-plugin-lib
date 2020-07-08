package log

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const loggerCtxKey ctxKey = "logger"

func WithCtx(ctx context.Context) logrus.FieldLogger {
	var log logrus.FieldLogger
	var err error

	log, err = withCtx(ctx)
	if err != nil {
		log = logrus.WithFields(logrus.Fields{}) // default: logger
		log.WithError(err).Info("Can't get logger from context, fallback to default")
	}

	return log
}

func ToCtx(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func withCtx(ctx context.Context) (logrus.FieldLogger, error) {
	logI := ctx.Value(loggerCtxKey)
	if logI == nil {
		return nil, errors.New("no logger in context")
	}

	log, ok := logI.(logrus.FieldLogger)
	if !ok {
		return nil, errors.New("invalid logger type")
	}

	return log, nil
}
