package snapshot

const (
	notAWorkingTree          = "%s is not inside a git working tree, so there is nothing to review: run scrat review from inside a repository"
	bareRepository           = "the repository at %s is bare, so it has no working tree to review: run scrat review from a working tree of it instead"
	noBaseCommit             = "the repository at %s has no commit for HEAD to resolve to, so there is no base to review against: commit something first"
	nothingToCheckOut        = "there is nothing to review, so there is no snapshot to check out"
	creatingAnIndexDirectory = "creating a directory for the snapshot's own index: %w"
	removing                 = "removing %s: %w"
	gitFailed                = "git %s in %s: %s"
)
