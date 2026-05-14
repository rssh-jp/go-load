package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	load "github.com/rssh-jp/go-load"
)

func main() {
	cpu := flag.Float64("cpu", 50, "CPU usage percentage (0-100)")
	mem := flag.Float64("mem", 50, "Memory usage percentage (0-100)")
	duration := flag.Duration("duration", 60*time.Second, "Duration of load (e.g. 30s, 1m, 2h)")
	flag.Parse()

	if *cpu < 0 || *cpu > 100 {
		log.Fatalf("cpu must be between 0 and 100, got %.2f", *cpu)
	}
	if *mem < 0 || *mem > 100 {
		log.Fatalf("mem must be between 0 and 100, got %.2f", *mem)
	}

	fmt.Printf("Starting load: CPU=%.0f%%, Memory=%.0f%%, Duration=%s\n", *cpu, *mem, *duration)

	l := load.New(
		load.OptionCPUPercentage(*cpu),
		load.OptionMemoryPercentage(*mem),
		load.OptionDuration(*duration),
	)

	if err := l.Run(); err != nil {
		log.Fatalf("load failed: %v", err)
	}

	fmt.Println("Done.")
}
