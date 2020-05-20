package lambda

import (
	"context"
	"database/sql"
	"gitwize-be/src/db"
	"os"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

type PullRequestService interface {
	List(owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

type GithubPullRequestService struct {
	githubClient *github.Client
}

func NewGithubPullRequestService() *GithubPullRequestService {
	return &GithubPullRequestService{
		githubClient: newGithubClient(),
	}
}

func newGithubClient() *github.Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func (s *GithubPullRequestService) List(owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return s.githubClient.PullRequests.List(context.Background(), owner, repo, opts)
}

func CollectPRs(prSvc PullRequestService, owner string, repo string) {
	prs, _, err := prSvc.List(owner, repo, &github.PullRequestListOptions{
		State: "all",
	})

	if err != nil {
		panic(err)
	}

	// find repo id
	conn := db.SqlDBConn()
	rows := conn.QueryRow("SELECT id FROM repository WHERE name = ?", repo) //FIXME repo should have owner
	var repoID int
	err = rows.Scan(&repoID)

	if err != nil {
		panic(err)
	}

	updatePRs(prSvc, repoID, prs, conn)
}

func updatePRs(prSvc PullRequestService, repoID int, prs []*github.PullRequest, conn *sql.DB) {
	// Prepare statements
	insertStmt, err := conn.Prepare("INSERT INTO pull_request (repository_id, url, pr_no, title, body, head, base, state, created_by, created_year, created_month, created_day, created_hour, closed_year, closed_month, closed_day, closed_hour) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	deleteStmt, err := conn.Prepare("DELETE FROM pull_request WHERE repository_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer insertStmt.Close()
	defer deleteStmt.Close()

	// delete old data
	_, err = deleteStmt.Exec(repoID)
	if err != nil {
		panic(err)
	}

	// insert again
	for _, pr := range prs {
		state := "open"
		if *pr.State == "closed" && pr.MergedAt != nil {
			state = "merged"
		} else if pr.State == nil {
			state = "rejected"
		}

		created := pr.CreatedAt.UTC()
		insertStmt.Exec(repoID, pr.HTMLURL, pr.Number, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, state, pr.User.Login, created.Year(), created.Month(), created.Day(), created.Hour(), nil, nil, nil, nil)
	}
}
