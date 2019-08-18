package testingexecx

import (
	"context"
	"os/exec"

	"github.com/izumin5210/execx"
)

type RunFunc func(ctx context.Context, cmd *exec.Cmd) error

func NewFakeProcessFactory(f RunFunc) *FakeProcessFactory {
	return &FakeProcessFactory{
		RunFunc: f,
	}
}

var (
	_ execx.Process        = (*FakeProcess)(nil)
	_ execx.ProcessFactory = (*FakeProcessFactory)(nil)
)

type FakeProcessFactory struct {
	RunFunc RunFunc
}

func (f *FakeProcessFactory) Create(c *exec.Cmd) execx.Process {
	return &FakeProcess{RunFunc: f.RunFunc}
}

type FakeProcess struct {
	RunFunc RunFunc

	ctx    context.Context
	cancel func()
	cmd    *exec.Cmd
	errCh  chan error
}

func (p *FakeProcess) Start() error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.errCh = make(chan error)

	go func() {
		defer p.cancel()
		defer close(p.errCh)
		err := p.RunFunc(p.ctx, p.cmd)
		if err != nil {
			p.errCh <- err
		}
	}()

	return nil
}

func (p *FakeProcess) Wait() <-chan *execx.ExitStatus {
	ch := make(chan *execx.ExitStatus)
	go func() {
		defer close(ch)
		if err := <-p.errCh; err != nil {
			ex := new(execx.ExitStatus)
			ex.Code = 1
			ex.Err = err
			ch <- ex
		}
	}()
	return ch
}

func (p *FakeProcess) Terminate() error {
	p.cancel()
	return nil
}

func (p *FakeProcess) Kill() error {
	p.cancel()
	return nil
}
