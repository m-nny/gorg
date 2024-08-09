package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/barasher/go-exiftool"
	"github.com/m-nny/gorg/pkg/gorg"
	"github.com/m-nny/gorg/pkg/measure"
)

var (
	targetFolder       = flag.String("folder", "", "target folder containing images")
	tryLoadingMetadata = flag.Bool("load_metadata", false, "try to load existing metadata json")
	cpuProfile         = flag.Bool("cpu_profile", false, "run cpu profiling")
)

func main() {
	flag.Parse()
	if *targetFolder == "" {
		log.Fatalf("should provide target folder")
	}
	if *cpuProfile {
		stopProfile := measure.StartCPUProfile()
		defer stopProfile()
	}
	measure.PrintMemUsage()

	tool, err := exiftool.NewExiftool()
	if err != nil {
		log.Fatalf("could not load exiftool: %v", err)
	}
	defer tool.Close()

	photos, err := gorg.ListAllPhotos(*targetFolder)
	if err != nil {
		log.Fatalf("could not list photos: %v", err)
	}
	log.Printf("Found %d photos", len(photos))

	if err := gorg.ReadPhotos(tool, photos, *tryLoadingMetadata); err != nil {
		log.Fatalf("could not load photos: %v", err)
	}

	measure.PrintMemUsage()

	// Force GC to clear up, should see a memory drop
	runtime.GC()
	measure.PrintMemUsage()
}
