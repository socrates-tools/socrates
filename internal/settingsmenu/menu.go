package settingsmenu

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type stage int

const (
	listing stage = iota
	choosing
	repairing
)

type option struct {
	label    string
	value    string
	defaults bool
	editable bool
	quits    bool
}

type menu struct {
	definition Definition
	path       string
	settings   []Setting
	cursor     int

	stage   stage
	chosen  Setting
	options []option
	at      int
	typed   string

	broken []unreadable
	said   []string
	failed error

	width int
}

func opened(definition Definition) (*menu, error) {
	standing := &menu{definition: definition, path: definition.Path}
	if err := standing.reread(); err != nil {
		return nil, err
	}
	return standing, nil
}

func (m *menu) reread() error {
	content, err := m.definition.held(m.path)
	if err != nil {
		return err
	}

	if err := m.definition.Parse(m.path, content); err != nil {
		m.broken = m.definition.unreadableIn(m.path, content)
		if len(m.broken) == 0 {
			return err
		}
		m.stage = repairing
		m.options = repairOptions(m.broken[0])
		m.at = 0
		return nil
	}

	loaded, err := m.definition.Settings(m.path)
	if err != nil {
		return err
	}
	m.settings = loaded
	m.broken = nil
	m.stage = listing
	if m.cursor >= len(m.settings) {
		m.cursor = 0
	}
	return nil
}

func (m *menu) Init() tea.Cmd {
	return nil
}

func (m *menu) View() tea.View {
	return tea.NewView(strings.Join(m.lines(), newline) + newline)
}

func (m *menu) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if resized, isResized := message.(tea.WindowSizeMsg); isResized {
		m.width = resized.Width
		return m, nil
	}
	pressed, isKey := message.(tea.KeyPressMsg)
	if !isKey {
		return m, nil
	}
	if pressed.Mod == tea.ModCtrl && (pressed.Code == 'c' || pressed.Code == 'd') {
		return m, tea.Quit
	}

	switch m.stage {
	case listing:
		return m.moveThroughTheSettings(pressed)
	case choosing:
		return m.moveThroughTheOptions(pressed)
	case repairing:
		return m.moveThroughTheRepair(pressed)
	}
	return m, nil
}

func (m *menu) moveThroughTheSettings(pressed tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case pressed.Code == tea.KeyUp:
		m.cursor = stepBack(m.cursor, len(m.settings))
	case pressed.Code == tea.KeyDown:
		m.cursor = stepOn(m.cursor, len(m.settings))
	case pressed.Code == tea.KeyEnter:
		m.choose(m.settings[m.cursor])
	case pressed.Text == quitKey:
		return m, tea.Quit
	}
	return m, nil
}

func (m *menu) choose(setting Setting) {
	m.stage = choosing
	m.chosen = setting
	m.options = optionsFor(setting)
	m.typed = setting.Value
	m.at = 0
	m.said = nil

	matched := false
	for position, offered := range m.options {
		if !offered.editable && !offered.defaults && offered.value == setting.Value {
			m.at = position
			matched = true
		}
	}
	if !matched && setting.Typeable && setting.Value != "" {
		m.at = len(m.options) - 2
	}
}

func (m *menu) moveThroughTheOptions(pressed tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	standing := m.options[m.at]

	switch {
	case pressed.Code == tea.KeyEsc:
		m.stage = listing
		m.said = nil
	case pressed.Code == tea.KeyUp:
		m.at = stepBack(m.at, len(m.options))
	case pressed.Code == tea.KeyDown:
		m.at = stepOn(m.at, len(m.options))
	case pressed.Code == tea.KeyEnter:
		return m, m.take(standing)
	case standing.editable && pressed.Code == tea.KeyBackspace:
		m.typed = withoutTheLastCharacter(m.typed)
	case standing.editable && pressed.Text != "":
		m.typed += pressed.Text
	}
	return m, nil
}

