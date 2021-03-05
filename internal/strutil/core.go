package strutil

import (
	"bufio"
	"strings"
)

// SplitLines the provided string.
// Works for any OS. "\n" for non Windows and "\r\n" for Windows.
func SplitLines(str string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
