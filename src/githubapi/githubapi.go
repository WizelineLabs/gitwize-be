package githubapi

import (
	"context"
	"errors"
	"github.com/google/go-github/v32/github"
	"gitwize-be/src/utils"
	"log"
	"strings"
)

func ParseGithubUrl(repoUrl string) (string, string, error) {
	var owner, repo string
	if strings.HasPrefix(repoUrl, "git") {
		s := strings.Split(repoUrl, ":")[1]
		owner = strings.Split(s, "/")[0]
		repo = strings.Split(s, "/")[1]
	} else if strings.HasPrefix(repoUrl, "https") {
		s := strings.Split(repoUrl, "/")
		owner = s[len(s)-2]
		repo = s[len(s)-1]
	} else {
		err := errors.New(notAbleToParseUrl)
		return "", "", err
	}
	repo = strings.Replace(repo, ".git", "", -1)
	log.Println(utils.GetFuncName(), ": owner=", owner, ", repo=", repo)
	return owner, repo, nil
}

func GetListBranches(owner, repoName, accessToken string) ([]string, error) {
	branches := make([]string, 0)

	githubSvcClient := newGithubSvc(accessToken)
	if branchInfos, _, err := githubSvcClient.GithubClient.Repositories.ListBranches(context.Background(),
		owner, repoName, &github.BranchListOptions{Protected: nil, ListOptions: github.ListOptions{Page: 1, PerPage: 100}}); err != nil {
		log.Println(utils.GetFuncName()+": ", err.Error())
		return nil, err
	} else {
		for _, branch := range branchInfos {
			branches = append(branches, branch.GetName())
		}
		return branches, nil
	}
}
