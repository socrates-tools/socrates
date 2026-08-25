package surface

const (
	noHomeDirectory = "agora keeps its session under your state directory, and it could not find your home directory: %s. Set XDG_STATE_HOME and run the command again"

	addressIsNotHostAndPort    = "the address %q is not a host and a port: %s. Give an address such as 127.0.0.1:0 and run the command again"
	addressNamesNoPort         = "the address %q names no port. Give an address such as 127.0.0.1:0 and run the command again"
	addressNamesEveryInterface = "the address %q names every interface, and the surface binds a loopback address only. Give an address such as 127.0.0.1:0 and run the command again"
	addressIsNotAnAddress      = "the address %q names the host %q, and the surface binds a loopback address only. Give an address such as 127.0.0.1:0 and run the command again"
	addressIsNotLoopback       = "the address %q is not a loopback address, and the surface binds a loopback address only. Give an address such as 127.0.0.1:0 and run the command again"
	bindingTheAddress          = "binding %s: %w"

	requestNamesNoHost  = "a request naming no host cannot be the surface's own. Use the address agora serve printed and run the request again"
	requestNamesTheHost = "a request naming the host %q did not reach the surface it was bound to. Use the address agora serve printed and run the request again"

	answerCarriesNoToken = "an answer carries the token the page was served with"
	cannotStream         = "this connection cannot carry a stream, so the page cannot be told when a reply arrives"
	waitIsNotADuration   = "the wait %q is not a length of time. Give one such as 30m and run the command again"

	mintingTheToken      = "minting the token this session answers on: %w"
	writingTheSession    = "writing the session: %w"
	readingTheSession    = "reading the session: %w"
	forgettingTheSession = "forgetting the session: %w"
	writingTheDocument   = "writing the document: %w"
	noSurfaceIsOpen      = "no surface is open. Run agora serve with a document and try again"

	exchangeIsOver                 = "the exchange is over, and a closed exchange does not reopen"
	answerNamesNoDecision          = "an answer names the decision %q, which this document does not hold"
	answerNamesNoOption            = "an answer names option %d of decision %q, which that decision does not offer"
	singleChoiceCarriesMoreThanOne = "an answer to the single-choice decision %q carries %d choices"

	reachingTheSurface  = "reaching the surface: %w"
	surfaceAnsweredWith = "the surface answered with %d"
	surfaceRefused      = "the surface refused with %d: %s"
	readingTheAnswer    = "reading what came back: %w"
)
