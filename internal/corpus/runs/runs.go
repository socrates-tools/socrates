package runs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rui-ktei/socrates/internal/corpus"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstore"
)

const (
	agentStarted   = "agent-started"
	reviewerAgent  = "reviewer"
	modelField     = "model"
	effortField    = "effort"
	writtenNewline = "\n"
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

	routing, err := routingOf(store, run.Name)
	if err != nil {
		return Answer{}, err
	}

	return Answer{
		Repository: run.Repository,
		Entry: corpus.Entry{
			Name:       request.Name,
			Model:      routing[modelField],
			Effort:     routing[effortField],
			Round:      request.Round,
			Base:       reviewed.Base,
			Snapshot:   reviewed.Snapshot,
			Run:        run.Name,
			RecordedOn: request.On,
			Subject:    request.Subject,
			Document:   []byte(strings.TrimSuffix(string(document), writtenNewline)),
		},
	}, nil
}

func routingOf(store runstore.Store, run string) (map[string]string, error) {
	path := filepath.Join(store.Path(run), runlog.FileName)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(readingTheLog, path, err)
	}
	defer file.Close()

	routing := map[string]string{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	for scanner.Scan() {
		var event runlog.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf(readingTheLog, path, err)
		}
		if event.Kind != agentStarted || event.Subject != reviewerAgent {
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
		return nil, fmt.Errorf(logStatesNothingFor, path, strings.Join(absent, ", "))
	}
	return routing, nil
}
