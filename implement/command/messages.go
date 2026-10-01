package command

const toolName = "implement"

const (
	skillShort         = "Install the agent definition an implementer runs from"
	rootShort          = "A fourth workflow: build a change's tickets through implementers"
	configShort        = "Print every setting, the value in force and where that value came from"
	configureShort     = "Change a setting from a menu, one at a time, without opening the file"
	defaultBranchShort = "Remember the branch this repository's features start from"

	defaultBranchTakesOneBranch   = "socrates implement default-branch takes one branch, and was given %d: run socrates implement default-branch <branch> from inside the repository"
	defaultBranchNeedsARepository = "socrates implement default-branch needs to run inside a git repository, and %s is not inside one: run it from the repository whose features start from the branch"
	defaultBranchNeedsAWorktree   = "socrates implement default-branch needs a repository with a main worktree, and %s is inside a bare one: run it from a clone with a working tree"
	defaultBranchRemembered       = "%s now starts its features from %s, remembered in %s"

	agentsToUsage     = "the agents directory to write to, instead of the one under your home directory"
	installedTheAgent = "the %s agent definition is at %s, and the overseer can spawn it from there"

	notACommand               = "%q is not a socrates implement command: run socrates implement --help to see what socrates implement does"
	unusableFlag              = "%s: run socrates implement --help to see what socrates implement accepts"
	configureTakesNoArguments = "socrates implement configure takes no arguments, and %q is one: run socrates implement configure on its own to change a setting"
	configTakesNoArguments    = "socrates implement config takes no arguments, and %q is one: run socrates implement config on its own to read every setting"
)
