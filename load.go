package load

import (
    "context"
	"log"
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

type Loader interface{
    Load(context.Context, chan <- error)
}

type ErrLoader struct{
    loader Loader
    err chan error
}

func (l *Load) Run() error {
	log.Println("START", l)
	defer log.Println("END")

    ticker := time.NewTicker(l.Duration)

    ctx, cancel := context.WithCancel(context.Background())

    defer cancel()

    chErrMemory := make(chan error)
    chErrCPU := make(chan error)

    defer close(chErrMemory)
    defer close(chErrCPU)

    errLoader := make([]ErrLoader, 0, 8)

    errLoader = append(errLoader, ErrLoader{memory.New(l.Memory), make(chan error)})
    errLoader = append(errLoader, ErrLoader{cpu.New(l.CPU), make(chan error)})


	// memory
    m := memory.New(l.Memory)
	go m.Load(ctx, chErrMemory)

    // cpu
    c := cpu.New(l.CPU)
    go c.Load(ctx, chErrCPU)

    select{
    case <-ticker.C:
    case err := <-chErrMemory:
        return err
    case err := <-chErrCPU:
        return err
    }

	return nil
}

