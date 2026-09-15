package reviewer

import (
	"fmt"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/socrates-review/dialogue"
	"github.com/socrates-tools/socrates/socrates-review/findings"
	"github.com/socrates-tools/socrates/socrates-review/runlog"
	"github.com/socrates-tools/socrates/socrates-review/runstate"
	"github.com/socrates-tools/socrates/socrates-review/runstore"
)

type Recording struct {
	Run        string
	Round      int
	Response   dialogue.Response
	Replaced   bool
	Unanswered []string
}

func Respond(store runstore.Store, name string, response dialogue.Response) (Recording, error) {
	run, round, document, err := parked(store, name)
	if err != nil {
		return Recording{}, err
	}

	if err := dialogue.Validate(response); err != nil {
		return Recording{}, err
	}
	reported, present := dialogue.Reported(document, response.Finding)
	if !present {
		return Recording{}, refusal.New(noSuchFinding, round.Number, run.Name, response.Finding)
	}
	if err := dialogue.Answering(reported, response, run.Authorship()); err != nil {
		return Recording{}, err
	}

	if err := round.AppendResponse(response); err != nil {
		return Recording{}, err
	}

	recorded, readBack := round.Responses()
	recording := Recording{
		Run:        run.Name,
		Round:      round.Number,
		Response:   response,
		Replaced:   recorded.Answers(response.Finding) > 1,
		Unanswered: recorded.Unanswered(document, run.Authorship()),
	}

	fields := map[string]any{
		"round":      round.Number,
		"finding":    response.Finding,
		"resolution": string(response.Resolution),
		"replaced":   recording.Replaced,
		"unanswered": len(recording.Unanswered),
	}
	if readBack != nil {
		fields["read_back"] = readBack.Error()
	}
	if err := record(store, run, runlog.Event{Kind: EventResponse, Subject: run.Name, Fields: fields}); err != nil {
		return recording, err
	}
	if readBack != nil {
		return recording, fmt.Errorf(storedAndUnreadable, round.ResponsesPath(), readBack)
	}
	return recording, nil
}

func Unanswered(store runstore.Store, name string) ([]string, error) {
	run, round, document, err := parked(store, name)
	if err != nil {
		return nil, err
	}

	recorded, err := round.Responses()
	if err != nil {
		return nil, err
	}
	return recorded.Unanswered(document, run.Authorship()), nil
}

func parked(store runstore.Store, name string) (*runstore.Run, runstore.Round, findings.Document, error) {
	run, err := store.Load(name)
	if err != nil {
		return nil, runstore.Round{}, findings.Document{}, err
	}
	if run.State != runstate.AwaitingResponse {
		return nil, runstore.Round{}, findings.Document{}, refusal.New(notAwaitingResponses, run.Name, run.State)
	}

	round, present, err := store.LatestRound(name)
	if err != nil {
		return nil, runstore.Round{}, findings.Document{}, err
	}
	if !present {
		return nil, runstore.Round{}, findings.Document{}, refusal.New(noRoundToRespondTo, run.Name)
	}

	document, reported, err := round.Findings()
	if err != nil {
		return nil, runstore.Round{}, findings.Document{}, err
	}
	if !reported {
		return nil, runstore.Round{}, findings.Document{}, refusal.New(noFindingsToRespondTo, round.Number, run.Name)
	}
	return run, round, document, nil
}
