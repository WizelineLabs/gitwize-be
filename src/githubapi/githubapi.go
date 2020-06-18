package githubapi

import (
	"context"
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
		owner, repoName, nil); err != nil {
		log.Println("GetListBranches: ", err.Error())
		return nil, err
	} else {
		for _, branch := range branchInfos {
			branches = append(branches, branch.GetName())
		}
		return branches, nil
	}
}
