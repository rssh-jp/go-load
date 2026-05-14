package load

import (
	"context"
	"time"

	"github.com/rssh-jp/go-load/cpu"
	"github.com/rssh-jp/go-load/memory"
)

type Load struct {
	Memory   float64
	CPU      float64
	Duration time.Duration
}

type Option func(*Load)

func OptionMemoryPercentage(m float64) Option {
	return func(l *Load) {
		l.Memory = m
	}
}
func OptionCPUPercentage(c float64) Option {
	return func(l *Load) {
		l.CPU = c
	}
}
func OptionDuration(d time.Duration) Option {
	return func(l *Load) {
		l.Duration = d
	}
}
func New(opts ...Option) *Load {
	l := new(Load)

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type Loader interface {
	Load(context.Context) error
}

func (l *Load) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chErr := make(chan error, 2)

	loaders := make([]Loader, 0, 2)
	loaders = append(loaders, memory.New(l.Memory))
	loaders = append(loaders, cpu.New(l.CPU))

	for _, loader := range loaders {
		go func(ld Loader) {
			if err := ld.Load(ctx); err != nil {
				chErr <- err
			}
		}(loader)
	}

	timer := time.NewTimer(l.Duration)
	defer timer.Stop()

	select {
	case <-timer.C:
	case err := <-chErr:
		return err
	}

	return nil
}
