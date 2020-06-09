package lambda

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"gitwize-be/src/utils"
	"log"
	"os"
	"time"
)

// GetRepo clone repo to local file sys to avoid memory issue
func GetRepo(repoName, repoURL, token string) *git.Repository {
	defer utils.TimeTrack(time.Now(), "GetRepo")

	repoPath := directory + "/" + repoName
	os.RemoveAll(repoPath)
	r, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "nonempty", // it need a non empty value here for private repo, look like gogit bug
			Password: token,
		},
		URL:               repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		log.Panic(err)
	}
	return r
}

//GetCommitIterFromBranch return CommitIter object
func GetCommitIterFromBranch(r *git.Repository, branch string, dateRange DateRange) object.CommitIter {
	defer utils.TimeTrack(time.Now(), "GetCommitIterFromBranch")

	ref, err := r.Head()
	if err != nil {
		log.Panic(err)
	}

	if len(branch) > 0 { // checkout branch
		w, err := r.Worktree()
		if err != nil {
			log.Panic(err)
		}

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
		})
		if err != nil {
			log.Panic(err)
		}

		ref, err = r.Head()
		if err != nil {
			log.Panic(err)
		}
	}

	commitIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: dateRange.Since, Until: dateRange.Until})
	if err != nil {
		log.Panic(err)
	}
	return commitIter
}
