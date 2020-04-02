package util

import (
	"fmt"
	"time"
)

// TimeFromFrequency returns the corresponding time.Duration
// of a frequency f (i.e. the inverse of f). It will be truncated to
// nearest millisecond.
func TimeFromFrequency(f float64) (time.Duration, error) {
	if !(f > 0) {
		return 0, fmt.Errorf("frequency must be greater than 0")
	}

	ms := int((1 / f) * 1000)
	return time.Duration(ms) * time.Millisecond, nil
}
