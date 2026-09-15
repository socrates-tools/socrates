package configure

const (
	newline      = "\n"
	indent       = "  "
	here         = ">"
	elsewhere    = " "
	typingCursor = "_"
	quitKey      = "q"

	settingRow = "%s %-*s  %-*s  %s"
	offeredRow = "%s %s%s"

	whatToConfigure    = "What do you want to configure?"
	movingThroughAMenu = "  (up and down to move, enter to choose, q to quit)"
	movingThroughAList = "  (up and down to move, enter to choose, esc to go back)"
	typingAValue       = "  (type to edit, up and down to move, enter to confirm, esc to go back)"
	inForce            = "  (in force)"

	useTheDefault     = "use the default"
	useTheDefaultOf   = "use the default (%s)"
	resetToItsDefault = "reset %s to its default"
	leaveTheFileAlone = "quit and leave the file alone"

	nothingWasTyped   = "%s needs a value. Type one, or choose %q."
	setTo             = "%s set to %s in %s."
	nowTheDefault     = "%s now uses the default in %s."
	alreadyTheDefault = "%s is already the default, and %s is unchanged."
	readFromNowOn     = "Every command started from now on reads it, including the next round of a run already going."

	cannotBeUnderstood = "%s states something socrates-review cannot read."
	theOffendingLine   = "  line %d: %s"

	cannotBeRead    = "socrates-review configure cannot read %s: %s. Fix its permissions or move it aside and run the command again"
	cannotBeWritten = "socrates-review configure cannot write %s: %s. Fix its permissions and run the command again"

	noPersonAtAKeyboard = "socrates-review configure draws a menu and waits for a person, and nothing here is a terminal. Run %s to read what is in force, and edit the file itself to change it"
	readingCommand      = "socrates-review config"
)
