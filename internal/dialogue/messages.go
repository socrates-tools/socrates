package dialogue

const (
	aResponseNamesItsFinding = "a response answers one finding, and this one names none: pass the identifier socrates findings printed"
	notAResolution           = "%q is not a way to answer a finding: choose one of %v"
	aRejectionIsAClaim       = "a rejection is a claim the next reviewer judges rather than a dismissal, so it needs a reason: give one saying why the finding is wrong"
	aDecisionIsAClaim        = "a decision is a claim the next reviewer judges rather than an answer that closes a finding, so it needs a reason: give one saying what was decided"
	theFindingAsksAPerson    = "finding %s asks a person rather than asking for a patch, and %q is not how a person's answer is recorded: ask them, and record what they said with socrates respond --%s and a reason"
	theFindingAsksForAPatch  = "finding %s asks for a patch rather than asking a person, and a decision is not an answer to it: choose --%s or --%s, and give a reason with a rejection"
	readingTheDocument       = "reading the responses: %w"
	writingTheDocument       = "writing the responses: %w"
)

const ResponseRules = `You did not review this code. An independent reviewer did, and it did not write any of it.
For each finding, either apply a fix and record it with socrates respond --fixed, or record a rejection with socrates respond --rejected and a concrete reason.
A rejection is a claim you are making to the next reviewer, which judges it. It is not a decision, and it does not make the finding go away.
A finding that asks a person is the person's to answer and not yours: take it to them, and record what they decided with socrates respond --decided and a reason.
Do not resolve a finding by weakening or deleting behaviour you added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause.`
