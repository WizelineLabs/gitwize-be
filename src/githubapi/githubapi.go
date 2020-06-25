package githubapi

import (
	"context"
	"github.com/google/go-github/v32/github"
	"gitwize-be/src/utils"
	"log"
)

func GetListBranches(repoUrl, accessToken string) ([]string, error) {
	var owner, repoName string
	var err error
	branches := make([]string, 0)

	if owner, repoName, err = parseGithubUrl(repoUrl); err != nil {
		return nil, err
	}

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
