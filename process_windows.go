package execx

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func newOSProcess(cmd *exec.Cmd) Process {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_windows.go#L9-L16
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT | 0x00000200,
		}
	}
	return &process{
		cmd: cmd,
	}
}

func (p *process) Terminate() error {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_windows.go#L18-L20
	return p.cmd.Process.Signal(p.Signal())
}

func (p *process) Kill() error {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_windows.go#L22-L24
	return exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(p.cmd.Process.Pid)).Run()
}

func (p *process) Signal() os.Signal { return os.Interrupt }
