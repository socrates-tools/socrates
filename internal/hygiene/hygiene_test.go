package hygiene_test

import (
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/hygiene"
)

func TestNeutralisingOneField(t *testing.T) {
	for _, field := range []struct {
		named string
		given string
		want  string
	}{
		{named: "a conflict marker", given: "<<<<<<< HEAD\nmine\n=======\ntheirs\n>>>>>>> other", want: "HEAD\nmine\n\ntheirs\nother"},
		{named: "a carriage return", given: "one\r\ntwo\rthree", want: "one\ntwo\nthree"},
		{named: "a run of spaces", given: "one     two\t\tthree", want: "one two three"},
		{named: "surrounding space", given: "   what is wrong   \n", want: "what is wrong"},
		{named: "nothing to do", given: "The guard reads before it checks.", want: "The guard reads before it checks."},
		{named: "an empty field", given: "", want: ""},
	} {
		if got := hygiene.Neutralise(field.given); got != field.want {
			t.Errorf("%s neutralises to %q, want %q", field.named, got, field.want)
		}
	}
}

func TestALineBreakIsNotAWhitespaceRun(t *testing.T) {
	neutralised := hygiene.Neutralise("the first claim\nthe second claim")

	if !strings.Contains(neutralised, "\n") {
		t.Errorf("collapsing whitespace ran two lines together: %q", neutralised)
	}
}

func TestNeutralisingIsIdempotent(t *testing.T) {
	const given = "<<<<<<< HEAD\r\nthe   guard  is wrong  \n=======\nand   so is the loop\n>>>>>>> theirs\r\n"

	once := hygiene.Neutralise(given)
	if twice := hygiene.Neutralise(once); twice != once {
		t.Errorf("neutralising twice differs from once:\nonce:  %q\ntwice: %q", once, twice)
	}
}
