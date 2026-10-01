package configure

const (
	cannotBeUnderstood = "%s states something socrates review cannot read."
	readFromNowOn      = "Every command started from now on reads it, including the next round of a run already going."

	cannotBeRead    = "socrates review configure cannot read %s: %s. Fix its permissions or move it aside and run the command again"
	cannotBeWritten = "socrates review configure cannot write %s: %s. Fix its permissions and run the command again"

	noPersonAtAKeyboard = "socrates review configure draws a menu and waits for a person, and nothing here is a terminal. Run socrates review config to read what is in force, and edit the file itself to change it"
)
