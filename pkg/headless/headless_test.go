package headless

import (
	"flag"
	"os"
	"testing"
)

func TestCheckCLIArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "no arguments",
			args:     []string{},
			expected: false,
		},
		{
			name:     "headless flag true",
			args:     []string{"-headless"},
			expected: true,
		},
		{
			name:     "headless flag false",
			args:     []string{"-headless=false"},
			expected: false,
		},
		{
			name:     "headless flag true explicit",
			args:     []string{"-headless=true"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag state for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			
			// Set up os.Args for this test
			oldArgs := os.Args
			os.Args = append([]string{"test"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			result := CheckCLIArguments()
			if result != tt.expected {
				t.Errorf("CheckCLIArguments() = %v, expected %v", result, tt.expected)
			}
		})
	}
}