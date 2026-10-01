package dispatcher

const (
	toolName = "socrates"

	failurePrefix = "socrates: "

	rootShort    = "Route to a workflow or a component: review, agora, design, or implement"
	versionShort = "Say which release this is, and which commit it was built from"

	notACommand             = "%q is not a socrates command: run socrates --help to see what socrates does"
	unusableFlag            = "%s: run socrates --help to see what socrates accepts"
	versionTakesNoArguments = "socrates version takes no arguments, and %q is one: run socrates version on its own to see which build this is"
)
