package lambda

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"log"
	"strings"
)

type remoteClient interface {
	List(options *git.ListOptions) ([]*plumbing.Reference, error)
}

func GetRemoteBranches(client remoteClient) []string {
	refList, err := client.List(&git.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	refPrefix := "refs/heads/"
	result := []string{}
	for _, ref := range refList {
		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]
		result = append(result, branchName)
	}
	return result
}

func GetRemoteBranchesNoAuth(url string) []string {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{Name: "origin", URLs: []string{url}})
	return GetRemoteBranches(remote)
}
