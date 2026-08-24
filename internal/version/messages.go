package version

const (
	builtFrom          = "scrat %s"
	builtFromDirtyTree = "scrat %s, built from a working tree with uncommitted changes"
	recordedNoCommit   = "scrat recorded no commit when it was linked, so which build this is cannot be read from the binary: build it with make install"
)
