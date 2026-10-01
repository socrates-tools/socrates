package settingsmenu

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
	noLongerStated    = "%s is no longer stated in %s."
	alreadyTheDefault = "%s is already the default, and %s is unchanged."

	theOffendingLine = "  line %d: %s"
)
