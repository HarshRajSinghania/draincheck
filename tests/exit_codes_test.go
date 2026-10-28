//go:build !windows

package tests

import (
	"os/exec"
	"strings"
	"testing"
)

func TestExitCodesPositionalArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"without debug", []string{"--nodes", "n1", "n2"}},
		{"with debug", []string{"--debug", "--nodes", "n1", "n2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("expected non-zero exit")
			}
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("not ExitError: %v", err)
			}
			if exitErr.ExitCode() != 2 {
				t.Fatalf("exit code = %d, want 2", exitErr.ExitCode())
			}
			if !strings.Contains(string(out), `unexpected argument "n2"`) {
				t.Fatalf("output = %q, want to name n2", out)
			}
			// No report on stdout (CombinedOutput mixes, but check no finding-like text or just the error)
			if strings.Contains(string(out), "BLOCKED") || strings.Contains(string(out), "findings") {
				t.Fatalf("unexpected report in output: %s", out)
			}
		})
	}
}
