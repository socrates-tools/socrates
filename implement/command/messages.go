package command

const toolName = "implement"

const (
	rootShort   = "A fourth workflow: build a change's tickets through implementers"
	skillShort  = "Install the agent definition an implementer runs from"
	configShort = "Print every setting, the value in force and where that value came from"

	agentsToUsage     = "the agents directory to write to, instead of the one under your home directory"
	installedTheAgent = "the %s agent definition is at %s, and the overseer can spawn it from there"

	notACommand            = "%q is not a socrates implement command: run socrates implement --help to see what socrates implement does"
	unusableFlag           = "%s: run socrates implement --help to see what socrates implement accepts"
	configTakesNoArguments = "socrates implement config takes no arguments, and %q is one: run socrates implement config on its own to read every setting"
)
