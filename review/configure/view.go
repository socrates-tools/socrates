package configure

import (
	"fmt"
	"strings"
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

	lines := []string{fmt.Sprintf(cannotBeUnderstood, m.path), ""}
	lines = append(lines, fmt.Sprintf(theOffendingLine, one.line, one.text))
	for _, said := range strings.Split(one.reason, newline) {
		lines = append(lines, indent+said)
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
	return append(append([]string{}, m.said...), "")
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
