// +build !windows

package execx

import (
	"os"
	"os/exec"
	"syscall"
)

func newOSProcess(cmd *exec.Cmd) Process {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_unix.go#L14-L19
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	return &process{
		cmd: cmd,
	}
}

func (p *process) Terminate() error {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_unix.go#L21-L35
	sig := os.Interrupt
	syssig, ok := sig.(syscall.Signal)
	if !ok {
		return p.cmd.Process.Signal(sig)
	}
	err := syscall.Kill(-p.cmd.Process.Pid, syssig)
	if err != nil {
		return err
	}
	if syssig != syscall.SIGKILL && syssig != syscall.SIGCONT {
		return syscall.Kill(-p.cmd.Process.Pid, syscall.SIGCONT)
	}
	return nil
}

func (p *process) Kill() error {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout_unix.go#L37-L39
	_ = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGKILL)
	_ = p.cmd.Process.Kill()
	return nil
}
