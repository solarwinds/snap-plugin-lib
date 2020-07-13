package runner

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/log"
)

var moduleFields = logrus.Fields{"layer": "lib", "module": "plugin-runner"}

func logger(ctx context.Context) logrus.FieldLogger {
	return log.WithCtx(ctx).WithFields(moduleFields)
}
