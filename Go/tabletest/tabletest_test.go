package tabletest

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  time.Duration
		expectErr bool
	}{
		{"empty string", "", 0, true},
		{"zero duration", "0", 0, false},
		{"unknown unit", "10x", 0, true},
		{"overflow case", "1000000000000h", 0, true},
		{"invalid format", "abc", 0, true},
		{"decimal without unit", "1.4", 0, true},
		{"negative duration with mixed units", "-1h32m", -92 * time.Minute, false},
		{"leading int overflow", "9223372036854775808", 0, true},
		{"leading int overflow2", "922337203685477582000", 0, true},
		{"fraction overflow", "0.9999999999999999999s", 1 * time.Second, false},
		{"fraction overflow large number", "0.9223372036854775809s", 922337203 * time.Nanosecond, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error status for input %q: got %v, want error: %v", tt.input, err, tt.expectErr)
			}
			if !tt.expectErr && result != tt.expected {
				t.Errorf("unexpected result for input %q: got %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
