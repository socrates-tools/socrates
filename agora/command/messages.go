package command

const toolName = "agora"

const (
	rootShort  = "Put a decision in front of a person and hand back what they answered"
	serveShort = "Serve one document as a page of its own on a loopback address, and return"
	pollShort  = "Wait for the person to send on one page and print what they answered, or print it again with --last"
	endShort   = "Close the exchange one page is holding"
	stopShort  = "Shut the surface down, taking every page with it, and free the port"

	documentUsage = "the document to serve, or - to read it from standard input"
	addressUsage  = "the loopback address to bind, such as 127.0.0.1:0"
	lastUsage     = "print again what the last poll of this session printed, without waiting, even once the surface has stopped"
	sessionUsage  = "the session to write to, as it appears at the end of the address socrates agora serve printed"

	notACommand               = "%q is not a socrates agora command: run socrates agora --help to see what socrates agora does"
	unusableFlag              = "%s: run socrates agora --help to see what socrates agora accepts"
	nameTheDocument           = "socrates agora serve renders one document, and none was named: pass --document with a path or - to read it from standard input"
	nameTheSession            = "socrates agora %s is about one page and the surface holds more than one: pass --session with the session at the end of the address socrates agora serve printed"
	noSurfaceHoldsThatSession = "no surface is open, so nothing holds the session %q. Run socrates agora serve without --session to open one"
	readingTheDocument        = "reading the document: %w"
	readingTheWorkingFolder   = "reading the folder this command was run in: %w"
	readingThisBinary         = "reading which binary socrates is running as: %w"
	startingTheSurface        = "starting the surface: %w"

	pageOpensOnThisMachineOnly = "the page is served on this machine and can be opened from it only. Where the browser that will open it is somewhere else, forward %s to where that browser is and state %s in %s, so that the address printed is one that opens"

	surfaceNeverAnswered = "the surface did not answer within %s, so nothing is being served. Run the command again"
	nobodyAnswered       = "nobody answered"
	nothingWasCollected  = "no poll has printed an answer from the session %q yet, so there is nothing to print again. Run socrates agora poll --session %s to wait for one"
	writingTheAnswer     = "writing out what came back: %w"
)
