package version

const (
	builtFrom                         = "%s %s, built from %s"
	builtFromDirtyTree                = "%s %s, built from %s, a working tree with uncommitted changes"
	builtBetweenReleases              = "%s, a development build from %s, which is no released version"
	builtBetweenReleasesFromDirtyTree = "%s, a development build from %s, a working tree with uncommitted changes, which is no released version"
	recordedNoCommit                  = "%s recorded no commit when it was linked, so which build this is cannot be read from the binary: build it with make install"
)
