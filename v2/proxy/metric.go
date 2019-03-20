package proxy

import (
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type Metric struct {
	Namespace string
	Value     interface{}
	Tags      plugin.Tags
	Timestamp time.Time
}
