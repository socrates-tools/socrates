package surface

const (
	noHomeDirectory      = "agora keeps everything it reads and writes under one home directory, and it could not find your own home directory to put it under: %s. Set %s to an absolute path and run the command again"
	statedHomeIsRelative = "%s names %q, which is not an absolute path, and agora would keep its home wherever the command happened to be run. Set it to an absolute path and run the command again"

	unreadableConfiguration = "%s exists and agora cannot read it: %s. Fix its permissions or move it aside and run the command again"
	unparsableConfiguration = "%s is not configuration agora can read: %s. Fix the file or move it aside, because agora will not quietly run on defaults while a configuration file says otherwise"

	addressIsNotHostAndPort    = "the address %q is not a host and a port: %s. Give an address such as 127.0.0.1:0 and run the command again"
	addressNamesNoPort         = "the address %q names no port. Give an address such as 127.0.0.1:0 and run the command again"
	addressNamesEveryInterface = "the address %q names every interface, and the surface binds a loopback address only. Give an address such as 127.0.0.1:0 and run the command again"
	addressIsNotAnAddress      = "the address %q names the host %q, and the surface binds a loopback address only. Give an address such as 127.0.0.1:0 and run the command again"
	addressIsNotLoopback       = "the address %q is not a loopback address, and the surface binds a loopback address only. Give an address such as 127.0.0.1:0 and run the command again"
	bindingTheAddress          = "binding %s: %w"

	addressIsNotAURL     = "the address %q the surface is reached at is not an address a person can open: %s. Give one such as https://agora.example:8443/ and run the command again"
	addressIsNotOpenable = "the address %q the surface is reached at names no http or https scheme, and it is handed to a person to open. Give one such as https://agora.example:8443/ and run the command again"
	addressNamesNoHost   = "the address %q the surface is reached at names no host, and it is handed to a person to open. Give one such as https://agora.example:8443/ and run the command again"

	requestNamesNoHost  = "a request naming no host cannot be the surface's own. Use the address agora serve printed and run the request again"
	requestNamesTheHost = "a request naming the host %q did not reach the surface it was bound to. Use the address agora serve printed and run the request again"

	answerCarriesNoToken    = "an answer carries the token the page was served with"
	requestCarriesNoControl = "this is asked of the surface itself, and it is asked by whoever started it"

	surfaceHoldsNoSuchSession = "the surface holds no session %q. A session comes into existence only when agora serve is given none, so run agora serve without --session to open one"
	portIsHeldByAnother       = "something is already listening on %s and it is not a surface agora knows of, so agora will not take that port from it. Free it, or name another loopback address with --address, and run the command again"
	bodyIsOverTheCeiling      = "what was sent is larger than the %d bytes the surface accepts. Send a smaller document"
	cannotStream              = "this connection cannot carry a stream, so the page cannot be told when a reply arrives"
	waitIsNotADuration        = "the wait %q is not a length of time. Give one such as 30m and run the command again"

	mintingTheToken      = "minting the token this session answers on: %w"
	mintingTheSession    = "minting the session the page belongs to: %w"
	writingTheRecord     = "writing what the surface is reached at: %w"
	readingTheRecord     = "reading what the surface is reached at: %w"
	forgettingTheRecord  = "forgetting what the surface was reached at: %w"
	writingTheSession    = "writing the session: %w"
	forgettingTheSession = "forgetting the session: %w"
	writingTheDocument   = "writing the document: %w"
	openingTheSession    = "opening a session on the surface: %w"
	noSurfaceIsOpen      = "no surface is open. Run agora serve with a document and try again"

	exchangeIsOver                 = "the exchange is over, and a closed exchange does not reopen"
	answerNamesNoDecision          = "an answer names the decision %q, which this document does not hold"
	answerNamesNoOption            = "an answer names option %d of decision %q, which that decision does not offer"
	singleChoiceCarriesMoreThanOne = "an answer to the single-choice decision %q carries %d choices"

	reachingTheSurface     = "reaching the surface: %w"
	addressIsNotTheSurface = "something answers at %s and it does not say it is a surface agora started"
	surfaceAnsweredWith    = "the surface answered with %d"
	surfaceRefused         = "the surface refused with %d: %s"
	readingTheAnswer       = "reading what came back: %w"
)
