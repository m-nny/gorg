package metadata

import (
	"bytes"
	"crypto"
	_ "crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/barasher/go-exiftool"
)

type Hash []byte

const METADATA_EXT = ".meta.json"

type Metadata struct {
	FullFilepath string
	CreatedAt    time.Time
	UniqueID     string
	FileSize     int64
	FileHash     Hash
	FileHashType string
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

	h, err := getFileHash(filename)
	if err != nil {
		return nil, fmt.Errorf("could not get hash: %w", err)
	}

	return &Metadata{
		FullFilepath: filename,
		CreatedAt:    createdAt,
		UniqueID:     uniqueID,
		FileSize:     fileSize,
		FileHash:     h,
		FileHashType: HASH_TYPE.String(),
	}, nil
}

func (meta *Metadata) Equal(other *Metadata) bool {
	if meta.CreatedAt != other.CreatedAt {
		return false
	}
	if meta.UniqueID != other.UniqueID {
		return false
	}
	if meta.FileSize != other.FileSize {
		return false
	}
	if meta.FileHashType != other.FileHashType {
		return false
	}
	if bytes.Equal(meta.FileHash, other.FileHash) {
		return false
	}
	return true
}

func (meta *Metadata) Save() error {
	jsonString, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	metaFilename := meta.FullFilepath + METADATA_EXT
	return os.WriteFile(metaFilename, jsonString, os.ModePerm)
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

const HASH_TYPE = crypto.MD5

func getFileHash(filename string) (Hash, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := HASH_TYPE.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
