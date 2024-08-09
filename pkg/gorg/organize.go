package gorg

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-nny/gorg/pkg/metadata"
)

const SUBFOLDER_FORMAT = "2006/2006-01-02"
const NEW_PERM os.FileMode = 0755

func NewFilename(folder string) (int, error) {
	f, err := os.Open(folder)
	if err != nil {
		return 0, err
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return 0, err
	}
	cnt := 0
	for _, file := range files {
		if strings.HasSuffix(file, metadata.METADATA_EXT) {
			continue
		}
		cnt++
	}
	return cnt, nil
}

func Organize(targetFolder string, meta *metadata.Metadata) error {
	newFolder := filepath.Join(targetFolder, meta.CreatedAt.Format(SUBFOLDER_FORMAT))
	if err := os.MkdirAll(newFolder, NEW_PERM); err != nil {
		return err
	}
	newFolderFileCnt, err := NewFilename(newFolder)
	if err != nil {
		return err
	}
	newFilename := fmt.Sprintf("%04d%s", newFolderFileCnt, filepath.Ext(meta.FullFilepath))
	newFilepath := filepath.Join(newFolder, newFilename)
	slog.Debug("Organize()", "newFilename", newFilename, "newFolder", newFolder, "newFilepath", newFilepath, "meta", meta)
	if err := os.Link(meta.FullFilepath, newFilepath); err != nil {
		return err
	}
	if err := meta.Copy(newFilepath).Save(); err != nil {
		return err
	}
	return nil
}

func BulkOrganize(targetFolder string, metas []*metadata.Metadata) error {
	for _, meta := range metas {
		if err := Organize(targetFolder, meta); err != nil {
			return err
		}
	}
	return nil
}
