package poll

import (
	git "github.com/libgit2/git2go"
)

func GetNewCommits(repoPath string) ([]*git.Commit, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	oldRefs, err := getRemoteReferences(repo)
	if err != nil {
		return nil, err
	}

	err = fetchAndPrune(repo)
	if err != nil {
		return nil, err
	}

	newRefs, err := getRemoteReferences(repo)
	if err != nil {
		return nil, err
	}

	return lookupCommits(repo, oldRefs, newRefs)
}

func lookupCommits(repo *git.Repository, oldRefs, newRefs []*git.Reference) ([]*git.Commit, error) {
	walk, err := repo.Walk()
	if err != nil {
		return nil, err
	}

	oid := new(git.Oid)
	var commits []*git.Commit

	//TODO don't insert duplicate commits into slice
	for _, ref := range newRefs {
		walk.Reset()

		walk.Push(ref.Target())
		for _, oldRef := range oldRefs {
			walk.Hide(oldRef.Target())
		}

		for err := walk.Next(oid); err == nil; err = walk.Next(oid) {
			commit, err := repo.LookupCommit(oid)
			if err != nil {
				//TODO may not want to ignore this
				continue
			}
			commits = append(commits, commit)
		}
	}

	return commits, nil
}

func getRemoteReferences(repo *git.Repository) ([]*git.Reference, error) {
	iter, err := repo.NewReferenceIterator()
	if err != nil {
		return nil, err
	}

	var references []*git.Reference

	for ref, _ := iter.Next(); ref != nil; ref, _ = iter.Next() {
		if ref.IsRemote() {
			target := ref.Target()
			name := ref.Name()
			if name != "" && target != nil {
				references = append(references, ref)
			}
		}
	}

	return references, nil
}

func fetchAndPrune(repo *git.Repository) error {
	remoteNames, err := repo.Remotes.List()

	if err != nil {
		return err
	}

	for _, name := range remoteNames {
		remote, err := repo.Remotes.Lookup(name)

		if err != nil {
			continue
		}

		var empty []string
		options := git.FetchOptions{}
		options.Prune = 1
		remote.Fetch(empty, &options, "")
	}

	return nil
}
