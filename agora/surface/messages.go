package surface

const (
	surfaceHasNoHomeOfItsOwn = "%s moved a home of the surface's own, and the surface has none: everything it reads and writes is under its own directory in the brand's root, and %s moves that root. Unset %s, set %s instead, and run the command again"

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

	sessionIsNotAName         = "%q is not a session: a session is the name at the end of the address agora serve printed. Pass that name to --session and run the command again"
	surfaceHoldsNoSuchSession = "the surface holds no session %q. A session comes into existence only when agora serve is given none, so run agora serve without --session to open one"
	portIsHeldByAnotherRoot   = "an agora surface of the root %s holds %s and agora will not take that port from it. Move this surface with %s in %s or with %s, and run the command again"
	surfaceStandsWithNoRecord = "an agora surface of this root stands, and what it is reached with was not written beside it. Stop that surface and run the command again"
	portIsHeldBySomethingElse = "something is already listening on %s and agora will not take that port from it. Free it, or move the surface with %s in %s or with %s, and run the command again"
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
	writingTheAnswer     = "keeping what the person answered in the session before handing it back: %w"
	readingTheKeptAnswer = "reading the answer the session kept: %w"
	openingTheSession    = "opening a session on the surface: %w"
	noSurfaceIsOpen      = "no surface is open. Run agora serve with a document and try again"

	surfaceIsBeingReplaced = "the surface is being replaced by a newer build of socrates, and takes nothing more. Run the command again"
	surfaceHoldsASession   = "the surface holds a session still being answered, and is not replaced until every session it holds is over"

	illustrationHasNoFolder  = "the illustration of decision %q names %q, a path relative to a folder nothing stated. Give an absolute path and run the command again"
	illustrationIsNotHTML    = "the illustration of decision %q names %q, which is not an HTML file. Point it at a file ending in .html and run the command again"
	illustrationDoesNotExist = "the illustration of decision %q names %q, which does not exist. Check the path and run the command again"
	illustrationIsNotAFile   = "the illustration of decision %q names %q, which is not a file. Point it at the HTML file and run the command again"
	illustrationCannotBeRead = "the illustration of decision %q names %q, which cannot be read: %v. Check the path and run the command again"
	illustrationIsTooLarge   = "the illustration of decision %q names %q, which is %d bytes, and at most %d are copied. Edit it to fit and run the command again"

	illustrationFilesAreTooLarge         = "the illustration of decision %q brings %q, and its files then total %d bytes, where at most %d are copied. Remove that file or shrink it and run the command again"
	illustrationReferenceLeavesItsFolder = "the illustration of decision %q references %q, which leaves the folder of %q. Move that file into the folder and run the command again"
	illustrationReferenceDoesNotExist    = "the illustration of decision %q references %q, which does not exist. Check the reference and run the command again"
	illustrationReferenceIsNotServed     = "the illustration of decision %q references %q, which is not an HTML, CSS, JavaScript, JSON, SVG, PNG, JPEG, WebP, GIF, WOFF or WOFF2 file. Use a file of one of those types and run the command again"
	illustrationReferenceIsNotAFile      = "the illustration of decision %q references %q, which is not a file. Point it at a file and run the command again"
	illustrationReferenceCannotBeRead    = "the illustration of decision %q references %q, which cannot be read: %v. Check the reference and run the command again"
	copyingTheIllustration               = "copying the illustration into the session: %w"
	illustrationIsNotHere                = "the session holds no such file of an illustration"

	exchangeIsOver                 = "the exchange is over, and a closed exchange does not reopen"
	answerNamesNoDecision          = "an answer names the decision %q, which this document does not hold"
	settledOtherwiseThanCollected  = "the document settles %q with a %s other than the one this session collected for it, %q. Write the settled answer exactly as agora poll handed it back, or ask the question again as a decision, and run the command again"
	answerNamesNoOption            = "an answer names option %d of decision %q, which that decision does not offer"
	singleChoiceCarriesMoreThanOne = "an answer to the single-choice decision %q carries %d choices"

	answerIsNotFormData      = "the answer was not sent as form data: %v. Send the answer as the form the page sends and run the request again"
	answerHasNoAnswerPart    = "what was sent carries no part named %q. Send the answer in that part and run the request again"
	answerPartIsRepeated     = "what was sent carries the part %q more than once. Send each part once and run the request again"
	answerIsOverTheCeiling   = "what was sent is larger than the %d bytes an answer may hold, files and all. Attach less and send again"
	answerNamesNoPart        = "the answer attaches %q, and no part of what was sent carries it. Attach the file again and send again"
	answerAttachesAPartTwice = "the answer attaches %q more than once. Attach each file in one place and send again"
	partIsNamedByNoAnswer    = "what was sent carries the file %q, and the answer does not attach it. Attach the file again and send again"
	stagingTheAttachment     = "keeping an attached file beside the session: %w"
	layingTheAttachments     = "laying the attached files in the answer's folder: %w"

	reachingTheSurface  = "reaching the surface: %w"
	surfaceAnsweredWith = "the surface answered with %d"
	surfaceRefused      = "the surface refused with %d: %s"
	readingTheAnswer    = "reading what came back: %w"
)
