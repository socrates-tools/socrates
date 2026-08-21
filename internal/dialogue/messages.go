package dialogue

const (
	aResponseNamesItsFinding = "a response answers one finding, and this one names none: pass the identifier scrat findings printed"
	notAResolution           = "%q is not a way to answer a finding: choose one of %v"
	aRejectionIsAClaim       = "a rejection is a claim the next reviewer judges rather than a dismissal, so it needs a reason: give one saying why the finding is wrong"
	readingTheDocument       = "reading the responses: %w"
	writingTheDocument       = "writing the responses: %w"
)

const ResponseRules = `You did not review this code. An independent reviewer did, and it did not write any of it.
For each finding, either apply a fix and record it with scrat respond --fixed, or record a rejection with scrat respond --rejected and a concrete reason.
A rejection is a claim you are making to the next reviewer, which judges it. It is not a decision, and it does not make the finding go away.
Do not resolve a finding by weakening or deleting behaviour you added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause.`
