package cpu

import(
    "context"
    "log"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type Instance struct{
    per float64
}
func New(per float64)*Instance{
    return &Instance{
        per: per,
    }
}

func (inst Instance)Load(ctx context.Context, chErr chan <- error){
	// cpu
	cpuinfos, err := cpu.Info()
	if err != nil {
		chErr <- err
	}

	for _, cpuinfo := range cpuinfos {
		log.Println(cpuinfo)
	}

	// process
	procs, err := process.Processes()
	if err != nil {
		chErr <- err
	}

	var cpupercent float64
	for _, proc := range procs {
		per, err := proc.CPUPercent()
		if err != nil {
		    chErr <- err
		}

		cpupercent += per
	}

	log.Println("cpu", cpupercent, "%")

    select{
    case <-ctx.Done():
    }
}

