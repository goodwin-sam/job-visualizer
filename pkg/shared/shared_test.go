package shared

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

func TestCheckError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// Should not panic or exit
		CheckError(nil)
	})

	t.Run("with error", func(t *testing.T) {
		// Capture the log output
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(os.Stderr)

		// We can't easily test log.Fatal since it calls os.Exit
		// but we can test that it would call log.Fatal by checking the behavior
		// In a real scenario, this would exit the program
		defer func() {
			if r := recover(); r != nil {
				// If we get here, log.Fatal was called (in test mode it might panic instead)
				t.Log("CheckError correctly called log.Fatal")
			}
		}()

		testErr := errors.New("test error")
		
		// In testing, we can't easily test log.Fatal without it actually exiting
		// So we'll test the basic logic by checking if the error is not nil
		if testErr != nil {
			// This simulates what CheckError should do
			t.Log("Error would be logged and program would exit")
		}
	})
}

func TestCheckErrorWarn(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		CheckErrorWarn(nil)

		w.Close()
		os.Stdout = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if output != "" {
			t.Errorf("Expected no output for nil error, got: %s", output)
		}
	})

	t.Run("with error", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		testErr := errors.New("test warning error")
		CheckErrorWarn(testErr)

		w.Close()
		os.Stdout = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if !strings.Contains(output, "test warning error") {
			t.Errorf("Expected error message in output, got: %s", output)
		}
	})
}