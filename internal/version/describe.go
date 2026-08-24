package version

import "fmt"

type Build struct {
	Commit      string
	Uncommitted bool
}

func Describe(build Build) string {
	if build.Commit == "" {
		return recordedNoCommit
	}
	if build.Uncommitted {
		return fmt.Sprintf(builtFromDirtyTree, build.Commit)
	}
	return fmt.Sprintf(builtFrom, build.Commit)
}
