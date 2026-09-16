package snapshot

const (
	notAWorkingTree          = "%s is not inside a git working tree, so there is nothing to review: run socrates review begin from inside a repository"
	bareRepository           = "the repository at %s is bare, so it has no working tree to review: run socrates review begin from a working tree of it instead"
	noBaseCommit             = "the repository at %s has no commit for HEAD to resolve to, so there is no base to review against: commit something first"
	nothingToCheckOut        = "there is nothing to review, so there is no snapshot to check out"
	noSuchBase               = "%s names no commit in the repository at %s, so there is nothing to review a change against: check how it is spelled, and name a branch or a commit that is in this repository"
	noCommonAncestor         = "%s and HEAD share no commit in the repository at %s, so a diff between them would be a diff of two whole trees rather than of a change: name a base this branch grew out of"
	baseHasGone              = "the commit %s this run was started against is no longer in the repository at %s, so a later round would be reviewed against something the first round never saw"
	creatingAnIndexDirectory = "creating a directory for the snapshot's own index: %w"
	removing                 = "removing %s: %w"
	gitFailed                = "git %s in %s: %s"
	unreadableTreeEntry      = "the tree entry for %s at %s does not say how large it is: %s"

	aFilesystemThatEnforcesNothing = "the filesystem holding the snapshot worktree does not enforce write permissions"
	theSealDoesNotHold             = "the snapshot worktree at %s can still be written to after every write permission in it was cleared, so either its filesystem does not enforce them or this process is root and ignores them, and socrates review cannot promise that the change a reviewer read is the change that was frozen: run socrates review begin as a user other than root, or set SOCRATES_HOME to a directory on a filesystem that enforces permissions"
	probingTheSeal                 = "attempting a write at %s to prove the snapshot worktree is sealed: %w"
	changingThePermissionsOf       = "changing the permissions of %s: %w"
	sealingTheWorktree             = "sealing the snapshot worktree at %s: %w"
)
