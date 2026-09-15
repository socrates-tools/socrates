package reviewer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/socrates-tools/socrates/socrates-review/agent"
	"github.com/socrates-tools/socrates/socrates-review/dialogue"
	"github.com/socrates-tools/socrates/socrates-review/findings"
	"github.com/socrates-tools/socrates/socrates-review/guidance"
	"github.com/socrates-tools/socrates/socrates-review/hunks"
	"github.com/socrates-tools/socrates/socrates-review/prompt"
	"github.com/socrates-tools/socrates/socrates-review/runlog"
	"github.com/socrates-tools/socrates/socrates-review/runstate"
	"github.com/socrates-tools/socrates/socrates-review/runstore"
	"github.com/socrates-tools/socrates/socrates-review/snapshot"
)

const firstRound = 1

type Supervision struct {
	Store    runstore.Store
	Name     string
	Process  int
	Ceiling  int
	Recorder runlog.Recorder
	Vendor   agent.Vendor
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
	run.ReviewerCLI = supervision.Vendor.Name
	run.ReviewerTier = string(supervision.Vendor.Tier)
	if err := run.Advance(runstate.Reviewing); err != nil {
		return err
	}
	if err := supervision.Store.Save(run); err != nil {
		return err
	}

	round, present, err := supervision.Store.LatestRound(run.Name)
	if err != nil {
		return err
	}
	if !present {
		return supervision.end(run, runstate.Failed, noRoundToReview)
	}

	reviewed, err := changeIn(ctx, run)
	if err != nil {
		return supervision.end(run, runstate.Failed, err.Error())
	}

	prior, err := priorRounds(supervision.Store, run.Name, round.Number)
	if err != nil {
		return supervision.end(run, runstate.Failed, err.Error())
	}

	guide, err := guidance.At(ctx, run.Repository, run.Base, run.GuidanceCeiling)
	if err != nil {
		return supervision.end(run, runstate.Failed, err.Error())
	}

	if err := round.WriteChange(hunks.Of(reviewed.diff)); err != nil {
		return err
	}

	asked := prompt.Build(prompt.Request{
		Intent:     run.Intent,
		NotMine:    run.NotMyChange,
		Repository: run.Repository,
		Base:       run.Base,
		Snapshot:   run.Snapshot,
		Changed:    reviewed.paths,
		Diff:       reviewed.diff,
		Ceiling:    supervision.Ceiling,
		Prior:      prior,
		Guidance:   guide.Text,

		DocumentEndsTheMessage: supervision.Vendor.Tier == agent.Parsed,
	})

	run.Diff = asked.Diff
	if err := supervision.Store.Save(run); err != nil {
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
		_ = stream.Close()
		return err
	}

	answer, agentErr := reviewing.Run(ctx, agent.Request{
		Role:       agent.Reviewer,
		Prompt:     asked.Text,
		Dir:        run.Worktree,
		Scratch:    round.Directory,
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

	classified := dialogue.Classified(document, run.Authorship())
	outcome := dialogue.After(classified, round.Number)

	if outcome == runstate.AwaitingResponse {
		removed := removeWorktree(ctx, worktree)
		return errors.Join(removed, supervision.park(run, round.Number, classified, removed))
	}

	if err := supervision.end(run, outcome, verdictReason(document, classified, round.Number)); err != nil {
		return err
	}
	return removeWorktree(ctx, worktree)
}

func (s Supervision) park(run *runstore.Run, round int, classified findings.Classification, teardown error) error {
	run.Reviewer = 0
	run.Agent = 0
	if err := run.Advance(runstate.AwaitingResponse); err != nil {
		return err
	}
	if err := s.Store.Save(run); err != nil {
		return err
	}

	lifted := liftTheSeal(run)
	if s.Recorder == nil {
		return lifted
	}

	fields := map[string]any{
		"round":            round,
		"patchable":        len(dialogue.AnsweredByAPatch(classified)),
		"needing_a_person": len(dialogue.AnsweredByAPerson(classified)),
	}
	if teardown != nil {
		fields["worktree"] = teardown.Error()
	}
	return errors.Join(lifted, s.Recorder.Record(runlog.Event{
		Kind:    EventRoundParked,
		Subject: run.Name,
		Fields:  fields,
	}))
}

func priorRounds(store runstore.Store, name string, current int) ([]prompt.PriorRound, error) {
	numbers, err := store.Rounds(name)
	if err != nil {
		return nil, err
	}

	var prior []prompt.PriorRound
	for _, number := range numbers {
		if number >= current {
			continue
		}

		round := store.Round(name, number)
		document, reported, err := round.Findings()
		if err != nil {
			return nil, err
		}
		if !reported {
			continue
		}
		recorded, err := round.Responses()
		if err != nil {
			return nil, err
		}
		prior = append(prior, prompt.PriorRound{
			Number:    number,
			Findings:  document.Findings,
			Responses: recorded.Latest(),
		})
	}
	return prior, nil
}

func (s Supervision) end(run *runstore.Run, outcome runstate.State, reason string) error {
	if err := run.End(outcome, reason); err != nil {
		return err
	}
	if err := s.Store.Save(run); err != nil {
		return err
	}

	lifted := liftTheSeal(run)
	if s.Recorder == nil {
		return lifted
	}
	return errors.Join(lifted, s.Recorder.Record(runlog.Event{
		Kind:    EventRunEnded,
		Subject: run.Name,
		Fields:  map[string]any{"state": string(outcome), "reason": reason},
	}))
}

func liftTheSeal(run *runstore.Run) error {
	if run.Worktree == "" {
		return nil
	}
	if err := snapshot.Unseal(run.Worktree); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

func verdictReason(document findings.Document, classified findings.Classification, round int) string {
	if len(document.Findings) == 0 {
		return nothingToReport
	}

	found := fmt.Sprintf(whatAReviewerFound, len(document.Findings), len(classified.NeedingAPerson), document.RiskLevel)
	if beyond := dialogue.BeyondAnswering(classified); len(beyond) > 0 {
		return found + fmt.Sprintf(nothingAddressesThem, len(beyond))
	}
	if dialogue.EscalatedByTheCap(classified, round) {
		return found + fmt.Sprintf(theCapWasReached, dialogue.Cap)
	}
	return found
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
