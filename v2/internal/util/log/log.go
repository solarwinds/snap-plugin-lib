/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

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
