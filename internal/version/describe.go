package version

import "fmt"

type Build struct {
	Name        string
	Version     string
	Commit      string
	Uncommitted bool
}

func Describe(build Build) string {
	if build.Commit == "" {
		return fmt.Sprintf(recordedNoCommit, build.Name)
	}
	if build.Version == "" {
		if build.Uncommitted {
			return fmt.Sprintf(builtBetweenReleasesFromDirtyTree, build.Name, build.Commit)
		}
		return fmt.Sprintf(builtBetweenReleases, build.Name, build.Commit)
	}
	if build.Uncommitted {
		return fmt.Sprintf(builtFromDirtyTree, build.Name, build.Version, build.Commit)
	}
	return fmt.Sprintf(builtFrom, build.Name, build.Version, build.Commit)
}
