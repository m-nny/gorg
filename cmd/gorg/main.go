package main

import (
	"flag"
	"io/fs"
	"log"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/barasher/go-exiftool"
	"github.com/m-nny/gorg/pkg/metadata"
)

var (
	target_folder = flag.String("folder", "", "target folder containing images")
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
		// slog.Info("isPhotoFile", "filename", filename, "ext", ext, "photoExt", photoExt)
		if ext == photoExt {
			return true
		}
	}
	return false
}

func listAllPhotos(dirName string) ([]string, error) {
	var photos []string
	walkFunc := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// slog.Info("walkFunc", "path", path, "entry", entry, "err", err)
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

func readPhotos(tool *exiftool.Exiftool, filenames []string) error {
	meta, err := metadata.NewOrLoad(tool, filenames[0])
	if err != nil {
		return err
	}
	slog.Info("readPhotos", "meta", meta)
	if err := meta.Save(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if *target_folder == "" {
		log.Fatalf("should provide target folder")
	}
	photos, err := listAllPhotos(*target_folder)
	if err != nil {
		log.Fatalf("could not list photos: %v", err)
	}
	tool, err := exiftool.NewExiftool()
	if err != nil {
		log.Fatalf("could not load exiftool: %v", err)
	}
	defer tool.Close()

	if err := readPhotos(tool, photos); err != nil {
		log.Fatalf("could not load photo: %v", err)
	}

}
