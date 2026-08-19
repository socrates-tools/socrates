package fidelity

const (
	notAnEvent          = "a line of the recorded stream is not an event: %s: %w"
	typeDiffers         = "the %s event %s"
	fieldDiffers        = "the %s event's %s field %s"
	noSuchEventRecorded = "is not in the recording at all"
	notInTheRecording   = "is not in the recording"
	readingARecording   = "reading %s: %w"
	noRecordings        = "there are no recordings under %s to compare against: run make fixtures against the real CLI first"
	writingARecording   = "writing %s: %w"
	runningTheRealCLI   = "running the real CLI: %w"
)
