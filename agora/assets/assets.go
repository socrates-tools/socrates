package assets

import (
	"embed"
	"io/fs"
	"strings"
)

const sumFile = "bundle.sum"

//go:embed bundle bundle.sum
var embedded embed.FS

func Embedded() fs.FS {
	return embedded
}

func Files() fs.FS {
	held, err := fs.Sub(embedded, "bundle")
	if err != nil {
		panic(err)
	}
	return held
}

func Sum() string {
	held, err := fs.ReadFile(embedded, sumFile)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(held))
}
