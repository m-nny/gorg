package metadata

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/barasher/go-exiftool"
)

type Metadata struct {
	FullFilepath string
	CreatedAt    time.Time
	UniqueID     string
	FileSize     int64
}

func Read(tool *exiftool.Exiftool, filename string) (*Metadata, error) {
	fileInfo := tool.ExtractMetadata(filename)[0]
	if err := fileInfo.Err; err != nil {
		slog.Error("fileInfo: coult not get metadata", "file", fileInfo.File, "err", err)
		return nil, err
	}

	// for k, v := range fileInfo.Fields {
	// 	slog.Info("fileInfo", "file", fileInfo.File, "key", k, "value", v, "value.type", fmt.Sprintf("%T", v))
	// }

	createdAt, err := parseTimeField(fileInfo.Fields["CreateDate"])
	if err != nil {
		return nil, fmt.Errorf("could not get createdAt: %w", err)
	}

	uniqueID, err := fileInfo.GetString("ImageUniqueID")
	if err != nil {
		return nil, fmt.Errorf("could not get uniqueID: %w", err)
	}

	fileSize, err := getFileSize(filename)
	if err != nil {
		return nil, fmt.Errorf("could not get fizeSize: %w", err)
	}

	return &Metadata{
		FullFilepath: filename,
		CreatedAt:    createdAt,
		UniqueID:     uniqueID,
		FileSize:     fileSize,
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

func getFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
