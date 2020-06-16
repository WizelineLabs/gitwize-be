package githubapi

import (
	"context"
	"errors"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

var notAbleToParseUrl string = "Not be able to parse repository url"

type githubSvc struct {
	GithubClient *github.Client
}

func newGithubSvc(token string) *githubSvc {
	return &githubSvc{
		GithubClient: newGithubClient(token),
	}
}

func newGithubClient(token string) *github.Client {
	ctx := context.Background()

	if token == "" && os.Getenv("USE_DEFAULT_API_TOKEN") == "true" {
		token = os.Getenv("DEFAULT_GITHUB_TOKEN")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func parseGithubUrl(repoUrl string) (string, string, error) {
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
	return owner, repo, nil
}
