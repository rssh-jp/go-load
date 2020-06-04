package load

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"log"
	"time"
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

func (l *Load) Run() error {
	log.Println("START", l)
	defer log.Println("END")

	// cpu
	cpuinfos, err := cpu.Info()
	if err != nil {
		return err
	}

	for _, cpuinfo := range cpuinfos {
		log.Println(cpuinfo)
	}

	// process
	procs, err := process.Processes()
	if err != nil {
		return err
	}

	var cpupercent float64
	for _, proc := range procs {
		per, err := proc.CPUPercent()
		if err != nil {
			return err
		}

		cpupercent += per
	}

	log.Println("cpu", cpupercent, "%")

	// memory
	m, err := memoryLoad(l.Memory)
	if err != nil {
		return err
	}

	defer m.Close()

	time.Sleep(l.Duration)

	return nil
}

type memory struct {
	buf []byte
}

func (m *memory) Close() {
	m.buf = nil
}
func memoryLoad(per float64) (*memory, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	log.Println("mem", vm)
	//expectUse := uint64(float64(vm.Total) * per / 100)
	expectUse := uint64(float64(vm.Available) * per / 100)
	diff := expectUse - vm.Used

	log.Println(expectUse, diff)

	m := new(memory)
	m.buf = make([]byte, diff, diff)

	for i := 0; i < int(diff); i++ {
		m.buf[i] = 1
	}

	vm, err = mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	log.Println("mem", vm)

	return m, nil
}
