package presentation

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/socrates-tools/socrates/review/runlog"
)

const shortened = 7

func Lines(events []runlog.Event) []string {
	var lines []string

	if slices.ContainsFunc(events, wasRotation) {
		lines = append(lines, rotationNotice)
	}
	for _, event := range events {
		if wasRotation(event) {
			continue
		}
		lines = append(lines, fmt.Sprintf(eventLine, event.At.UTC().Format(clockLayout), sentence(event)))
	}
	return lines
}

func wasRotation(event runlog.Event) bool {
	return event.Kind == rotatedKind
}

func sentence(event runlog.Event) string {
	held := event.Fields

	switch event.Kind {
	case runStartedKind:
		started := fmt.Sprintf(runStartedSentence, stated(held, repositoryField), short(stated(held, snapshotField)), short(stated(held, baseField)))
		if flagged(held, nobodyToAskField) {
			return started + runStartedWithNobodyToAsk
		}
		return started
	case runEndedKind:
		if reason := stated(held, reasonField); reason != "" {
			return fmt.Sprintf(runEndedSentence, stated(held, stateField), reason)
		}
		return fmt.Sprintf(runEndedWithNoReason, stated(held, stateField))
	case roundStartedKind:
		return fmt.Sprintf(roundStartedSentence, counted(held, roundField), short(stated(held, snapshotField)), short(stated(held, baseField)))
	case roundParkedKind:
		asked := counted(held, aPersonField)
		if asked > 0 && flagged(held, nobodyToAskField) {
			return fmt.Sprintf(roundParkedForTheSession, counted(held, roundField), counted(held, patchableField), asked)
		}
		if asked > 0 {
			return fmt.Sprintf(roundParkedOnAPerson, counted(held, roundField), counted(held, patchableField), asked)
		}
		return fmt.Sprintf(roundParkedSentence, counted(held, roundField), counted(held, patchableField))
	case responseKind:
		return fmt.Sprintf(responseSentence, stated(held, findingField), counted(held, roundField), stated(held, resolutionField), counted(held, unansweredField))
	case agentStartedKind:
		return fmt.Sprintf(agentStartedSentence, event.Subject, stated(held, modelField), stated(held, effortField), counted(held, processField))
	case agentCompletedKind:
		return fmt.Sprintf(agentCompletedSentence, event.Subject, stated(held, modelField), counted(held, inputTokensField), counted(held, outputTokensField))
	case agentFailedKind:
		return fmt.Sprintf(agentFailedSentence, event.Subject, stated(held, reasonField))
	case unexpectedToolsKind:
		return fmt.Sprintf(unexpectedToolsSentence, stated(held, toolsField))
	case unexpectedSourcesKind:
		return fmt.Sprintf(unexpectedSourcesSentence, stated(held, sourcesField))
	case guidanceUnreadKind:
		return fmt.Sprintf(guidanceUnreadSentence, stated(held, guideField))
	case guidanceTooLargeKind:
		return fmt.Sprintf(guidanceTooLargeSentence, stated(held, guideField), counted(held, sizeField), counted(held, ceilingField))
	case toolsUnreportedKind:
		return fmt.Sprintf(toolsUnreportedSentence, stated(held, askedForField))
	default:
		return unrecognised(event)
	}
}

func unrecognised(event runlog.Event) string {
	stated := fmt.Sprintf(unknownSentence, event.Kind, event.Subject)
	if len(event.Fields) == 0 {
		return stated
	}

	var pairs []string
	for _, key := range slices.Sorted(maps.Keys(event.Fields)) {
		pairs = append(pairs, fmt.Sprintf(fieldPair, key, event.Fields[key]))
	}
	return fmt.Sprintf(unknownFields, stated, strings.Join(pairs, fieldSeparator))
}

func stated(fields map[string]any, key string) string {
	text, isText := fields[key].(string)
	if !isText {
		return ""
	}
	return text
}

func flagged(fields map[string]any, key string) bool {
	set, isFlag := fields[key].(bool)
	return isFlag && set
}

func counted(fields map[string]any, key string) int64 {
	switch held := fields[key].(type) {
	case float64:
		return int64(held)
	case int:
		return int64(held)
	case int64:
		return held
	default:
		return 0
	}
}

func short(commit string) string {
	if len(commit) <= shortened {
		return commit
	}
	return commit[:shortened]
}
