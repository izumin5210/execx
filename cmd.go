package execx

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/Songmu/wrapcommander"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
	ErrNotStarted    = errors.New("command not started")
)

type Cmd struct {
	*exec.Cmd
	*Config
	ctx context.Context
	p   Process
}

func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *Cmd) Start() error {
	c.p = c.CreateProcessFunc(c.Cmd)
	if err := c.p.Start(); err != nil {
		return &ExitStatus{
			Code: wrapcommander.ResolveExitCode(err),
			Err:  err,
		}
	}
	return nil
}

func (c *Cmd) Wait() error {
	if c.p == nil {
		return ErrNotStarted
	}

	// https://github.com/Songmu/timeout/blob/v0.4.0/timeout.go#L132-L174
	ex := &ExitStatus{}

	killCh := make(chan struct{}, 2)

	done := make(chan struct{})
	defer close(done)

	exitCh := c.p.Wait()

	for {
		select {
		case st := <-exitCh:
			ex.Code = wrapcommander.WaitStatusToExitCode(st)
			ex.Signaled = st.Signaled()

			if ex.Code == wrapcommander.ExitNormal {
				return nil
			}

			return ex

		case <-killCh:
			c.p.Kill()
			ex.Killed = true

		case <-c.ctx.Done():
			c.p.Terminate()
			ex.Err = c.ctx.Err()

			go func() {
				select {
				case <-done:
					return
				case <-time.After(c.GracePeriod):
					killCh <- struct{}{}
				}
			}()
		}
	}
}

func (c *Cmd) CombinedOutput() ([]byte, error) {
	buf := new(bytes.Buffer)
	c.Stdout = buf
	c.Stderr = buf
	err := c.Run()
	return buf.Bytes(), err
}

func (c *Cmd) Output() ([]byte, error) {
	buf := new(bytes.Buffer)
	c.Stdout = buf
	err := c.Run()
	return buf.Bytes(), err
}
