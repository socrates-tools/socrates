package assets

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed bundle
var bundled embed.FS

//go:embed bundle.sum
var sum string

func Files() fs.FS {
	held, err := fs.Sub(bundled, "bundle")
	if err != nil {
		panic(err)
	}
	return held
}

func Sum() string {
	return strings.TrimSpace(sum)
}
