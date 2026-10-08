package harness

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creack/pty"
)

const (
	screenRows    = 60
	screenColumns = 220

	cursor = "> "

	Up     = "\x1b[A"
	Down   = "\x1b[B"
	Enter  = "\r"
	Escape = "\x1b"

	drawnWait   = 15 * time.Second
	drawnPoll   = 5 * time.Millisecond
	changedWait = 2 * time.Second
	settledWait = 20 * time.Millisecond
)

type Terminal struct {
	t         *testing.T
	arguments []string
	command   *exec.Cmd
	device    *os.File
	screen    *screen
	finished  chan error
	stopped   sync.Once
	typedOver string
	typed     bool
}

func (s *Session) SocratesInATerminal(arguments ...string) *Terminal {
	s.t.Helper()
	return s.SocratesInATerminalWith(nil, arguments...)
}

func (s *Session) SocratesInATerminalWith(environment []string, arguments ...string) *Terminal {
	s.t.Helper()
	return s.startInATerminal(s.Repository.Root, environment, append([]string{"review"}, arguments...))
}

func (s *Session) TopInATerminal(arguments ...string) *Terminal {
	s.t.Helper()
	return s.TopInATerminalIn(s.Repository.Root, arguments...)
}

func (s *Session) TopInATerminalIn(directory string, arguments ...string) *Terminal {
	s.t.Helper()
	return s.startInATerminal(directory, nil, arguments)
}

func (s *Session) startInATerminal(directory string, environment []string, arguments []string) *Terminal {
	s.t.Helper()

	command := exec.Command(s.Binary, arguments...)
	command.Dir = directory
	command.Env = s.Environment(environment...)

	device, err := pty.StartWithSize(command, &pty.Winsize{Rows: screenRows, Cols: screenColumns})
	if err != nil {
		s.t.Fatalf("starting socrates %s in a pseudo-terminal: %v", strings.Join(arguments, " "), err)
	}

	attached := &Terminal{
		t:         s.t,
		arguments: arguments,
		command:   command,
		device:    device,
		screen:    newScreen(),
		finished:  make(chan error, 1),
	}

	go attached.read()
	go func() { attached.finished <- command.Wait() }()
	s.t.Cleanup(attached.stop)
	return attached
}

func (a *Terminal) read() {
	held := make([]byte, 4096)
	for {
		count, err := a.device.Read(held)
		if count > 0 {
			a.screen.write(held[:count])
		}
		if err != nil {
			return
		}
	}
}

func (a *Terminal) Type(keys string) {
	a.t.Helper()

	a.typedOver, a.typed = a.Drawn(), true
	if _, err := io.WriteString(a.device, keys); err != nil {
		a.t.Fatalf("typing %q at socrates %s: %v", keys, strings.Join(a.arguments, " "), err)
	}
}

func (a *Terminal) Drawn() string {
	return a.screen.text()
}

func (a *Terminal) Await(wanted string) string {
	a.t.Helper()

	deadline := time.Now().Add(drawnWait)
	for time.Now().Before(deadline) {
		if drawn := a.Drawn(); strings.Contains(drawn, wanted) {
			return drawn
		}
		time.Sleep(drawnPoll)
	}
	a.t.Fatalf("socrates %s never drew %q, and the screen reads:\n%s", strings.Join(a.arguments, " "), wanted, a.Drawn())
	return ""
}

func (a *Terminal) AwaitSettling() string {
	a.t.Helper()

	a.awaitAnswering()

	settled := a.Drawn()
	for {
		time.Sleep(settledWait)
		drawn := a.Drawn()
		if drawn == settled {
			return settled
		}
		settled = drawn
	}
}

func (a *Terminal) awaitAnswering() {
	if !a.typed {
		return
	}
	a.typed = false

	deadline := time.Now().Add(changedWait)
	for time.Now().Before(deadline) {
		if a.Drawn() != a.typedOver {
			return
		}
		time.Sleep(drawnPoll)
	}
}

func (a *Terminal) MoveTo(named string) string {
	a.t.Helper()

	for range screenRows {
		drawn := a.AwaitSettling()
		if strings.HasPrefix(Under(drawn), named) {
			return drawn
		}
		a.Type(Down)
	}
	a.t.Fatalf("the cursor never reached %s, and the screen reads:\n%s", named, a.Drawn())
	return ""
}

func Under(drawn string) string {
	for line := range strings.SplitSeq(drawn, "\n") {
		if after, marked := strings.CutPrefix(line, cursor); marked {
			return strings.TrimSpace(after)
		}
	}
	return ""
}

