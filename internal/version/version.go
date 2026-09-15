package version

const recordedUncommitted = "yes"

var (
	version     string
	commit      string
	uncommitted string
)

func Reported(name string) string {
	return Describe(Build{Name: name, Version: version, Commit: commit, Uncommitted: uncommitted == recordedUncommitted})
}
