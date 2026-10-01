package config

const (
	unreadableFile          = "%s exists and socrates implement cannot read it: %s. Fix its permissions or move it aside and run the command again"
	unparsableFile          = "%s is not configuration socrates implement can read: %s. Fix the file or move it aside, because socrates implement will not quietly run on defaults while a configuration file says otherwise"
	statesSomethingUnusable = "%s states something socrates implement cannot use: %s. Fix the file or move it aside and run the command again"
	concurrencyNotPositive  = "concurrency is %d, and it is a positive whole number"
	notAWholeNumber         = "line %d: %q is not a whole number, and concurrency is a positive whole number"
	modelHasWhitespace      = "roles.implementer.model is %q, and it is a Claude Code model alias or a full model identifier, which holds no whitespace"
	unknownEffort           = "roles.implementer.effort is %q, and the efforts are %v"
	repositoryNotAbsolute   = "default_branches names %q, and a repository is named by the absolute path of its top level"
	branchNotUsable         = "default_branches holds %q for %s, and it is the name of one branch, with no whitespace"
	notARepository          = "this directory is not inside a git repository"
	bareRepository          = "this is a bare repository, which has no main worktree to name it by"
	gitCannotRun            = "socrates implement could not run git: %v. Check that git is on PATH and run the command again"
	savingCommand           = "socrates implement default-branch"

	settingLine = "%-*s  %-*s  %s"
)
