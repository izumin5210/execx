package execx_test

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/izumin5210/execx"
)

var (
	stubCmd = filepath.Join(".", "testdata", "echo", "bin", "echo")
)

func init() {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_test.go#L22-L32
	if runtime.GOOS == "windows" {
		stubCmd += ".exe"
	}
	err := exec.Command("go", "build", "-o", stubCmd, "./testdata/echo").Run()
	if err != nil {
		panic(err)
	}
}

func TestCommand(t *testing.T) {
	cases := []struct {
		test       string
		cmd        func() *execx.Cmd
		wantStdout []string
		wantStderr []string
		wantStatus *execx.ExitStatus
	}{
		{
			test:       "simple",
			cmd:        func() *execx.Cmd { return execx.Command(stubCmd, "-sleep", "1s", "It Works!") },
			wantStdout: []string{"It Works!"},
			wantStderr: []string{},
		},
		{
			test: "timeout",
			cmd: func() *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
				return execx.CommandContext(ctx, stubCmd, "-sleep", "1s", "It Works!")
			},
			wantStdout: []string{},
			wantStderr: []string{},
			wantStatus: &execx.ExitStatus{Signaled: true, Timeout: true},
		},
		{
			test: "trap timeout",
			cmd: func() *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
				return execx.CommandContext(ctx, stubCmd, "-trap", "-sleep", "1s", "It Works!")
			},
			wantStdout: []string{"It Works!"},
			wantStderr: []string{"signal received"},
		},
		{
			test: "over grace period",
			cmd: func() *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
				return execx.New(execx.WithGracePeriod(200*time.Millisecond)).
					CommandContext(ctx, stubCmd, "-trap", "-sleep", "1s", "It Works!")
			},
			wantStdout: []string{""},
			wantStderr: []string{"signal received"},
			wantStatus: &execx.ExitStatus{Signaled: true, Killed: true, Timeout: true},
		},
		{
			test: "cancel",
			cmd: func() *execx.Cmd {
				ctx, cancel := context.WithCancel(context.Background())
				go func() { time.Sleep(50 * time.Millisecond); cancel() }()
				return execx.CommandContext(ctx, stubCmd, "-sleep", "1s", "It Works!")
			},
			wantStdout: []string{},
			wantStderr: []string{},
			wantStatus: &execx.ExitStatus{Signaled: true, Canceled: true},
		},
		{
			test: "trap cancel",
			cmd: func() *execx.Cmd {
				ctx, cancel := context.WithCancel(context.Background())
				go func() { time.Sleep(50 * time.Millisecond); cancel() }()
				return execx.CommandContext(ctx, stubCmd, "-trap", "-sleep", "1s", "It Works!")
			},
			wantStdout: []string{"It Works!"},
			wantStderr: []string{"signal received"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			outW, errW := new(bytes.Buffer), new(bytes.Buffer)
			cmd := tc.cmd()
			cmd.Stdout = outW
			cmd.Stderr = errW
			err := cmd.Run()

			if err != nil {
				switch {
				case tc.wantStatus == nil:
					t.Errorf("Run() returned unknown error: %v", err)
				case tc.wantStatus != nil:
					if gotStatus, ok := err.(*execx.ExitStatus); ok {
						if got, want := gotStatus.Signaled, tc.wantStatus.Signaled; got != want {
							t.Errorf("(*ExitStatus).Signaled got %t, want %t", got, want)
						}
						if got, want := gotStatus.Killed, tc.wantStatus.Killed; got != want {
							t.Errorf("(*ExitStatus).Killed got %t, want %t", got, want)
						}
						if got, want := gotStatus.Timeout, tc.wantStatus.Timeout; got != want {
							t.Errorf("(*ExitStatus).Timeout got %t, want %t", got, want)
						}
						if got, want := gotStatus.Canceled, tc.wantStatus.Canceled; got != want {
							t.Errorf("(*ExitStatus).Canceled got %t, want %t", got, want)
						}
					} else {
						t.Errorf("Run() returned unknown error: %v", err)
					}
				}
			}

			if diff := cmp.Diff(strings.TrimSpace(outW.String()), strings.Join(tc.wantStdout, "\n")); diff != "" {
				t.Errorf("Stdout diff:\n%s", diff)
			}

			if diff := cmp.Diff(strings.TrimSpace(errW.String()), strings.Join(tc.wantStderr, "\n")); diff != "" {
				t.Errorf("Stderr diff:\n%s", diff)
			}
		})
	}
}
