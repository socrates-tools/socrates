package version

const recordedUncommitted = "yes"

var (
	commit      string
	uncommitted string
)

func Reported() string {
	return Describe(Build{Commit: commit, Uncommitted: uncommitted == recordedUncommitted})
}
