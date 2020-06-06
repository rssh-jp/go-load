package memory

import (
	"context"

	"github.com/shirou/gopsutil/mem"
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
	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	expectUse := uint64(float64(vm.Total) * inst.per / 100)
	diff := expectUse - vm.Used

	buf := make([]byte, diff, diff)

	for i := 0; i < int(diff); i++ {
		buf[i] = 1
	}

	select {
	case <-ctx.Done():
	}

	return nil
}
