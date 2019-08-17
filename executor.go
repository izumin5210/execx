package execx

import (
	"context"
	"os/exec"
)

func New(opts ...Option) *Executor {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &Executor{
		Config: cfg,
	}
}

func Command(cmd string, args ...string) *Cmd {
	return New().Command(cmd, args...)
}

func CommandContext(ctx context.Context, cmd string, args ...string) *Cmd {
	return New().CommandContext(ctx, cmd, args...)
}

type Executor struct {
	Config *Config
}

func (e *Executor) Command(cmd string, args ...string) *Cmd {
	return e.CommandContext(context.Background(), cmd, args...)
}

func (e *Executor) CommandContext(ctx context.Context, cmd string, args ...string) *Cmd {
	return &Cmd{
		Cmd:    exec.Command(cmd, args...),
		Config: e.Config,
		ctx:    ctx,
	}
}
