package plugin

import "time"

type Tags map[string]string

type Namespace string

type MetricValue interface{}

type Metric struct {
	Namespace Namespace
	Value     MetricValue
	Tags      Tags
	Timestamp time.Time
}
