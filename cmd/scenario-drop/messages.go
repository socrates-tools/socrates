package main

const (
	failurePrefix = "scenario-drop: "

	dropHeading      = "%s"
	dropLine         = "  %s: %s"
	scenarioLine     = "    the promoted specification would stop holding: %s"
	nothingUnderWay  = "no change is under way, so archiving one would drop nothing: the rule itself is held by the tests over a fabricated tree"
	nothingDropped   = "%d requirements would be substituted and %d added, and the promoted specifications would lose nothing"
	somethingDropped = "%d of the requirements a change would promote lose something, out of %d substituted and %d added"
)
