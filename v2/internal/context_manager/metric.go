package context_manager

import (
	"time"
)

type Metric struct {
	Namespace string
	Value     interface{}
	Tags      map[string]string
	Timestamp time.Time
}