func (m *menu) take(taken option) tea.Cmd {
	if taken.defaults {
		return m.record("", true)
	}

	value := taken.value
	if taken.editable {
		value = strings.TrimSpace(m.typed)
	}
	if value == "" {
		m.said = []string{fmt.Sprintf(nothingWasTyped, m.chosen.Key, theDefaultLabel(m.chosen))}
		return nil
	}
	return m.record(value, false)
}

func (m *menu) record(value string, removing bool) tea.Cmd {
	key := m.chosen.Key
	if !removing && m.definition.Validate != nil {
		if err := m.definition.Validate(m.chosen, value); err != nil {
			m.said = strings.Split(err.Error(), newline)
			return nil
		}
	}
	written, changed, err := m.definition.withSetting(m.path, m.chosen.writtenAt(), value, removing)
	if err != nil {
		m.said = strings.Split(err.Error(), newline)
		return nil
	}
	if !changed {
		m.said = []string{fmt.Sprintf(alreadyTheDefault, key, m.path)}
		m.stage = listing
		return nil
	}
	if err := m.definition.writeAtomically(m.path, written); err != nil {
		return m.giveUp(err)
	}

	said := append([]string{fmt.Sprintf(setTo, key, value, m.path)}, m.definition.afterAWriteTo(m.chosen)...)
	if removing {
		said = append([]string{fmt.Sprintf(nowTheDefault, key, m.path)}, m.definition.afterAWriteTo(m.chosen)...)
	}
	if err := m.reread(); err != nil {
		return m.giveUp(err)
	}
	m.said = said
	return nil
}

func (m *menu) giveUp(err error) tea.Cmd {
	m.failed = err
	return tea.Quit
}

func (m *menu) moveThroughTheRepair(pressed tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case pressed.Code == tea.KeyUp:
		m.at = stepBack(m.at, len(m.options))
	case pressed.Code == tea.KeyDown:
		m.at = stepOn(m.at, len(m.options))
	case pressed.Text == quitKey:
		return m, tea.Quit
	case pressed.Code == tea.KeyEnter:
		if m.options[m.at].quits {
			return m, tea.Quit
		}
		return m, m.repair()
	}
	return m, nil
}

func (m *menu) repair() tea.Cmd {
	one := m.broken[0]

	content, err := m.definition.held(m.path)
	if err != nil {
		return m.giveUp(err)
	}
	written, err := repaired(content, one)
	if err != nil {
		return m.giveUp(err)
	}
	if err := m.definition.writeAtomically(m.path, written); err != nil {
		return m.giveUp(err)
	}
	if err := m.reread(); err != nil {
		return m.giveUp(err)
	}
	if m.stage == listing {
		m.said = append([]string{fmt.Sprintf(noLongerStated, one.key, m.path)}, m.definition.afterAWriteTo(Setting{Key: one.key})...)
	}
	return nil
}

func optionsFor(setting Setting) []option {
	var offered []option
	for _, choice := range setting.Choices {
		offered = append(offered, option{label: choice, value: choice})
	}
	if len(offered) == 0 || setting.Typeable {
		offered = append(offered, option{editable: true})
	}
	return append(offered, option{label: theDefaultLabel(setting), defaults: true})
}

func theDefaultLabel(setting Setting) string {
	if setting.Default == "" {
		return useTheDefault
	}
	return fmt.Sprintf(useTheDefaultOf, setting.Default)
}

func repairOptions(one unreadable) []option {
	return []option{
		{label: fmt.Sprintf(resetToItsDefault, one.key)},
		{label: leaveTheFileAlone, quits: true},
	}
}

func stepBack(at, total int) int {
	if at <= 0 {
		return total - 1
	}
	return at - 1
}

func stepOn(at, total int) int {
	if at >= total-1 {
		return 0
	}
	return at + 1
}

func withoutTheLastCharacter(typed string) string {
	runes := []rune(typed)
	if len(runes) == 0 {
		return typed
	}
	return string(runes[:len(runes)-1])
}
