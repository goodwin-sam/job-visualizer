// package arguments provides tests for command-line argument parsing functionality
package arguments

import (
	"flag"
	"os"
	"testing"
)

// TestCheckCLIArguments tests the CheckCLIArguments function with different flag combinations
func TestCheckCLIArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			// tests default behavior when no headless flag is provided
			name:     "no headless flag",
			args:     []string{"program"},
			expected: false,
		},
		{
			// tests explicit headless flag set to true
			name:     "headless flag set to true",
			args:     []string{"program", "-headless=true"},
			expected: true,
		},
		{
			// tests explicit headless flag set to false
			name:     "headless flag set to false",
			args:     []string{"program", "-headless=false"},
			expected: false,
		},
		{
			// tests headless flag without value (should default to true)
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
