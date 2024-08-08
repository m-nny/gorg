package gorg

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/barasher/go-exiftool"
	"github.com/m-nny/gorg/pkg/measure"
	"github.com/m-nny/gorg/pkg/metadata"
)

var (
	photoExtentions = []string{"CR2", "CR3", "JPEG", "JPG"}
)

func isPhotoFile(filename string) bool {
	// TODO: actualy check by content, not just by extention of filename
	ext := strings.ToUpper(filepath.Ext(filename))
	if ext == "" || ext == "." {
		return false
	}
	if ext[0] != '.' {
		slog.Error("filepath.Ext() returned something not starting with .",
			"filename", filename,
			"ext", ext)
	}
	ext = ext[1:]

	for _, photoExt := range photoExtentions {
		if ext == photoExt {
			return true
		}
	}
	return false
}

func ListAllPhotos(dirName string) ([]string, error) {
	defer measure.TimerStop(measure.Timer("ListAllPhotos"))
	var photos []string
	walkFunc := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !isPhotoFile(entry.Name()) {
			return nil
		}
		photos = append(photos, path)
		return nil
	}
	// TODO:  WalkDir() reads entire directory into memory before proceeding to wald the dir. We probably want to yield these as we go
	if err := filepath.WalkDir(dirName, walkFunc); err != nil {
		return nil, err
	}
	return photos, nil
}

func ReadPhotos(tool *exiftool.Exiftool, filenames []string, tryLoading bool) error {
	defer measure.TimerStop(measure.Timer("ReadPhotos"))
	for _, filename := range filenames {
		meta, err := metadata.NewOrLoad(tool, filename, tryLoading)
		if err != nil {
			return fmt.Errorf("could not load metadata for file %q: %w", filename, err)
		}
		slog.Debug("readPhotos", "meta", meta)
		if err := meta.Save(); err != nil {
			return err
		}
	}
	return nil
}
