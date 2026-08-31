package version

import "fmt"

type Build struct {
	Version     string
	Commit      string
	Uncommitted bool
}

func Describe(build Build) string {
	if build.Commit == "" {
		return recordedNoCommit
	}
	if build.Version == "" {
		if build.Uncommitted {
			return fmt.Sprintf(builtBetweenReleasesFromDirtyTree, build.Commit)
		}
		return fmt.Sprintf(builtBetweenReleases, build.Commit)
	}
	if build.Uncommitted {
		return fmt.Sprintf(builtFromDirtyTree, build.Version, build.Commit)
	}
	return fmt.Sprintf(builtFrom, build.Version, build.Commit)
}
