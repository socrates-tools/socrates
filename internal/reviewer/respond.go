package reviewer

import (
	"github.com/rui-ktei/socrates/internal/dialogue"
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
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
	if !dialogue.Names(document, response.Finding) {
		return Recording{}, refusal.New(noSuchFinding, round.Number, run.Name, response.Finding)
	}

	recorded, err := round.Responses()
	if err != nil {
		return Recording{}, err
	}
	replaced := recorded.Record(response)
	if err := round.SaveResponses(recorded); err != nil {
		return Recording{}, err
	}

	recording := Recording{
		Run:        run.Name,
		Round:      round.Number,
		Response:   response,
		Replaced:   replaced,
		Unanswered: recorded.Unanswered(document),
	}
	return recording, record(store, run, runlog.Event{
		Kind:    EventResponse,
		Subject: run.Name,
		Fields: map[string]any{
			"round":      round.Number,
			"finding":    response.Finding,
			"resolution": string(response.Resolution),
			"replaced":   replaced,
			"unanswered": len(recording.Unanswered),
		},
	})
}

func Unanswered(store runstore.Store, name string) ([]string, error) {
	_, round, document, err := parked(store, name)
	if err != nil {
		return nil, err
	}

	recorded, err := round.Responses()
	if err != nil {
		return nil, err
	}
	return recorded.Unanswered(document), nil
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
