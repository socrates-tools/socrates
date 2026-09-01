package guidance

import (
	"context"
	"fmt"
	"slices"

	"github.com/ktei/socrates/internal/hygiene"
	"github.com/ktei/socrates/internal/refusal"
	"github.com/ktei/socrates/internal/snapshot"
)

const Path = ".socrates/review-guide.md"

const Ceiling = 16 << 10

type Outcome string

const (
	None            Outcome = "none"
	Carried         Outcome = "carried"
	OnlyInTheChange Outcome = "only-in-the-change"
)

type Guide struct {
	Text      string
	AtTheBase bool
}

func At(ctx context.Context, repositoryRoot, base string) (Guide, error) {
	size, present, err := snapshot.BlobSizeAt(ctx, repositoryRoot, base, Path)
	if err != nil {
		return Guide{}, fmt.Errorf(readingTheGuide, Path, base, err)
	}
	if !present {
		return Guide{}, nil
	}
	if size > Ceiling {
		return Guide{}, refusal.New(overTheCeiling, Path, size, Ceiling)
	}

	content, _, err := snapshot.FileAt(ctx, repositoryRoot, base, Path)
	if err != nil {
		return Guide{}, fmt.Errorf(readingTheGuide, Path, base, err)
	}
	return Guide{Text: hygiene.Neutralise(string(content)), AtTheBase: true}, nil
}

func (g Guide) Reached(changed []string) Outcome {
	if g.Text != "" {
		return Carried
	}
	if !g.AtTheBase && slices.Contains(changed, Path) {
		return OnlyInTheChange
	}
	return None
}
