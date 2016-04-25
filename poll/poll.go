package poll

import (
	git "github.com/libgit2/git2go"
	"strings"
)

//TODO: error handling should be reviewed and updated for this whole file

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
	commits := make(map[git.Oid]*git.Commit)

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
			commits[*commit.Id()] = commit
		}
	}

	commitsArray := make([]*git.Commit, len(commits))
	i := 0
	for k := range commits {
		commitsArray[i] = commits[k]
		i++
	}
	return commitsArray, nil
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
		var callbacks git.RemoteCallbacks
		callbacks.CredentialsCallback = credentialsCallback
		callbacks.CertificateCheckCallback = certificateCheckCallback
		options.RemoteCallbacks = callbacks
		options.Prune = 1
		//TODO: check error
		remote.Fetch(empty, &options, "")
	}

	return nil
}

func credentialsCallback(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
	err, cred := git.NewCredSshKeyFromAgent("git")
	errorCode := git.ErrorCode(err)
	return errorCode, &cred
}

//TODO: check certificate
func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	var rc git.ErrorCode
	if strings.Compare(hostname, "github.com") == 0 {
		rc = git.ErrorCode(0)
	} else {
		rc = git.ErrorCode(1)
	}
	return rc
}
