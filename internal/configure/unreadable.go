package configure

import (
	"slices"
	"strings"

	"github.com/socrates-tools/socrates/internal/config"
)

const (
	indented   = " "
	tabbed     = "\t"
	afterAKey  = ":"
	aCommentAt = "#"
)

type unreadable struct {
	key    string
	line   int
	text   string
	reason string
	from   int
	to     int
}

type region struct {
	key       string
	depth     int
	from      int
	to        int
	ancestors []string
}

func unreadableIn(path string, content []byte) []unreadable {
	lines := strings.Split(string(content), newline)

	var found []unreadable
	for _, one := range regionsIn(lines) {
		if slices.ContainsFunc(found, func(reported unreadable) bool {
			return reported.from >= one.from && reported.to <= one.to
		}) {
			continue
		}

		err := config.Parse(path, one.inItsOwnRight(lines))
		if err == nil {
			continue
		}
		found = append(found, unreadable{
			key:    one.key,
			line:   one.from + 1,
			text:   strings.TrimSpace(lines[one.from]),
			reason: err.Error(),
			from:   one.from,
			to:     one.to,
		})
	}
	return found
}

func (r region) inItsOwnRight(lines []string) []byte {
	stated := append(append([]string{}, r.ancestors...), lines[r.from:r.to]...)
	return []byte(strings.Join(stated, newline) + newline)
}

func regionsIn(lines []string) []region {
	var open []region
	var found []region

	close := func(down int, at int) {
		for len(open) > 0 && open[len(open)-1].depth >= down {
			closing := open[len(open)-1]
			open = open[:len(open)-1]
			closing.to = endOf(lines, closing.from, at)
			found = append(found, closing)
		}
	}

	for position, line := range lines {
		if !startsAKey(line) {
			continue
		}

		depth := indentOf(line)
		close(depth, position)

		var ancestors []string
		for _, above := range open {
			ancestors = append(ancestors, lines[above.from])
		}
		open = append(open, region{
			key:       dotted(open, lines, keyOf(line)),
			depth:     depth,
			from:      position,
			ancestors: ancestors,
		})
	}
	close(0, len(lines))

	slices.SortFunc(found, func(one, other region) int {
		if one.depth != other.depth {
			return other.depth - one.depth
		}
		return one.from - other.from
	})
	return found
}

func dotted(open []region, lines []string, name string) string {
	if len(open) == 0 {
		return name
	}
	return open[len(open)-1].key + keySeparator + name
}

func startsAKey(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed != "" && !strings.HasPrefix(trimmed, aCommentAt)
}

func indentOf(line string) int {
	return len(line) - len(strings.TrimLeft(line, indented+tabbed))
}

func endOf(lines []string, from, next int) int {
	to := next
	for to > from+1 && !startsAKey(lines[to-1]) {
		to--
	}
	for to < next && strings.TrimSpace(lines[to]) == "" {
		to++
	}
	return to
}

func keyOf(line string) string {
	named, _, stated := strings.Cut(line, afterAKey)
	if !stated {
		return strings.TrimSpace(line)
	}
	return strings.TrimSpace(named)
}

func repaired(content []byte, one unreadable) ([]byte, error) {
	if document, err := documentOf(content); err == nil && document.states(one.key) {
		document.remove(one.key)
		return document.render()
	}
	return withoutTheLines(content, one), nil
}

func withoutTheLines(content []byte, one unreadable) []byte {
	lines := strings.Split(string(content), newline)
	kept := append(append([]string{}, lines[:one.from]...), lines[one.to:]...)

	written := strings.Join(kept, newline)
	if strings.TrimSpace(written) == "" {
		return nil
	}
	return []byte(written)
}
