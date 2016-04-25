package parse

import (
	git "github.com/libgit2/git2go"
	"golang.org/x/text/unicode/norm"
	"strings"
)

func IsFirstLineTooLong(commit *git.Commit) bool {
	rc := false

	firstLine := strings.TrimRight(strings.Split(commit.Message(), "\n")[0], "\r")

	firstLineIter := norm.Iter{}
	firstLineIter.InitString(norm.NFKD, firstLine)
	characterCount := 0
	for !firstLineIter.Done() {
		characterCount += 1
		firstLineIter.Next()
	}

	if characterCount > 50 {
		rc = true
	}

	return rc
}
