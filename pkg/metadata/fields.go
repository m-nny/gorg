package metadata

import (
	"crypto"
	"fmt"
	"io"
	"os"
	"time"
)

var datetimeFormats = []string{
	"2006:01:02 15:04:05-07:00",
	"2006:01:02 15:04:05",
}

func parseTimeField(value any) (t time.Time, err error) {
	val, ok := value.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("value is not string, but %T", value)
	}
	for _, format := range datetimeFormats {
		t, err = time.Parse(format, val)
		if err == nil {
			return t, nil
		}
	}
	return t, fmt.Errorf("no format matched, last error: %w", err)
}

func getFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

type Hash []byte

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
