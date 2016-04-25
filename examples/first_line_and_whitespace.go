package main

import (
	"flag"
	"fmt"
	git "github.com/libgit2/git2go"
	"os"
	"shamebot/parse"
	"shamebot/poll"
	"time"
)

var repoName = flag.String("repo", "", "Path to git repository")

func main() {
	flag.Parse()

	if flag.NFlag() < 1 {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	repo, err := git.OpenRepository(*repoName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open repository\n")
		fmt.Fprintf(os.Stderr, "Got error: %v\n", err)
		os.Exit(1)
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find remote repo \"origin\"\n")
		fmt.Fprintf(os.Stderr, "Got error: %v\n", err)
		os.Exit(1)
	}

	for {
		commits, _ := poll.GetNewCommits(*repoName)

		for _, commit := range commits {
			fmt.Println(commit.Author().Name, "pushed", commit.Id(), "\""+commit.Summary()+"\"", "to", remote.Url())
			if parse.IsFirstLineTooLong(commit) {
				fmt.Println(commit.Id(), "First line is too long.")
			}
			if parse.ContainsTrailingWhiteSpace(repo, commit) {
				fmt.Println(commit.Id(), "Contains trailing whitespace")
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
