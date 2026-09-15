package main

const (
	failurePrefix = "corpus-record: "

	runFlag     = "the run whose answer is recorded"
	roundFlag   = "the round of that run"
	nameFlag    = "the name the corpus holds the answer under"
	subjectFlag = "one line naming what was reviewed"
	onFlag      = "the day the answer is recorded, as YYYY-MM-DD"

	nothingToRecord = "name the run, the entry and the subject: corpus-record -run <run> -name <entry> -subject <one line>"
	recorded        = "recorded %s: %s at %s effort, round %d, in %s"
)
