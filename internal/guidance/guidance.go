package guidance

import (
	"context"
	"fmt"
	"slices"

	"github.com/ktei/socrates/internal/hygiene"
	"github.com/ktei/socrates/internal/snapshot"
)

const Path = ".socrates/review-guide.md"

const DefaultCeiling = 16 << 10

func ceiling(stated int) int64 {
	if stated > 0 {
		return int64(stated)
	}
	return DefaultCeiling
}

type Outcome string

const (
	None            Outcome = "none"
	Carried         Outcome = "carried"
	OnlyInTheChange Outcome = "only-in-the-change"
	OverTheCeiling  Outcome = "over-the-ceiling"
)

type Guide struct {
	Text      string
	Size      int64
	Ceiling   int64
	AtTheBase bool
	TooLarge  bool
}

func At(ctx context.Context, repositoryRoot, base string, stated int) (Guide, error) {
	size, present, err := snapshot.BlobSizeAt(ctx, repositoryRoot, base, Path)
	if err != nil {
		return Guide{}, fmt.Errorf(readingTheGuide, Path, base, err)
	}
	if !present {
		return Guide{}, nil
	}
	if size > ceiling(stated) {
		return Guide{Size: size, Ceiling: ceiling(stated), AtTheBase: true, TooLarge: true}, nil
	}

	content, _, err := snapshot.FileAt(ctx, repositoryRoot, base, Path)
	if err != nil {
		return Guide{}, fmt.Errorf(readingTheGuide, Path, base, err)
	}
	return Guide{Text: hygiene.Neutralise(string(content)), Size: size, Ceiling: ceiling(stated), AtTheBase: true}, nil
}

func (g Guide) Reached(changed []string) Outcome {
	if g.TooLarge {
		return OverTheCeiling
	}
	if g.Text != "" {
		return Carried
	}
	if !g.AtTheBase && slices.Contains(changed, Path) {
		return OnlyInTheChange
	}
	return None
}
