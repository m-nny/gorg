package gorg

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/m-nny/gorg/pkg/metadata"
)

const SUBFOLDER_FORMAT = "2006/2006-01-02"
const NEW_PERM os.FileMode = 0755

func isFolderEmpty(folder string) (bool, error) {
	f, err := os.Open(folder)
	if err != nil {
		return false, err
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return false, err
	}
	return len(files) == 0, nil
}

type MetaMap map[string]metadata.MetadataList

func (m MetaMap) Add(meta *metadata.Metadata) {
	newFolder := meta.CreatedAt.Format(SUBFOLDER_FORMAT)
	m[newFolder] = m[newFolder].Add(meta)
}

func (m MetaMap) organize(destFolder string) error {
	for subFolder, metas := range m {
		newFolder := filepath.Join(destFolder, subFolder)
		if err := os.MkdirAll(newFolder, NEW_PERM); err != nil {
			return err
		}
		isEmpty, err := isFolderEmpty(newFolder)
		if err != nil {
			return err
		}
		if !isEmpty {
			return fmt.Errorf("subfolder %s is not empty", newFolder)
		}
		for i, meta := range metas {
			newFilename := fmt.Sprintf("%04d%s", i, filepath.Ext(meta.FullFilepath))
			newFilepath := filepath.Join(newFolder, newFilename)
			slog.Debug("Organize()", "newFilename", newFilename, "newFolder", newFolder, "newFilepath", newFilepath, "meta", meta)
			if _, err := meta.CloneTo(newFilepath); err != nil {
				return err
			}
		}
	}
	return nil
}

func BulkOrganize(targetFolder string, metas []*metadata.Metadata) error {
	m := MetaMap{}
	for _, meta := range metas {
		m.Add(meta)
	}
	return m.organize(targetFolder)
}
