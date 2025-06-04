package api

import (
	"testing"
)

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name     string
		args     QueryArgs
		expected string
	}{
		{
			name:     "Basic query",
			args:     QueryArgs{MinStars: 10, MaxStars: 100},
			expected: "is:public fork:true stars:10..100 ",
		},
		{
			name:     "Zero stars range",
			args:     QueryArgs{MinStars: 0, MaxStars: 0},
			expected: "is:public fork:true stars:0..0 ",
		},
		{
			name:     "Large stars range",
			args:     QueryArgs{MinStars: 1000, MaxStars: 50000},
			expected: "is:public fork:true stars:1000..50000 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQuery(tt.args)
			if got != tt.expected {
				t.Errorf("buildQuery(%+v) = %q, want %q", tt.args, got, tt.expected)
			}
		})
	}
}
