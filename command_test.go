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

	"github.com/izumin5210/execx"
)

var (
	stubCmd = filepath.Join(".", "testdata", "echo", "bin", "echo")
	isWin   = runtime.GOOS == "windows"
)

func init() {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_test.go#L22-L32
	if isWin {
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
			cmd:        func() *execx.Cmd { return execx.Command(stubCmd, "-sleep", "3s", "It Works!") },
			wantStdout: []string{"It Works!"},
			wantStderr: []string{},
		},
		{
			test: "timeout",
			cmd: func() *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
				return execx.CommandContext(ctx, stubCmd, "-sleep", "3s", "It Works!")
			},
			wantStdout: []string{},
			wantStderr: []string{},
			wantStatus: &execx.ExitStatus{Signaled: true, Timeout: true},
		},
		{
			test: "trap timeout",
			cmd: func() *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
				return execx.CommandContext(ctx, stubCmd, "-trap", "-sleep", "3s", "It Works!")
			},
			wantStdout: []string{"It Works!"},
			wantStderr: []string{"signal received"},
		},
		{
			test: "over grace period",
			cmd: func() *execx.Cmd {
				ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
				return execx.New(execx.WithGracePeriod(1500*time.Millisecond)).
					CommandContext(ctx, stubCmd, "-trap", "-sleep", "3s", "It Works!")
			},
			wantStdout: []string{""},
			wantStderr: []string{"signal received"},
			wantStatus: &execx.ExitStatus{Signaled: true, Killed: true, Timeout: true},
		},
		{
			test: "cancel",
			cmd: func() *execx.Cmd {
				ctx, cancel := context.WithCancel(context.Background())
				go func() { time.Sleep(100 * time.Millisecond); cancel() }()
				return execx.CommandContext(ctx, stubCmd, "-sleep", "3s", "It Works!")
			},
			wantStdout: []string{},
			wantStderr: []string{},
			wantStatus: &execx.ExitStatus{Signaled: true, Canceled: true},
		},
		{
			test: "trap cancel",
			cmd: func() *execx.Cmd {
				ctx, cancel := context.WithCancel(context.Background())
				go func() { time.Sleep(100 * time.Millisecond); cancel() }()
				return execx.CommandContext(ctx, stubCmd, "-trap", "-sleep", "3s", "It Works!")
			},
			wantStdout: []string{"It Works!"},
			wantStderr: []string{"signal received"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			defer func(l execx.Logger) { execx.DefaultErrorLog = l }(execx.DefaultErrorLog)
			execx.DefaultErrorLog = &testLogger{t: t}

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
						if isWin {
							t.Log("(*execx.ExitStatus).Signaled does not work on windows")
						} else {
							if got, want := gotStatus.Signaled, tc.wantStatus.Signaled; got != want {
								t.Errorf("(*ExitStatus).Signaled got %t, want %t", got, want)
							}
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

			if got, want := strings.TrimSpace(outW.String()), strings.Join(tc.wantStdout, "\n"); got != want {
				t.Errorf("Stdout was:\n%s\nwant:\n%s", got, want)
			}

			if got, want := strings.TrimSpace(errW.String()), strings.Join(tc.wantStderr, "\n"); got != want {
				t.Errorf("Stderr was:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

type testLogger struct {
	t *testing.T
}

func (l *testLogger) Print(args ...interface{}) {
	l.t.Helper()
	l.t.Log(args...)
}
