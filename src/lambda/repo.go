package lambda

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"log"
	"os"
	"time"
)

func getPrivateRepo(url, user, token string) *git.Repository {
	defer timeTrack(time.Now(), "getPrivateRepo")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: user,
			Password: token,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Panic(err)
	}
	return r
}

func getPublicRepo(url string) *git.Repository {
	defer timeTrack(time.Now(), "getPublicRepo")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Panic(err)
	}
	return r
}

func getCommitIterFromBranch(r *git.Repository, branch string, dateRange DateRange) object.CommitIter {
	defer timeTrack(time.Now(), "getCommitIterFromBranch")
	ref, err := r.Head()
	if err != nil {
		log.Fatalln(err)
	}

	if len(branch) > 0 { // checkout branch
		w, err := r.Worktree()
		if err != nil {
			log.Fatalln(err)
		}

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
		})
		if err != nil {
			log.Fatalln(err)
		}

		ref, err = r.Head()
		if err != nil {
			log.Fatalln(err)
		}
	}

	commitIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: dateRange.Since, Until: dateRange.Until})
	if err != nil {
		log.Fatalln(err)
	}
	return commitIter
}

func getRepoBySSH(url, sshFilePath, sshPassPhase string) *git.Repository {
	authenticator, _ := ssh.NewPublicKeysFromFile("git", sshFilePath, sshPassPhase)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		Auth:     authenticator,
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return r
}

// getRepoLocal use for uploaded/downloaded repository
func getRepoLocal(repoPath string) *git.Repository {
	defer timeTrack(time.Now(), "getRepoLocal")
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalln(err)
	}
	return r
}

// GitCloneLocal clone repo to local file sys
func GitCloneLocal(directory, publicURL string) {
	defer timeTrack(time.Now(), "GitCloneLocal")
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               publicURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		log.Panic(err)
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		log.Panic(err)
	}
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Panic(err)
	}
	log.Println("commit", commit)
}
