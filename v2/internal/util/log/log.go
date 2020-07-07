package log

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

const loggerKey = "logger"

func WithCtx(ctx context.Context) *logrus.Entry {
	var log *logrus.Entry
	var err error

	log, err = withCtx(ctx)
	if err != nil {
		log = logrus.WithFields(logrus.Fields{}) // default: logger
		log.WithError(err).Info("Can't get logger from context, fallback to default")
	}

	return log
}

func withCtx(ctx context.Context) (*logrus.Entry, error) {
	logI := ctx.Value(loggerKey)
	if logI == nil {
		return nil, errors.New("no logger in context")
	}

	log, ok := logI.(*logrus.Entry)
	if !ok {
		return nil, errors.New("invalid logger type")
	}

	return log, nil
}
