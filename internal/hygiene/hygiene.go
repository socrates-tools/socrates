package hygiene

import "strings"

var conflictMarkers = []string{"<<<<<<<", "=======", ">>>>>>>", "|||||||"}

const headingMarker = "#"

func Neutralise(field string) string {
	var lines []string

	for line := range strings.SplitSeq(withoutCarriageReturns(field), "\n") {
		lines = append(lines, withoutHeadingMarker(collapsed(withoutConflictMarkers(line))))
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func collapsed(line string) string {
	return strings.Join(strings.Fields(line), " ")
}

func withoutHeadingMarker(line string) string {
	named := strings.TrimLeft(line, headingMarker)
	if len(named) == len(line) {
		return line
	}
	if named != "" && !strings.HasPrefix(named, " ") {
		return line
	}
	return strings.TrimSpace(named)
}

func withoutCarriageReturns(field string) string {
	return strings.ReplaceAll(strings.ReplaceAll(field, "\r\n", "\n"), "\r", "\n")
}

func withoutConflictMarkers(line string) string {
	cleaned := line
	for _, marker := range conflictMarkers {
		cleaned = strings.ReplaceAll(cleaned, marker, " ")
	}
	return cleaned
}
