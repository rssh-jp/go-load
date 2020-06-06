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
	ticker := time.NewTicker(l.Duration)

	ctx, cancel := context.WithCancel(context.Background())

	chErr := make(chan error, 2)

	defer close(chErr)

	loaders := make([]Loader, 0, 8)

	loaders = append(loaders, memory.New(l.Memory))
	loaders = append(loaders, cpu.New(l.CPU))

	for _, loader := range loaders {
		go func(l Loader) {
			err := l.Load(ctx)
			if err != nil {
				chErr <- err
			}
		}(loader)
	}

	select {
	case <-ticker.C:
		cancel()
	case err := <-chErr:
		return err
	}

	return nil
}
