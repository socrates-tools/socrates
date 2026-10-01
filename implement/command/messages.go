package command

const toolName = "implement"

const (
	rootShort          = "A fourth workflow: build a change's tickets through implementers"
	configShort        = "Print every setting, the value in force and where that value came from"
	defaultBranchShort = "Remember the branch this repository's features start from"

	defaultBranchTakesOneBranch   = "socrates implement default-branch takes one branch, and was given %d: run socrates implement default-branch <branch> from inside the repository"
	defaultBranchNeedsARepository = "socrates implement default-branch needs to run inside a git repository, and %s is not inside one: run it from the repository whose features start from the branch"
	defaultBranchNeedsAWorktree   = "socrates implement default-branch needs a repository with a main worktree, and %s is inside a bare one: run it from a clone with a working tree"
	repositoryHoldsNoSuchBranch   = "%s holds no local or remote-tracking branch named %q, and a default branch is one the repository holds, such as dev or origin/dev: run git branch --all to see which it holds"
	defaultBranchRemembered       = "%s now starts its features from %s, remembered in %s"

	notACommand            = "%q is not a socrates implement command: run socrates implement --help to see what socrates implement does"
	unusableFlag           = "%s: run socrates implement --help to see what socrates implement accepts"
	configTakesNoArguments = "socrates implement config takes no arguments, and %q is one: run socrates implement config on its own to read every setting"
)
