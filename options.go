package execx

import (
	"os/exec"
	"time"
)

var (
	DefaultGracePeriod time.Duration = 30 * time.Second
	DefaultErrorLog    Logger        = new(nopLogger)
)

type Option func(*Config)

func defaultConfig() *Config {
	return &Config{
		GracePeriod:    DefaultGracePeriod,
		ProcessFactory: ProcessFactoryFunc(newOSProcess),
		ErrorLog:       DefaultErrorLog,
	}
}

type Config struct {
	GracePeriod time.Duration

	ProcessFactory ProcessFactory

	ErrorLog Logger
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

func WithErrorLog(l Logger) Option {
	return func(c *Config) { c.ErrorLog = l }
}
