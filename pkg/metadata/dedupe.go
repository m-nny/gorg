package metadata

import "log/slog"

type MetadataList []*Metadata

func (metas MetadataList) Add(newMeta *Metadata) MetadataList {
	for _, existing := range metas {
		if existing.Equal(newMeta) {
			slog.Info("MetaDataList: found duplicate", "existing", existing.FullFilepath, "newMeta", newMeta.FullFilepath)
			return metas
		}
	}
	return append(metas, newMeta)
}

func Dedupe(metas []*Metadata) []*Metadata {
	res := MetadataList{}
	for _, meta := range metas {
		res.Add(meta)
	}
	return res
}
