package version

import (
	"fmt"
	"strings"
)

var (
	tag    string
	commit string
	branch string
)

func FormatFullVersion(name string) string {
	var parts = []string{name}

	if tag != "" {
		parts = append(parts, tag)
	} else if branch != "" {
		parts = append(parts, branch)
	} else {
		parts = append(parts, "unknown")
	}

	if commit != "" {
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s)", commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}
