package settingsmenu

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func (m *menu) lines() []string {
	switch m.stage {
	case choosing:
		return m.optionLines()
	case repairing:
		return m.repairLines()
	}
	return m.settingLines()
}

func (m *menu) settingLines() []string {
	lines := m.saidLines()
	lines = append(lines, whatToConfigure, "")

	keyColumn, valueColumn := 0, 0
	for _, setting := range m.settings {
		keyColumn = max(keyColumn, len(setting.Key))
		valueColumn = max(valueColumn, len(setting.Value))
	}

	for position, setting := range m.settings {
		lines = append(lines, fmt.Sprintf(settingRow, marker(position == m.cursor), keyColumn, setting.Key, valueColumn, setting.Value, setting.Source))
	}
	return append(lines, "", movingThroughAMenu)
}

func (m *menu) optionLines() []string {
	lines := m.saidLines()
	lines = append(lines, m.chosen.Key, "")

	for position, offered := range m.options {
		at := position == m.at
		if offered.editable {
			lines = append(lines, fmt.Sprintf(offeredRow, marker(at), m.typed, cursorFor(at)))
			continue
		}
		lines = append(lines, fmt.Sprintf(offeredRow, marker(at), offered.label, inForceMark(offered, m.chosen.Value)))
	}
	return append(lines, "", m.movingThroughTheOptions())
}

func (m *menu) movingThroughTheOptions() string {
	for _, offered := range m.options {
		if offered.editable {
			return typingAValue
		}
	}
	return movingThroughAList
}

func (m *menu) repairLines() []string {
	one := m.broken[0]

	lines := m.wrapped(fmt.Sprintf(m.definition.Wording.CannotBeUnderstood, m.path), "")
	lines = append(lines, "")
	lines = append(lines, m.wrapped(fmt.Sprintf(theOffendingLine, one.line, one.text), indent)...)
	for _, said := range strings.Split(one.reason, newline) {
		lines = append(lines, m.wrapped(indent+said, indent)...)
	}
	lines = append(lines, "")

	for position, offered := range m.options {
		lines = append(lines, fmt.Sprintf(offeredRow, marker(position == m.at), offered.label, ""))
	}
	return append(lines, "", movingThroughAMenu)
}

func (m *menu) saidLines() []string {
	if len(m.said) == 0 {
		return nil
	}
	var lines []string
	for _, said := range m.said {
		lines = append(lines, m.wrapped(said, "")...)
	}
	return append(lines, "")
}

func (m *menu) wrapped(line, continuation string) []string {
	if m.width <= len(continuation) {
		return []string{line}
	}
	wrapped := strings.Split(ansi.Wrap(line, m.width-len(continuation), string(filepath.Separator)), newline)
	for at := 1; at < len(wrapped); at++ {
		wrapped[at] = continuation + wrapped[at]
	}
	return wrapped
}

func marker(at bool) string {
	if at {
		return here
	}
	return elsewhere
}

func cursorFor(at bool) string {
	if at {
		return typingCursor
	}
	return ""
}

func inForceMark(offered option, value string) string {
	if !offered.defaults && offered.value == value {
		return inForce
	}
	return ""
}
