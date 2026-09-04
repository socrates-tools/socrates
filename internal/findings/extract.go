package findings

import (
	"encoding/json"
	"errors"
	"strings"
)

const fenceCharacter = '`'

const shortestFence = 3

func Extract(message string) (json.RawMessage, error) {
	candidates := fencedBlocksIn(message)
	if len(candidates) == 0 {
		return nil, errors.New(noDocumentInTheFinalMessage)
	}

	last := strings.TrimSpace(candidates[len(candidates)-1])
	if last == "" {
		return nil, errors.New(theLastBlockIsEmpty)
	}
	if !json.Valid([]byte(last)) {
		return nil, errors.New(theLastBlockIsNotJSON)
	}
	return json.RawMessage(last), nil
}

func fencedBlocksIn(message string) []string {
	var blocks []string

	lines := strings.Split(message, "\n")
	for position := 0; position < len(lines); position++ {
		opening := fenceWidth(lines[position])
		if opening < shortestFence {
			continue
		}

		var held []string
		closed := false
		for position++; position < len(lines); position++ {
			if closes(lines[position], opening) {
				closed = true
				break
			}
			held = append(held, lines[position])
		}
		if closed {
			blocks = append(blocks, strings.Join(held, "\n"))
		}
	}
	return blocks
}

func fenceWidth(line string) int {
	trimmed := strings.TrimLeft(line, " ")

	width := 0
	for width < len(trimmed) && trimmed[width] == fenceCharacter {
		width++
	}
	return width
}

func closes(line string, opening int) bool {
	trimmed := strings.TrimSpace(line)
	if fenceWidth(trimmed) < opening {
		return false
	}
	return strings.Trim(trimmed, string(fenceCharacter)) == ""
}

func DeliveryInstruction() string {
	return strings.Join([]string{
		"End your final message with the findings document and nothing after it.",
		"Write it as one fenced block of JSON, opened with ```json and closed with ```, holding the whole document and nothing else.",
		"socrates reads the last fenced block of your final message, so anything you quote while reasoning must come before it.",
		"A final message with no fenced block, or one whose last block is not the document, fails the round.",
	}, " ")
}
