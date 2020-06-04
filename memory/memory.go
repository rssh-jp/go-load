package memory

import(
    "context"
    "log"

	"github.com/shirou/gopsutil/mem"
)

type Instance struct{
    per float64
}
func New(per float64)*Instance{
    return &Instance{
        per: per,
    }
}

func (inst Instance)Load(ctx context.Context, chErr chan <- error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		chErr <- err
	}

	log.Println("mem", vm)
	expectUse := uint64(float64(vm.Total) * inst.per / 100)
	diff := expectUse - vm.Used

	log.Println(expectUse, diff)

	buf := make([]byte, diff, diff)

	for i := 0; i < int(diff); i++ {
		buf[i] = 1
	}

	vm, err = mem.VirtualMemory()
	if err != nil {
		chErr <- err
	}

	log.Println("mem", vm)

    select{
    case <-ctx.Done():
    }
}

