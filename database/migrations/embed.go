package migrations

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed *.sql
var rawFS embed.FS

// FS is the embedded SQL migration filesystem with macOS AppleDouble files filtered out.
var FS fs.FS = filteredFS{rawFS}

type filteredFS struct {
	base embed.FS
}

func (f filteredFS) Open(name string) (fs.File, error) {
	return f.base.Open(name)
}

func (f filteredFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := f.base.ReadDir(name)
	if err != nil {
		return nil, err
	}

	out := make([]fs.DirEntry, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "_") {
			out = append(out, entry)
		}
	}
	return out, nil
}
