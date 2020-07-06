package core

import (
	"os"
	"strings"
)

// byName and its associated functions provide a quick method of sorting FileInfos by name
type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SplitBy takes an input string and a delimiter, returning the resulting two strings from the split (ALWAYS 2)
func SplitBy(input, delim string) (string, string) {
	split := strings.SplitN(input, delim, 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], ""
}
