//go:build !agoradev

package assets

import "io/fs"

func Files() fs.FS {
	held, err := fs.Sub(embedded, "bundle")
	if err != nil {
		panic(err)
	}
	return held
}
