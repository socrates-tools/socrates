package snapshot

const (
	notAWorkingTree          = "%s is not inside a git working tree, so there is nothing to review: run scrat review from inside a repository"
	bareRepository           = "the repository at %s is bare, so it has no working tree to review: run scrat review from a working tree of it instead"
	noBaseCommit             = "the repository at %s has no commit for HEAD to resolve to, so there is no base to review against: commit something first"
	nothingToCheckOut        = "there is nothing to review, so there is no snapshot to check out"
	creatingAnIndexDirectory = "creating a directory for the snapshot's own index: %w"
	removing                 = "removing %s: %w"
	gitFailed                = "git %s in %s: %s"

	aFilesystemThatEnforcesNothing = "the filesystem holding the snapshot worktree does not enforce write permissions"
	theSealDoesNotHold             = "the snapshot worktree at %s can still be written to after every write permission in it was cleared, so either its filesystem does not enforce them or this process is root and ignores them, and socrates cannot promise that the change a reviewer read is the change that was frozen: run scrat review as a user other than root, or set SOCRATES_STATE_HOME to a directory on a filesystem that enforces permissions"
	probingTheSeal                 = "attempting a write at %s to prove the snapshot worktree is sealed: %w"
	changingThePermissionsOf       = "changing the permissions of %s: %w"
	sealingTheWorktree             = "sealing the snapshot worktree at %s: %w"
)
