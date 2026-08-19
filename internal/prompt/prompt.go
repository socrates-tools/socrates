package prompt

import (
	"fmt"
	"strings"
)

type Request struct {
	Intent     string
	Repository string
	Base       string
	Snapshot   string
	Changed    []string
}

func Build(request Request) string {
	var built strings.Builder

	built.WriteString("You are reviewing a change you did not write.\n\n")
	built.WriteString("Your working directory is a worktree of a commit that cannot move under you.\n")
	fmt.Fprintf(&built, "It holds the repository %s as it stands now, and its parent %s is what the change is measured against.\n\n", request.Repository, short(request.Base))

	if request.Intent != "" {
		built.WriteString("What the person asked for:\n")
		built.WriteString(request.Intent)
		built.WriteString("\n\n")
	}

	built.WriteString("The files this change touches:\n")
	for _, path := range request.Changed {
		fmt.Fprintf(&built, "- %s\n", path)
	}
	built.WriteString("\nReport what you find as the schema requires, and nothing else.\n")

	return built.String()
}

func short(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}
