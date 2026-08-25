package render

import (
	"strconv"
	"strings"
)

type Kind string

const (
	Context Kind = "context"
	Removed Kind = "removed"
	Added   Kind = "added"
	Absent  Kind = "absent"
)

type Side struct {
	Kind   Kind
	Number int
	Text   string
}

type Row struct {
	Left  Side
	Right Side
}

type Hunk struct {
	Heading string
	Rows    []Row
}

type Diff struct {
	From  string
	To    string
	Hunks []Hunk
}

func ReadPatch(patch string) Diff {
	var read Diff
	var hunk *Hunk
	var removals, additions []Side

	oldNumber, newNumber := 0, 0

	flush := func() {
		if hunk == nil {
			return
		}
		for position := 0; position < len(removals) || position < len(additions); position++ {
			row := Row{Left: Side{Kind: Absent}, Right: Side{Kind: Absent}}
			if position < len(removals) {
				row.Left = removals[position]
			}
			if position < len(additions) {
				row.Right = additions[position]
			}
			hunk.Rows = append(hunk.Rows, row)
		}
		removals, additions = nil, nil
	}

	end := func() {
		flush()
		if hunk != nil {
			read.Hunks = append(read.Hunks, *hunk)
			hunk = nil
		}
	}

	for line := range strings.SplitSeq(patch, "\n") {
		switch {
		case strings.HasPrefix(line, "--- "):
			end()
			read.From = strings.TrimPrefix(line, "--- ")
		case strings.HasPrefix(line, "+++ "):
			end()
			read.To = strings.TrimPrefix(line, "+++ ")
		case strings.HasPrefix(line, "@@"):
			end()
			heading, from, to := readHunkHeading(line)
			oldNumber, newNumber = from, to
			hunk = &Hunk{Heading: heading}
		case hunk == nil, line == "":
			continue
		case strings.HasPrefix(line, "-"):
			removals = append(removals, Side{Kind: Removed, Number: oldNumber, Text: line[1:]})
			oldNumber++
		case strings.HasPrefix(line, "+"):
			additions = append(additions, Side{Kind: Added, Number: newNumber, Text: line[1:]})
			newNumber++
		case strings.HasPrefix(line, "\\"):
			continue
		default:
			flush()
			text := strings.TrimPrefix(line, " ")
			hunk.Rows = append(hunk.Rows, Row{
				Left:  Side{Kind: Context, Number: oldNumber, Text: text},
				Right: Side{Kind: Context, Number: newNumber, Text: text},
			})
			oldNumber++
			newNumber++
		}
	}
	end()

	return read
}

func readHunkHeading(line string) (string, int, int) {
	within := strings.TrimPrefix(line, "@@")
	ranges, trailing, closed := strings.Cut(within, "@@")
	if !closed {
		return "", 1, 1
	}

	from, to := 1, 1
	for stated := range strings.FieldsSeq(ranges) {
		start, _, _ := strings.Cut(strings.TrimLeft(stated, "-+"), ",")
		counted, err := strconv.Atoi(start)
		if err != nil {
			continue
		}
		if strings.HasPrefix(stated, "-") {
			from = counted
		}
		if strings.HasPrefix(stated, "+") {
			to = counted
		}
	}
	return strings.TrimSpace(trailing), from, to
}

func (s Side) Numbered() string {
	if s.Kind == Absent {
		return ""
	}
	return strconv.Itoa(s.Number)
}

func (s Side) Marker() string {
	switch s.Kind {
	case Removed:
		return "-"
	case Added:
		return "+"
	case Context:
		return " "
	default:
		return ""
	}
}

func (d Diff) Empty() bool {
	return len(d.Hunks) == 0
}
