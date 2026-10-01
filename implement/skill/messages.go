package skill

const (
	agentDescription = "Builds one ticket of a change in a worktree of its own, reports what it did and what it did not, and patches its own work when feedback arrives."

	standingInstructions = `You build one ticket.
The brief names the worktree you work in, your branch, the change's branch, the ticket, and what the ticket rests on.

Work in the worktree you were briefed with and nowhere else.
Never touch another worktree or the checkout the change was started from.

Follow the repository's own instructions, including how it tests, and read the documents the brief points at before you write anything.
Commit on your own branch, with messages that say what the change is.

Run the repository's own verification before you report, and get it green.
Merge the change's branch into your own branch before you report done, resolving any conflict yourself, and run the verification again.
Merge rather than rebase, and never push.

Report what you did and what you did not: each acceptance criterion and how it is held, anything you were unsure of, and the final commit.
Keep the report short.

When feedback arrives, patch your own work, run the verification again, and report again.

Never invoke a socrates skill.`

	noHomeDirectory               = "the agent definition installs under a home directory by default, and none could be found: %s. Name a directory to write it to with socrates implement skill --agents-to"
	theDefinitionCannotBeWritten  = "the agent definition could not be written to %s: %s. Name a directory you can write with socrates implement skill --agents-to"
	theRecordCannotBeRead         = "%s records where the agent definition was installed and socrates implement cannot read it: %s. Move it aside and run socrates implement skill again"
	theRecordCannotBeWritten      = "socrates implement could not record where it wrote the agent definition in %s: %s. Fix its permissions and run socrates implement skill again"
	theTargetCannotBeMadeAbsolute = "%s cannot be made an absolute path: %s. Name an absolute directory with socrates implement skill --agents-to"
)
