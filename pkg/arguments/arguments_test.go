package arguments

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
			name:     "no headless flag",
			args:     []string{"program"},
			expected: false,
		},
		{
			name:     "headless flag set to true",
			args:     []string{"program", "-headless=true"},
			expected: true,
		},
		{
			name:     "headless flag set to false",
			args:     []string{"program", "-headless=false"},
			expected: false,
		},
		{
			name:     "headless flag without value (should default to true)",
			args:     []string{"program", "-headless"},
			expected: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() {
				os.Args = originalArgs
			}()
			os.Args = testCase.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			result := CheckCLIArguments()
			if result != testCase.expected {
				t.Errorf("CheckCLIArguments() = %v, want %v", result, testCase.expected)
			}
		})
	}
}