func (a *Terminal) Wait() (string, int) {
	a.t.Helper()

	select {
	case err := <-a.finished:
		code := 0
		if err != nil {
			var exited *exec.ExitError
			if !errors.As(err, &exited) {
				a.t.Fatalf("running socrates %s in a pseudo-terminal: %v", strings.Join(a.arguments, " "), err)
			}
			code = exited.ExitCode()
		}
		time.Sleep(drawnPoll)
		return a.Drawn(), code
	case <-time.After(terminalWait):
		a.t.Fatalf("socrates %s is still running after %s, and the screen reads:\n%s", strings.Join(a.arguments, " "), terminalWait, a.Drawn())
		return "", 0
	}
}

func (a *Terminal) stop() {
	a.stopped.Do(func() {
		select {
		case err := <-a.finished:
			a.finished <- err
		default:
			_ = a.command.Process.Kill()
		}
		_ = a.device.Close()
	})
}

type screen struct {
	guard       sync.Mutex
	rows        [][]rune
	row         int
	column      int
	savedRow    int
	savedColumn int
	pending     []byte
}

var controlSequence = regexp.MustCompile(`^\x1b(\[[0-9;?]*[@-~]|\][^\x07\x1b]*(\x07|\x1b\\)?|[()#][0-9A-Za-z]|[0-9A-Za-z=><])`)

func newScreen() *screen {
	return &screen{rows: blankRows()}
}

func blankRows() [][]rune {
	rows := make([][]rune, screenRows)
	for at := range rows {
		rows[at] = blankRow()
	}
	return rows
}

func blankRow() []rune {
	row := make([]rune, screenColumns)
	for at := range row {
		row[at] = ' '
	}
	return row
}

func (s *screen) write(held []byte) {
	s.guard.Lock()
	defer s.guard.Unlock()

	s.pending = append(s.pending, held...)
	for len(s.pending) > 0 {
		if s.pending[0] == 0x1b {
			matched := controlSequence.Find(s.pending)
			if matched == nil {
				if len(s.pending) < 8 {
					return
				}
				s.pending = s.pending[1:]
				continue
			}
			s.apply(string(matched))
			s.pending = s.pending[len(matched):]
			continue
		}

		character, width := decode(s.pending)
		if width == 0 {
			return
		}
		s.put(character)
		s.pending = s.pending[width:]
	}
}

func decode(held []byte) (rune, int) {
	character := rune(held[0])
	if character < 0x80 {
		return character, 1
	}
	for width := 2; width <= 4; width++ {
		if len(held) < width {
			return 0, 0
		}
		if decoded := []rune(string(held[:width])); len(decoded) == 1 && decoded[0] != 0xFFFD {
			return decoded[0], width
		}
	}
	return '?', 1
}

func (s *screen) put(character rune) {
	switch character {
	case '\n':
		s.row = min(s.row+1, screenRows-1)
	case '\r':
		s.column = 0
	case '\b':
		s.column = max(s.column-1, 0)
	case '\t':
		s.column = min(s.column+8-s.column%8, screenColumns-1)
	default:
		if character < ' ' {
			return
		}
		if s.column >= screenColumns {
			s.column = 0
			s.row = min(s.row+1, screenRows-1)
		}
		s.rows[s.row][s.column] = character
		s.column++
	}
}

func (s *screen) apply(sequence string) {
	if !strings.HasPrefix(sequence, "\x1b[") {
		s.applySingle(sequence)
		return
	}

	final := sequence[len(sequence)-1]
	counted := numbersIn(sequence[2 : len(sequence)-1])
	first := 1
	if len(counted) > 0 && counted[0] > 0 {
		first = counted[0]
	}

	switch final {
	case 'A':
		s.row = max(s.row-first, 0)
	case 'B':
		s.row = min(s.row+first, screenRows-1)
	case 'C':
		s.column = min(s.column+first, screenColumns-1)
	case 'D':
		s.column = max(s.column-first, 0)
	case 'E':
		s.row, s.column = min(s.row+first, screenRows-1), 0
	case 'F':
		s.row, s.column = max(s.row-first, 0), 0
	case 'G':
		s.column = min(max(first-1, 0), screenColumns-1)
	case 'H', 'f':
		row, column := 1, 1
		if len(counted) > 0 && counted[0] > 0 {
			row = counted[0]
		}
		if len(counted) > 1 && counted[1] > 0 {
			column = counted[1]
		}
		s.row, s.column = min(row-1, screenRows-1), min(column-1, screenColumns-1)
	case 'J':
		s.erasePage(modeIn(counted))
	case 'K':
		s.eraseRow(modeIn(counted))
	case 'L':
		s.insertRows(first)
	case 'M':
		s.deleteRows(first)
	case 'P':
		s.deleteCharacters(first)
	case '@':
		s.insertCharacters(first)
	case 'X':
		s.eraseCharacters(first)
	case 'S':
		s.scrollUp(first)
	case 'T':
		s.scrollDown(first)
	case 's':
		s.savedRow, s.savedColumn = s.row, s.column
	case 'u':
		s.row, s.column = s.savedRow, s.savedColumn
	}
}

