package exec

import (
	"os/exec"
	"syscall"

	"github.com/Songmu/wrapcommander"
)

type Process interface {
	Start() error
	Wait() <-chan syscall.WaitStatus
	Terminate() error
	Kill() error
}

type process struct {
	cmd *exec.Cmd
}

func (p *process) Start() error {
	return p.cmd.Start()
}

func (p *process) Wait() <-chan syscall.WaitStatus {
	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout.go#L185-L191
	ch := make(chan syscall.WaitStatus)
	go func() {
		err := p.cmd.Wait()
		st, _ := wrapcommander.ErrorToWaitStatus(err)
		ch <- st
	}()
	return ch
}
