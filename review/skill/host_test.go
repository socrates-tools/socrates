package skill_test

import (
	"strings"
	"testing"

	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/skill"
)

func TestEveryRenderingPassesItsOwnHostToBegin(t *testing.T) {
	beginning := skillkit.Dispatcher + " review begin"

	for _, one := range skill.All() {
		started := 0
		for _, line := range strings.Split(one.Rendered(), "\n") {
			if !strings.HasPrefix(line, skillkit.CodeIndent+beginning) {
				continue
			}
			started++
			want := "review begin --" + skillkit.HostFlag + " " + one.Host.Flag
			if !strings.Contains(line, want) {
				t.Errorf("%s on %s starts a review with %q, want it to carry %q", one.Name, one.Host.Name, strings.TrimSpace(line), want)
			}
			for _, other := range skillkit.Hosts() {
				if other != one.Host && strings.Contains(line, "--"+skillkit.HostFlag+" "+other.Flag) {
					t.Errorf("%s on %s passes %s's host: %q", one.Name, one.Host.Name, other.Name, strings.TrimSpace(line))
				}
			}
		}
		if started == 0 {
			t.Errorf("%s on %s holds no indented %q to check", one.Name, one.Host.Name, beginning)
		}
	}
}
