package reviewer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/prompt"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/snapshot"
)

const firstRound = 1

type Supervision struct {
	Store    runstore.Store
	Name     string
	Process  int
	Ceiling  int
	Recorder runlog.Recorder
	NewAgent func(runlog.Recorder) (agent.Agent, error)
}

func Supervise(ctx context.Context, supervision Supervision) error {
	run, err := supervision.Store.Load(supervision.Name)
	if err != nil {
		return err
	}
	if run.State.Terminal() {
		return fmt.Errorf(alreadyEnded, run.Name, run.State)
	}

	run.Reviewer = supervision.Process
	if err := run.Advance(runstate.Reviewing); err != nil {
		return err
	}
	if err := supervision.Store.Save(run); err != nil {
		return err
	}

	reviewed, err := changeIn(ctx, run)
	if err != nil {
		return supervision.end(run, runstate.Failed, err.Error())
	}

	asked := prompt.Build(prompt.Request{
		Intent:     run.Intent,
		Repository: run.Repository,
		Base:       run.Base,
		Snapshot:   run.Snapshot,
		Changed:    reviewed.paths,
		Diff:       reviewed.diff,
		Ceiling:    supervision.Ceiling,
	})

	run.Diff = asked.Diff
	if err := supervision.Store.Save(run); err != nil {
		return err
	}

	round, err := supervision.Store.CreateRound(run.Name, firstRound)
	if err != nil {
		return err
	}
	if err := round.Write(round.PromptPath(), asked.Text); err != nil {
		return err
	}

	stream, err := os.Create(round.StreamPath())
	if err != nil {
		return fmt.Errorf(openingTheStream, round.StreamPath(), err)
	}

	noting := &noteAgentProcess{inner: supervision.Recorder, store: supervision.Store, run: run}
	reviewing, err := supervision.NewAgent(noting)
	if err != nil {
		stream.Close()
		return err
	}

	answer, agentErr := reviewing.Run(ctx, agent.Request{
		Role:       agent.Reviewer,
		Prompt:     asked.Text,
		Dir:        run.Worktree,
		Schema:     json.RawMessage(findings.Schema),
		Honours:    honoursTheSchema,
		Transcript: stream,
	})
	if err := stream.Close(); err != nil {
		return fmt.Errorf(closingTheStream, round.StreamPath(), err)
	}

	if len(answer.Structured) > 0 {
		if err := round.Write(round.FindingsPath(), string(answer.Structured)+"\n"); err != nil {
			return err
		}
	}

	worktree := snapshot.Worktree{RepositoryRoot: run.Repository, Path: run.Worktree, Commit: run.Snapshot}
	written, compareErr := worktree.Differences(ctx)
	if compareErr != nil {
		return supervision.end(run, runstate.Failed, fmt.Sprintf(couldNotCompare, compareErr))
	}

	if len(written) > 0 {
		reason := fmt.Sprintf(theReviewerWrote, strings.Join(written, ", "), run.Worktree)
		if agentErr != nil {
			reason += fmt.Sprintf(andTheReviewFailed, agentErr)
		}
		return supervision.end(run, runstate.Failed, reason)
	}

	if agentErr != nil {
		reason := agentErr.Error()
		if len(answer.Structured) > 0 {
			reason += fmt.Sprintf(whatItReturnedIsKept, round.FindingsPath())
		}
		if err := supervision.end(run, runstate.Failed, reason); err != nil {
			return err
		}
		return removeWorktree(ctx, worktree)
	}

	document, err := findings.Parse(answer.Structured)
	if err != nil {
		if endErr := supervision.end(run, runstate.Failed, fmt.Sprintf(schemaNotSatisfied, err, round.FindingsPath())); endErr != nil {
			return endErr
		}
		return removeWorktree(ctx, worktree)
	}

	if err := supervision.end(run, document.Verdict(), verdictReason(document)); err != nil {
		return err
	}
	return removeWorktree(ctx, worktree)
}

func (s Supervision) end(run *runstore.Run, outcome runstate.State, reason string) error {
	if err := run.End(outcome, reason); err != nil {
		return err
	}
	if err := s.Store.Save(run); err != nil {
		return err
	}
	if s.Recorder == nil {
		return nil
	}
	return s.Recorder.Record(runlog.Event{
		Kind:    EventRunEnded,
		Subject: run.Name,
		Fields:  map[string]any{"state": string(outcome), "reason": reason},
	})
}

func verdictReason(document findings.Document) string {
	classified := findings.Classify(document)
	if len(document.Findings) == 0 {
		return nothingToReport
	}
	return fmt.Sprintf(whatAReviewerFound, len(document.Findings), len(classified.NeedingAPerson), document.RiskLevel)
}

func honoursTheSchema(structured json.RawMessage) error {
	_, err := findings.Parse(structured)
	return err
}

type change struct {
	paths []string
	diff  string
}

func changeIn(ctx context.Context, run *runstore.Run) (change, error) {
	paths, err := snapshot.Changed(ctx, run.Repository, run.Base, run.Snapshot)
	if err != nil {
		return change{}, fmt.Errorf(noChangeList, err)
	}

	diff, err := snapshot.Diff(ctx, run.Repository, run.Base, run.Snapshot)
	if err != nil {
		return change{}, fmt.Errorf(noDiff, err)
	}

	return change{paths: paths, diff: diff}, nil
}

func removeWorktree(ctx context.Context, worktree snapshot.Worktree) error {
	if worktree.Path == "" {
		return nil
	}
	return worktree.Remove(ctx)
}

type noteAgentProcess struct {
	inner runlog.Recorder
	store runstore.Store
	run   *runstore.Run
}

func (n *noteAgentProcess) Record(event runlog.Event) error {
	if event.Kind == agent.EventStarted {
		if pid, stated := event.Fields["process"].(int); stated && pid != n.run.Agent {
			n.run.Agent = pid
			if err := n.store.Save(n.run); err != nil {
				return err
			}
		}
	}
	if n.inner == nil {
		return nil
	}
	return n.inner.Record(event)
}
