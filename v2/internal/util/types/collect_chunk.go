package types

type CollectChunk struct {
	Metrics  []*Metric
	Warnings []Warning
	Err      error
}
