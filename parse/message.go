package parse

import (
	"regexp"
	"strings"

	git "github.com/libgit2/git2go"
	"golang.org/x/text/unicode/norm"
)

//IsFirstLineTooLong checks file line length
func IsFirstLineTooLong(commit *git.Commit) bool {
	rc := false

	firstLine := strings.TrimRight(strings.Split(commit.Message(), "\n")[0], "\r")

	firstLineIter := norm.Iter{}
	firstLineIter.InitString(norm.NFKD, firstLine)
	characterCount := 0
	for !firstLineIter.Done() {
		characterCount++
		firstLineIter.Next()
	}

	if characterCount > 50 {
		rc = true
	}

	return rc
}

//DoesFirstLineEndWithPeriod checks for unwanted punctuation
func DoesFirstLineEndWithPeriod(commit *git.Commit) bool {
	rc := false

	firstLine := strings.TrimRight(strings.Split(commit.Message(), "\n")[0], "\r")

	contains, _ := regexp.MatchString(".*\\.$", firstLine)

	if contains {
		rc = true
	}

	return rc
}
