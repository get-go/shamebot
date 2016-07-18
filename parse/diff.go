package parse

import (
	"errors"
	git "github.com/libgit2/git2go"
	"regexp"
	"strings"
)

//ContainsTrailingWhiteSpace checks for extra whitespace
func ContainsTrailingWhiteSpace(repo *git.Repository, commit *git.Commit) (bool, error) {
	//TODO: Error handling
	changeset, err := getChangeset(repo, commit)
	if err != nil {
		return false, err
	}

	rc := false
	changesetLines := strings.Split(changeset, "\n")

	for _, line := range changesetLines {
		if endsWithWhiteSpace(line) {
			rc = true
			break
		}
	}

	return rc, nil
}

func endsWithWhiteSpace(s string) bool {
	contains, _ := regexp.MatchString("^\\+.*[ \t]$", s)
	rc := false
	if contains {
		rc = true
	}

	return rc
}

//TODO: Error handling
func getChangeset(repo *git.Repository, commit *git.Commit) (string, error) {
	if commit.ParentCount() < 1 {
		//TODO: Add error?
		return "", nil
	} else if commit.ParentCount() > 1 {
		//TODO: Figure out how to handle merge commits. libgit2 doesn't seem to support this.
		return "", errors.New("Changeset for merge is unsupported")
	}

	parentCommit := commit.Parent(0)

	parentObject, _ := parentCommit.Peel(git.ObjectTree)
	currentObject, _ := commit.Peel(git.ObjectTree)

	parentTree, _ := parentObject.AsTree()
	currentTree, _ := currentObject.AsTree()

	//TODO: Handle commits with multiple parents
	diff, _ := repo.DiffTreeToTree(parentTree, currentTree, nil)

	var patchString string

	numDeltas, _ := diff.NumDeltas()
	for i := 0; i < numDeltas; i++ {
		patch, _ := diff.Patch(i)
		patchS, _ := patch.String()
		patchString += patchS
	}

	return patchString, nil
}
