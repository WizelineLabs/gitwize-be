package githubapi

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"os"
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
