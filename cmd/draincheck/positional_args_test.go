package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRejectPositionalArguments(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
		wantOK  bool
	}{
		{
			name:    "one argument",
			args:    []string{"n1"},
			wantErr: `unexpected argument "n1"`,
		},
		{
			name:    "several arguments",
			args:    []string{"n1", "n2", "n3"},
			wantErr: `unexpected argument "n1"`,
		},
		{
			name:    "empty-string argument",
			args:    []string{""},
			wantErr: `unexpected argument ""`,
		},
		{
			name:    "argument after --",
			args:    []string{"--nodes", "n1", "--", "n2"},
			wantErr: `unexpected argument "n2"`,
		},
		{
			name:    "argument next to a valid --nodes",
			args:    []string{"--nodes", "n1", "n2"},
			wantErr: `unexpected argument "n2"`,
		},
		{
			name:   "no argument",
			args:   []string{"--nodes", "n1"},
			wantOK: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got options
			cmd := newRootCmd("test", func(_ context.Context, o options) error { got = o; return nil })
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantOK {
				if err != nil {
					t.Fatalf("Execute() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
			var ee *exitError
			if !errors.As(err, &ee) || ee.code != exitFailure {
				t.Fatalf("exit code = %d, want %d", exitCodeOf(err), exitFailure)
			}
		})
	}
}
