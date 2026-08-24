package version

const (
	builtFrom          = "socrates %s"
	builtFromDirtyTree = "socrates %s, built from a working tree with uncommitted changes"
	recordedNoCommit   = "socrates recorded no commit when it was linked, so which build this is cannot be read from the binary: build it with make install"
)
