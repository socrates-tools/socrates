//go:build agoradev

package assets

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

func Files() fs.FS {
	_, source, _, known := runtime.Caller(0)
	if !known {
		panic("the development build cannot locate the bundle on disk")
	}
	return os.DirFS(filepath.Join(filepath.Dir(source), "bundle"))
}
