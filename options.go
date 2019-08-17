package exec

import (
	"os/exec"
	"time"
)

type Option func(*Config)

func defaultConfig() *Config {
	return &Config{
		TerminateTimeout:  30 * time.Second,
		CreateProcessFunc: newOSProcess,
	}
}

type Config struct {
	TerminateTimeout time.Duration

	CreateProcessFunc CreateProcessFunc
}

func (c *Config) apply(opts []Option) {
	for _, f := range opts {
		f(c)
	}
}

type CreateProcessFunc func(*exec.Cmd) Process

func WithTerminateTimeout(d time.Duration) Option {
	return func(c *Config) { c.TerminateTimeout = d }
}

func WithCreateProcessFunc(f CreateProcessFunc) Option {
	return func(c *Config) { c.CreateProcessFunc = f }
}
