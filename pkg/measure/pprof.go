package measure

import (
	"log"
	"os"
	"runtime/pprof"
)

const CPU_PPROF_FILE = "out/cpu.pprof"

// StartCPUProfile starts CPU profiling using pprof and returns callback to stop profiling
func StartCPUProfile() func() {
	f, err := os.Create(CPU_PPROF_FILE)
	if err != nil {
		log.Fatalf("Could not create cpu profiling file: %v", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalf("Could not start cpu profiling: %v", err)
	}
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}
