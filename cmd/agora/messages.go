package main

const (
	failurePrefix = "agora: "

	rootShort    = "Put a decision in front of a person and hand back what they answered"
	serveShort   = "Serve one document as a page of its own on a loopback address, and return"
	pollShort    = "Wait for the person to send on one page, and print what they answered"
	endShort     = "Close the exchange one page is holding"
	stopShort    = "Shut the surface down, taking every page with it, and free the port"
	versionShort = "Say which release this is, and which commit it was built from"

	documentUsage = "the document to serve, or - to read it from standard input"
	addressUsage  = "the loopback address to bind, such as 127.0.0.1:0"
	sessionUsage  = "the session to write to, as it appears at the end of the address agora serve printed"

	notACommand               = "%q is not an agora command: run agora --help to see what agora does"
	unusableFlag              = "%s: run agora --help to see what agora accepts"
	nameTheDocument           = "agora serve renders one document, and none was named: pass --document with a path or - to read it from standard input"
	nameTheSession            = "agora %s is about one page and the surface holds more than one: pass --session with the session at the end of the address agora serve printed"
	noSurfaceHoldsThatSession = "no surface is open, so nothing holds the session %q. Run agora serve without --session to open one"
	readingTheDocument        = "reading the document: %w"
	readingThisBinary         = "reading which binary agora is running as: %w"
	startingTheSurface        = "starting the surface: %w"

	pageOpensOnThisMachineOnly = "the page is served on this machine and can be opened from it only. Where the browser that will open it is somewhere else, forward %s to where that browser is and state %s in %s, so that the address printed is one that opens"

	surfaceNeverAnswered = "the surface did not answer within %s, so nothing is being served. Run the command again"
	nobodyAnswered       = "nobody answered"
	writingTheAnswer     = "writing out what came back: %w"

	builtFrom                         = "agora %s, built from %s"
	builtFromDirtyTree                = "agora %s, built from %s, a working tree with uncommitted changes"
	builtBetweenReleases              = "agora, a development build from %s, which is no released version"
	builtBetweenReleasesFromDirtyTree = "agora, a development build from %s, a working tree with uncommitted changes, which is no released version"
	recordedNoCommit                  = "agora recorded no commit when it was linked, so which build this is cannot be read from the binary: build it with make install"
)
