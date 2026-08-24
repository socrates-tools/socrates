package agent

const (
	notARole                    = "%q is not a role socrates has, and its roles are %v"
	notARoleAtTheCall           = "%q is not a role socrates has"
	noMessageInAnAssistantEvent = "an assistant event carries no message: %w"
	readingTheStream            = "reading the agent's output stream: %w"
	truncatedMidLine            = "the agent's output stream is truncated part-way through a line"
	truncatedMidLineAt          = "the agent's output stream is truncated part-way through a line: %s"
	notAnEventLine              = "a line of the agent's output stream is not an event: %s"
	readingWhatItWrites         = "reading what %s writes: %w"
	startingTheAgent            = "starting %s: %w"
	theAgentExited              = "%s exited: %w: %s"
	theAgentExitedSayingNothing = "%s exited: %w, and said nothing about why"
	noResultEvent               = "%s wrote a complete stream carrying no result event, so no review completed"
	theResultReported           = "%s reported the result %q, so no review completed"
	theResultReportedSaying     = "%s reported an error, so no review completed: %s"
	noStructuredOutput          = "a schema was required and %s returned no structured output, and its text is not a substitute"
	structuredOutputIsNotJSON   = "%s returned structured output that is not JSON"
	schemaNotSatisfied          = "the structured output does not satisfy the schema: %w"
	theDeadlinePassed           = "the review was given a deadline and it passed"
	silenceUntilTheDeadline     = "%s answered nothing within %s, so its process group was killed"
)
