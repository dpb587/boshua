package semverutil

import "strings"

func IsConstraint(s string) bool {
	// TODO complicate for better certainty?
	return strings.HasSuffix(s, ".x") || strings.Contains(s, ">") || strings.Contains(s, "<") || strings.Contains(s, "+") || strings.Contains(s, ",") || strings.Contains(s, "~")
}
