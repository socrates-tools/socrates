package main

const (
	failurePrefix = "agora: "

	rootShort    = "Put a decision in front of a person and hand back what they answered"
	serveShort   = "Serve one document as one page on a loopback address, and return"
	pollShort    = "Wait for the person to send, and print what they answered"
	endShort     = "Close the exchange the page is holding"
	stopShort    = "Shut the surface down and free the port"
	versionShort = "Say which build this is, by the commit it was built from"

	documentUsage = "the document to serve, or - to read it from standard input"
	addressUsage  = "the loopback address to bind, such as 127.0.0.1:0"

	notACommand        = "%q is not an agora command: run agora --help to see what agora does"
	unusableFlag       = "%s: run agora --help to see what agora accepts"
	nameTheDocument    = "agora serve renders one document, and none was named: pass --document with a path or - to read it from standard input"
	readingTheDocument = "reading the document: %w"
	readingThisBinary  = "reading which binary agora is running as: %w"
	startingTheSurface = "starting the surface: %w"

	surfaceNeverAnswered = "the surface did not answer within %s, so nothing is being served. Run the command again"
	nobodyAnswered       = "nobody answered"
	writingTheAnswer     = "writing out what came back: %w"

	builtFrom          = "agora %s"
	builtFromDirtyTree = "agora %s, built from a working tree with uncommitted changes"
	recordedNoCommit   = "agora recorded no commit when it was linked, so which build this is cannot be read from the binary: build it with make install"
)
