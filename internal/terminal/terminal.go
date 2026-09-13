package terminal

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Output struct {
	To       io.Writer
	Failures io.Writer
	Terminal bool
	Pager    []string
}

func For(out, failures io.Writer) Output {
	return Output{To: out, Failures: failures, Terminal: Attached(out), Pager: strings.Fields(os.Getenv(pagerVariable))}
}

func Attached(stream any) bool {
	file, isFile := stream.(*os.File)
	if !isFile {
		return false
	}

	written, err := file.Stat()
	return err == nil && written.Mode()&os.ModeCharDevice != 0
}

func (o Output) Write(ctx context.Context, lines []string) error {
	if o.Terminal && len(o.Pager) > 0 {
		if err := o.page(ctx, lines); err == nil {
			return nil
		}
	}
	return o.plainly(lines)
}

func (o Output) page(ctx context.Context, lines []string) error {
	paging := exec.CommandContext(ctx, o.Pager[0], o.Pager[1:]...)
	paging.Stdin = strings.NewReader(text(lines))
	paging.Stdout = o.To
	paging.Stderr = o.Failures

	if err := paging.Start(); err != nil {
		return fmt.Errorf(pagerWouldNotStart, o.Pager[0], err)
	}
	_ = paging.Wait()
	return nil
}

func (o Output) plainly(lines []string) error {
	if _, err := io.WriteString(o.To, text(lines)); err != nil {
		return fmt.Errorf(writing, err)
	}
	return nil
}

func text(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}
