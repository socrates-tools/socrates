package runs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/socrates-tools/socrates/internal/agent"
	"github.com/socrates-tools/socrates/internal/corpus"
	"github.com/socrates-tools/socrates/internal/findings"
	"github.com/socrates-tools/socrates/internal/reviewer"
	"github.com/socrates-tools/socrates/internal/runlog"
	"github.com/socrates-tools/socrates/internal/runstore"
)

const (
	reviewerAgent  = "reviewer"
	modelField     = "model"
	effortField    = "effort"
	roundField     = "round"
	writtenNewline = "\n"
	firstRound     = 1
)

type Request struct {
	Run     string
	Round   int
	Name    string
	Subject string
	On      string
}

type Answer struct {
	Entry      corpus.Entry
	Repository string
}

func Answered(store runstore.Store, request Request) (Answer, error) {
	run, err := store.Load(request.Run)
	if err != nil {
		return Answer{}, err
	}

	reviewed, found := run.Reviewed(request.Round)
	if !found {
		return Answer{}, fmt.Errorf(noSuchRound, request.Round, run.Name)
	}

	document, err := os.ReadFile(store.Round(run.Name, request.Round).FindingsPath())
	if err != nil {
		return Answer{}, fmt.Errorf(roundAnsweredNothing, request.Round, run.Name, err)
	}
	returned := strings.TrimSuffix(string(document), writtenNewline)
	if _, err := findings.Parse([]byte(returned)); err != nil {
		return Answer{}, fmt.Errorf(roundAnsweredSomethingUnreadable, request.Round, run.Name, err)
	}

	routing, err := routingOf(store, run.Name, request.Round)
	if err != nil {
		return Answer{}, err
	}

	return Answer{
		Repository: run.Repository,
		Entry: corpus.Entry{
			Name:       request.Name,
			Reviewer:   run.Vendor(),
			Model:      routing[modelField],
			Effort:     routing[effortField],
			Round:      request.Round,
			Base:       reviewed.Base,
			Snapshot:   reviewed.Snapshot,
			Run:        run.Name,
			RecordedOn: request.On,
			Subject:    request.Subject,
			Document:   []byte(returned),
		},
	}, nil
}

func routingOf(store runstore.Store, run string, round int) (map[string]string, error) {
	path := filepath.Join(store.Path(run), runlog.FileName)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(readingTheLog, path, err)
	}
	defer func() { _ = file.Close() }()

	routing := map[string]string{}
	reached := firstRound
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	for scanner.Scan() {
		var event runlog.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf(readingTheLog, path, err)
		}

		if event.Kind == reviewer.EventRoundStarted {
			started, stated := event.Fields[roundField].(float64)
			if !stated {
				return nil, fmt.Errorf(logStatesNoRound, path)
			}
			reached = int(started)
			continue
		}
		if event.Kind != agent.EventStarted || event.Subject != reviewerAgent || reached != round {
			continue
		}
		for _, field := range []string{modelField, effortField} {
			if stated, said := event.Fields[field].(string); said && stated != "" {
				routing[field] = stated
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(readingTheLog, path, err)
	}

	var absent []string
	for _, field := range []string{modelField, effortField} {
		if routing[field] == "" {
			absent = append(absent, field)
		}
	}
	if len(absent) > 0 {
		return nil, fmt.Errorf(logStatesNothingFor, round, path, strings.Join(absent, ", "))
	}
	return routing, nil
}
