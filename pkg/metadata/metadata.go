package metadata

import (
	"bytes"
	_ "crypto/md5"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/barasher/go-exiftool"
)

const METADATA_EXT = ".meta.json"

type Metadata struct {
	FullFilepath     string
	OriginalFilepath string `json:",omitempty"`
	CreatedAt        time.Time
	FileSize         int64
	FileHash         Hash
	FileHashType     string
}

func New(tool *exiftool.Exiftool, filename string) (*Metadata, error) {
	if tool == nil {
		return nil, fmt.Errorf("exiftool is nil")
	}
	fileInfo := tool.ExtractMetadata(filename)[0]
	if err := fileInfo.Err; err != nil {
		slog.Error("fileInfo: coult not get metadata", "file", fileInfo.File, "err", err)
		return nil, err
	}

	createdAt, err := parseTimeField(fileInfo.Fields["CreateDate"])
	if err != nil {
		return nil, fmt.Errorf("could not get createdAt: %w", err)
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
		FullFilepath:     filename,
		OriginalFilepath: filename,
		CreatedAt:        createdAt,
		FileSize:         fileSize,
		FileHash:         h,
		FileHashType:     HASH_TYPE.String(),
	}, nil
}

func (meta *Metadata) CloneTo(newFilepath string) (*Metadata, error) {
	newMeta := &Metadata{
		FullFilepath:     newFilepath,
		OriginalFilepath: meta.FullFilepath,
		CreatedAt:        meta.CreatedAt,
		FileSize:         meta.FileSize,
		FileHash:         meta.FileHash,
		FileHashType:     meta.FileHashType,
	}
	if err := os.Link(
		newMeta.OriginalFilepath,
		newMeta.FullFilepath); err != nil {
		return nil, err
	}
	if err := newMeta.Save(); err != nil {
		return nil, err
	}
	return newMeta, nil
}

func (meta *Metadata) MetaFilepath() string {
	return meta.FullFilepath + METADATA_EXT
}

func (meta *Metadata) Equal(other *Metadata) bool {
	if meta.CreatedAt != other.CreatedAt {
		return false
	}
	if meta.FileSize != other.FileSize {
		return false
	}
	if meta.FileHashType != other.FileHashType {
		return false
	}
	if !bytes.Equal(meta.FileHash, other.FileHash) {
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

func Load(filename string) (*Metadata, error) {
	metaFilename := filename + METADATA_EXT
	jsonString, err := os.ReadFile(metaFilename)
	if err != nil {
		return nil, err
	}
	meta := &Metadata{}
	if err := json.Unmarshal(jsonString, meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func NewOrLoad(tool *exiftool.Exiftool, filename string, tryLoading bool) (*Metadata, error) {
	if tryLoading {
		// Try loading existing
		if meta, err := Load(filename); err == nil {
			return meta, nil
		}
	}
	return New(tool, filename)
}
