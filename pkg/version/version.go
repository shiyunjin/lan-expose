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
	} else {
		parts = append(parts, "unknown")
	}

	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}
