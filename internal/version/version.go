package version

const recordedUncommitted = "yes"

var (
	version     string
	commit      string
	uncommitted string
)

func Linked(name string) Build {
	return Build{Name: name, Version: version, Commit: commit, Uncommitted: uncommitted == recordedUncommitted}
}

func Reported(name string) string {
	return Describe(Linked(name))
}
