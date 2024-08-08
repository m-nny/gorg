package metadata

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/barasher/go-exiftool"
)

type Metadata struct {
	FullFilepath string
	CreatedAt    time.Time
}

func Read(tool *exiftool.Exiftool, filename string) (*Metadata, error) {
	fileInfo := tool.ExtractMetadata(filename)[0]
	if err := fileInfo.Err; err != nil {
		slog.Error("fileInfo: coult not get metadata", "file", fileInfo.File, "err", err)
		return nil, err
	}
	createdAt, err := parseTimeField(fileInfo.Fields["CreateDate"])
	if err != nil {
		return nil, fmt.Errorf("could not get created at: %w", err)
	}
	slog.Info("Metadata", "createdAt", createdAt)
	return &Metadata{
		FullFilepath: filename,
		CreatedAt:    createdAt,
	}, nil
}

var datetimeFormat = "2006:01:02 15:04:05-07:00"

func parseTimeField(value any) (time.Time, error) {
	val, ok := value.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("value is not string, but %T", value)
	}
	return time.Parse(datetimeFormat, val)
}
