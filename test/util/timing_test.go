package util_test

import (
	"testing"
	"time"
	"github.com/lunjon/go-blast/pkg/util"
)

func TestTimeFromFrequency(t *testing.T) {
	tests := []struct {
		name         string
		frequency    float64
		expected     time.Duration
		expectsError bool
	}{
		{"10 Hz", 10, 100 * time.Millisecond, false},
		{"1 Hz", 1, 1000 * time.Millisecond, false},
		{"0 Hz", 0, time.Millisecond, true},
		{"0 Hz", -1, time.Millisecond, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.TimeFromFrequency(tt.frequency)

			if tt.expectsError && err == nil {
				t.Errorf("TimeFromFrequency() = %v, want %v", got, tt.expected)
			} else if !tt.expectsError && got != tt.expected {
				t.Errorf("TimeFromFrequency() expected to return error for %f but did not", tt.frequency)
			}
		})
	}
}