func (s *screen) applySingle(sequence string) {
	switch sequence {
	case "\x1bM":
		if s.row == 0 {
			s.scrollDown(1)
			return
		}
		s.row--
	case "\x1bD":
		if s.row >= screenRows-1 {
			s.scrollUp(1)
			return
		}
		s.row++
	case "\x1bE":
		s.column = 0
		if s.row >= screenRows-1 {
			s.scrollUp(1)
			return
		}
		s.row++
	case "\x1b7":
		s.savedRow, s.savedColumn = s.row, s.column
	case "\x1b8":
		s.row, s.column = s.savedRow, s.savedColumn
	}
}

func (s *screen) eraseCharacters(count int) {
	for at := s.column; at < s.column+count && at < screenColumns; at++ {
		s.rows[s.row][at] = ' '
	}
}

func (s *screen) deleteCharacters(count int) {
	row := s.rows[s.row]
	for at := s.column; at < screenColumns; at++ {
		if at+count < screenColumns {
			row[at] = row[at+count]
			continue
		}
		row[at] = ' '
	}
}

func (s *screen) insertCharacters(count int) {
	row := s.rows[s.row]
	for at := screenColumns - 1; at >= s.column; at-- {
		if at-count >= s.column {
			row[at] = row[at-count]
			continue
		}
		row[at] = ' '
	}
}

func (s *screen) scrollUp(count int) {
	for range count {
		s.rows = append(s.rows[1:], blankRow())
	}
}

func (s *screen) scrollDown(count int) {
	for range count {
		s.rows = append([][]rune{blankRow()}, s.rows[:screenRows-1]...)
	}
}

func modeIn(counted []int) int {
	if len(counted) == 0 {
		return 0
	}
	return counted[0]
}

func numbersIn(parameters string) []int {
	var counted []int
	for _, one := range strings.Split(strings.TrimPrefix(parameters, "?"), ";") {
		read, err := strconv.Atoi(one)
		if err != nil {
			read = 0
		}
		counted = append(counted, read)
	}
	return counted
}

func (s *screen) erasePage(mode int) {
	switch mode {
	case 0:
		s.eraseRow(0)
		for at := s.row + 1; at < screenRows; at++ {
			s.rows[at] = blankRow()
		}
	case 1:
		s.eraseRow(1)
		for at := 0; at < s.row; at++ {
			s.rows[at] = blankRow()
		}
	default:
		s.rows = blankRows()
	}
}

func (s *screen) eraseRow(mode int) {
	switch mode {
	case 0:
		for at := s.column; at < screenColumns; at++ {
			s.rows[s.row][at] = ' '
		}
	case 1:
		for at := 0; at <= s.column && at < screenColumns; at++ {
			s.rows[s.row][at] = ' '
		}
	default:
		s.rows[s.row] = blankRow()
	}
}

func (s *screen) insertRows(count int) {
	for range count {
		s.rows = append(s.rows[:s.row], append([][]rune{blankRow()}, s.rows[s.row:]...)...)
	}
	s.rows = s.rows[:screenRows]
}

func (s *screen) deleteRows(count int) {
	for range count {
		if s.row >= len(s.rows) {
			return
		}
		s.rows = append(s.rows[:s.row], s.rows[s.row+1:]...)
		s.rows = append(s.rows, blankRow())
	}
}

func (s *screen) text() string {
	s.guard.Lock()
	defer s.guard.Unlock()

	drawn := make([]string, 0, len(s.rows))
	for _, row := range s.rows {
		drawn = append(drawn, strings.TrimRight(string(row), " "))
	}
	for len(drawn) > 0 && drawn[len(drawn)-1] == "" {
		drawn = drawn[:len(drawn)-1]
	}
	return strings.Join(drawn, "\n")
}
