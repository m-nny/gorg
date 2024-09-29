package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"runtime"

	"github.com/barasher/go-exiftool"
	"github.com/m-nny/gorg/pkg/gorg"
	"github.com/m-nny/gorg/pkg/logutils"
	"github.com/m-nny/gorg/pkg/measure"
)

var (
	inputFolder        = flag.String("input", "", "target folder containing images")
	destFolder         = flag.String("output", "", "destination folder")
	tryLoadingMetadata = flag.Bool("load_metadata", false, "try to load existing metadata json")
	cpuProfile         = flag.Bool("cpu_profile", false, "run cpu profiling")
	logLevel           = logutils.Level("log", slog.LevelInfo, "logging level")
	limit              = flag.Int("limit", -1, "limit number of photos to process")
	organize           = flag.Bool("organize", false, "sets if script should reorganize images")
)

func main() {
	flag.Parse()
	if *inputFolder == "" {
		log.Fatalf("should provide target folder")
	}
	if *organize && *destFolder == "" {
		log.Fatalf("should provide target folder")
	}
	slog.SetLogLoggerLevel(*logLevel)
	if *cpuProfile {
		stopProfile := measure.StartCPUProfile()
		defer stopProfile()
	}
	measure.PrintMemUsage()

	tool, err := exiftool.NewExiftool(
		exiftool.Api("IgnoreTags=All"), exiftool.Api("RequestTags=CreateDate"),
	)
	if err != nil {
		log.Fatalf("could not load exiftool: %v", err)
	}
	defer tool.Close()

	photos, err := gorg.ListAllPhotos(*inputFolder)
	if err != nil {
		log.Fatalf("could not list photos: %v", err)
	}
	slog.Info(fmt.Sprintf("Found %d photos", len(photos)))
	if 0 < *limit && *limit < len(photos) {
		photos = photos[:*limit]
		slog.Info("Limiting photos", "limit", limit)
	}

	metas, err := gorg.ReadPhotos(tool, photos, *tryLoadingMetadata)
	if err != nil {
		log.Fatalf("could not load photos: %v", err)
	}

	measure.PrintMemUsage()

	if *organize {
		if err := gorg.BulkOrganize(*destFolder, metas); err != nil {
			log.Fatalf("could not organize photos: %v", err)
		}
	}

	// Force GC to clear up, should see a memory drop
	runtime.GC()
	measure.PrintMemUsage()
}
