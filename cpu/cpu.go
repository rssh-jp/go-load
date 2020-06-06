package cpu

import (
	"context"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type Instance struct {
	per float64
}

func New(per float64) *Instance {
	return &Instance{
		per: per,
	}
}

func (inst Instance) Load(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ticker.C:
			err := inst.useCPU()
			if err != nil {
				return err
			}
		case <-ctx.Done():
			break LOOP
		}
	}

	return nil
}

func (inst Instance) useCPU() error {
	// cpu
	cpuinfos, err := cpu.Times(true)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < len(cpuinfos); i++ {
		wg.Add(1)
		go func(cpuindex int) {
			defer wg.Done()
			load(cpuindex, inst.per)
		}(i)
	}
	wg.Wait()

	return nil
}

func load(cpuindex int, per float64) {
	cpuinfos, err := cpu.Times(true)
	if err != nil {
		return
	}

	if cpuindex >= len(cpuinfos) {
		return
	}

	cpuinfo := cpuinfos[cpuindex]

	expectUse := (cpuinfo.Total() * per) / 100
	diff := expectUse - (cpuinfo.Total() - cpuinfo.Idle)

	diffbyte := int(diff * 1000 * 1000)

	// どれくらい負荷をかけるかの目安
	count := int(diffbyte * 72 / 100)

	work := 100
	for i := 0; i < count; i++ {
		work /= 2
	}
}
