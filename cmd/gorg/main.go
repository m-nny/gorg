package main

import (
	"flag"
	"io/fs"
	"log"
	"log/slog"
	"path/filepath"
	"strings"
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
		photos = append(photos, filepath.Join(dirName, entry.Name()))
		return nil
	}
	if err := filepath.WalkDir(dirName, walkFunc); err != nil {
		return nil, err
	}
	return photos, nil
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
	log.Printf("photos: %+v", photos)
}
