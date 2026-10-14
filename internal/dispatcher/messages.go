package dispatcher

const (
	toolName = "socrates"

	failurePrefix = "socrates: "

	rootShort    = "Route to a workflow or a component: review, agora, design, or implement"
	versionShort = "Say which release this is, and which commit it was built from"
	skillShort   = "Install every skill every namespace offers, where each namespace's own skill puts it"

	notACommand             = "%q is not a socrates command: run socrates --help to see what socrates does"
	unusableFlag            = "%s: run socrates --help to see what socrates accepts"
	skillTakesNoArguments   = "socrates skill takes no arguments, and %q is one: run socrates skill on its own to install every skill, or socrates <namespace> skill --to to put one namespace's elsewhere"
	versionTakesNoArguments = "socrates version takes no arguments, and %q is one: run socrates version on its own to see which build this is"
)
