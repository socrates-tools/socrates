package hygiene

import "strings"

var conflictMarkers = []string{"<<<<<<<", "=======", ">>>>>>>", "|||||||"}

func Neutralise(field string) string {
	var lines []string

	for line := range strings.SplitSeq(withoutCarriageReturns(field), "\n") {
		lines = append(lines, strings.Join(strings.Fields(withoutConflictMarkers(line)), " "))
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
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
