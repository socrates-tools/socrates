package runstore

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	repositoryKeyLength = 8
	instantLayout       = "20060102T150405Z"
)

func RepositoryKey(repositoryRoot string) (string, error) {
	absolute, err := filepath.Abs(repositoryRoot)
	if err != nil {
		return "", fmt.Errorf(resolving, repositoryRoot, err)
	}

	sum := sha256.Sum256([]byte(filepath.Clean(absolute)))
	return hex.EncodeToString(sum[:])[:repositoryKeyLength], nil
}

func Name(startedAt time.Time, repositoryKey string) string {
	return startedAt.UTC().Format(instantLayout) + "-" + repositoryKey
}

func Latest(names []string, repositoryKey string) (string, bool) {
	suffix := "-" + repositoryKey

	latest := ""
	for _, name := range names {
		if !strings.HasSuffix(name, suffix) {
			continue
		}
		if name > latest {
			latest = name
		}
	}

	return latest, latest != ""
}
