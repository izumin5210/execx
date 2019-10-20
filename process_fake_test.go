package execx_test

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/izumin5210/execx"
)

func TestFakeProcess(t *testing.T) {
	cases := []struct {
		test       string
		run        func(ctx context.Context, cmd *exec.Cmd) error
		cmd        func(exec *execx.Executor) *execx.Cmd
		wantStdout []string
		wantStderr []string
		wantError  error
	}{
		{
			test: "simple",
			run: func(ctx context.Context, cmd *exec.Cmd) error {
				_, err := fmt.Fprintln(cmd.Stdout, "2")
				return err
			},
			cmd: func(exec *execx.Executor) *execx.Cmd {
				return exec.Command("echo", "1")
			},
			wantStdout: []string{"2"},
			wantStderr: []string{},
			wantError:  nil,
		},
		{
			test: "cancel",
			run: func(ctx context.Context, cmd *exec.Cmd) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(1 * time.Second):
					fmt.Fprintln(cmd.Stdout, "2")
				}
				return nil
			},
			cmd: func(exec *execx.Executor) *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
				return exec.CommandContext(ctx, "echo", "1")
			},
			wantStdout: []string{},
			wantStderr: []string{},
			wantError:  context.Canceled,
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			cmd := tc.cmd(execx.New(
				execx.WithFakeProcess(tc.run),
			))

			var outW, errW bytes.Buffer
			cmd.Stdout = &outW
			cmd.Stderr = &errW
			err := cmd.Run()

			switch {
			case err == nil && tc.wantError != nil:
				t.Errorf("Run() returned nil, want %v", tc.wantError)
			case err != nil && tc.wantError == nil:
				t.Errorf("Run() returned %v, want nil", err)
			case err != nil && tc.wantError != nil:
				if ex, ok := err.(*execx.ExitStatus); ok {
					if got, want := ex.Err, tc.wantError; got != want {
						t.Errorf("Run() returned %v, want %v", got, want)
					}
				} else {
					t.Errorf("Run() returned an unknown error %v, want nil", err)
				}
			}

			if got, want := strings.TrimSpace(outW.String()), strings.Join(tc.wantStdout, "\n"); got != want {
				t.Errorf("stdout received output %q, want %q", got, want)
			}

			if got, want := strings.TrimSpace(errW.String()), strings.Join(tc.wantStderr, "\n"); got != want {
				t.Errorf("stderr received errput %q, want %q", got, want)
			}
		})
	}
}
