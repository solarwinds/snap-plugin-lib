package types

import "time"

// Contains additional information about warnings raised during collect / publish processes
type ProcessingError struct {
	Error    error
	Warnings []Warning
}

type Warning struct {
	Message   string
	Timestamp time.Time
}
