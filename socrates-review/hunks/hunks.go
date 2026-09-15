package hunks

import (
	"strconv"
	"strings"
)

const (
	OldSide = "a/"
	NewSide = "b/"
)

const (
	fileMarker   = "+++ "
	hunkMarker   = "@@ "
	changeMarker = "diff --git "
	nothing      = "/dev/null"
	noNewline    = '\\'
)

type Span struct {
	First int `json:"first"`
	Last  int `json:"last"`
}

type Change struct {
	Touched map[string][]Span `json:"touched"`
}

type hunk struct {
	oldLines int
	newFirst int
	newLines int
}

func Of(diff string) Change {
	touched := map[string][]Span{}

	file := ""
	remainingOld, remainingNew := 0, 0

	for line := range strings.SplitSeq(diff, "\n") {
		if remainingOld > 0 || remainingNew > 0 {
			left, right, body := bodyLine(line, remainingOld, remainingNew)
			if body {
				remainingOld, remainingNew = left, right
				continue
			}
			remainingOld, remainingNew = 0, 0
		}

		switch {
		case strings.HasPrefix(line, changeMarker):
			file = ""
		case strings.HasPrefix(line, fileMarker):
			file = newSidePath(strings.TrimPrefix(line, fileMarker))
		case strings.HasPrefix(line, hunkMarker):
			read, understood := headerOf(line)
			if !understood {
				continue
			}
			remainingOld, remainingNew = read.oldLines, read.newLines
			if file == "" || read.newLines < 1 {
				continue
			}
			touched[file] = append(touched[file], Span{First: read.newFirst, Last: read.newFirst + read.newLines - 1})
		}
	}

	return Change{Touched: touched}
}

func (c Change) Holds(file string, line int) bool {
	for _, span := range c.Touched[file] {
		if line >= span.First && line <= span.Last {
			return true
		}
	}
	return false
}

func bodyLine(line string, remainingOld, remainingNew int) (int, int, bool) {
	if line == "" {
		return remainingOld - 1, remainingNew - 1, true
	}
	switch line[0] {
	case ' ':
		return remainingOld - 1, remainingNew - 1, true
	case '-':
		return remainingOld - 1, remainingNew, true
	case '+':
		return remainingOld, remainingNew - 1, true
	case noNewline:
		return remainingOld, remainingNew, true
	default:
		return 0, 0, false
	}
}

func newSidePath(stated string) string {
	named, _, _ := strings.Cut(stated, "\t")
	if named == nothing {
		return ""
	}
	if unquoted, err := strconv.Unquote(named); strings.HasPrefix(named, `"`) && err == nil {
		named = unquoted
	}
	return strings.TrimPrefix(named, NewSide)
}

func headerOf(line string) (hunk, bool) {
	stated, _, _ := strings.Cut(strings.TrimPrefix(line, hunkMarker), " @@")

	var read hunk
	sides := 0

	for field := range strings.FieldsSeq(stated) {
		side := field[:1]
		first, count, understood := rangeOf(field[1:])
		if !understood {
			return hunk{}, false
		}

		switch side {
		case "-":
			read.oldLines = count
			sides++
		case "+":
			read.newFirst, read.newLines = first, count
			sides++
		}
	}

	return read, sides == 2
}

func rangeOf(stated string) (int, int, bool) {
	start, counted, ranged := strings.Cut(stated, ",")

	first, err := strconv.Atoi(start)
	if err != nil {
		return 0, 0, false
	}
	if !ranged {
		return first, 1, true
	}
	count, err := strconv.Atoi(counted)
	if err != nil {
		return 0, 0, false
	}
	return first, count, true
}
