package version

const recordedUncommitted = "yes"

var (
	version     string
	commit      string
	uncommitted string
)

func Reported() string {
	return Describe(Build{Version: version, Commit: commit, Uncommitted: uncommitted == recordedUncommitted})
}
