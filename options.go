package execx

import (
	"os/exec"
	"time"
)

type Option func(*Config)

func defaultConfig() *Config {
	return &Config{
		GracePeriod:    30 * time.Second,
		ProcessFactory: ProcessFactoryFunc(newOSProcess),
	}
}

type Config struct {
	GracePeriod time.Duration

	ProcessFactory ProcessFactory
}

type ProcessFactory interface {
	Create(*exec.Cmd) Process
}

type ProcessFactoryFunc func(*exec.Cmd) Process

func (f ProcessFactoryFunc) Create(c *exec.Cmd) Process { return f(c) }

func (c *Config) apply(opts []Option) {
	for _, f := range opts {
		f(c)
	}
}

func WithGracePeriod(d time.Duration) Option {
	return func(c *Config) { c.GracePeriod = d }
}

func WithProcessFactory(f ProcessFactory) Option {
	return func(c *Config) { c.ProcessFactory = f }
}
