package version

const (
	builtFrom                         = "socrates %s, built from %s"
	builtFromDirtyTree                = "socrates %s, built from %s, a working tree with uncommitted changes"
	builtBetweenReleases              = "socrates, a development build from %s, which is no released version"
	builtBetweenReleasesFromDirtyTree = "socrates, a development build from %s, a working tree with uncommitted changes, which is no released version"
	recordedNoCommit                  = "socrates recorded no commit when it was linked, so which build this is cannot be read from the binary: build it with make install"
)
