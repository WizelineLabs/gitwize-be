package lambda

import (
	"context"
	"database/sql"
	"fmt"
	"gitwize-be/src/db"
	"log"
	"os"
	"strings"

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

func CollectPRs(prSvc PullRequestService) {
	// find repo id
	conn := db.SqlDBConn()
	rows, _ := conn.Query("SELECT id, url FROM repository")

	var url string
	var id int
	for rows.Next() {
		err := rows.Scan(&id, &url)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		} else {
			CollectPRsOfRepo(prSvc, id, url, conn)
		}
	}
}

func CollectPRsOfRepo(prSvc PullRequestService, id int, url string, conn *sql.DB) {
	var owner, repo string
	if strings.HasPrefix(url, "git") {
		s := strings.Split(url, ":")[1]
		owner = strings.Split(s, "/")[0]
		repo = strings.Split(s, "/")[1]
	}
	if strings.HasPrefix(url, "https") {
		s := strings.Split(url, "/")
		owner = s[len(s)-2]
		repo = s[len(s)-1]
	}
	repo = strings.Replace(repo, ".git", "", -1)
	log.Printf("Collecting PRs: owner=%s, repo=%s", owner, repo)
	collectPRsOfRepo(prSvc, id, owner, repo, conn)
}

func collectPRsOfRepo(prSvc PullRequestService, id int, owner string, repo string, conn *sql.DB) {
	deleteStmt, err := conn.Prepare("DELETE FROM pull_request WHERE repository_id = ?")
	defer deleteStmt.Close()

	// delete old data
	_, err = deleteStmt.Exec(id)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}

	// fetch PRs page by page, 100 per_page
	listOpt := github.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	for {
		prs, _, err := prSvc.List(owner, repo, &github.PullRequestListOptions{
			State:       "all",
			Sort:        "created",
			Direction:   "desc",
			ListOptions: listOpt,
		})

		var prnumbers []int
		for _, pr := range prs {
			prnumbers = append(prnumbers, *pr.Number)
		}
		fmt.Println(prnumbers)

		if err != nil {
			log.Printf("[ERROR] %s", err)
			return
		}
		if len(prs) == 0 {
			break
		}
		insertPRs(prSvc, id, prs, conn)

		listOpt.Page++
	}
}

func insertPRs(prSvc PullRequestService, repoID int, prs []*github.PullRequest, conn *sql.DB) {
	// Prepare statements
	insertStmt, err := conn.Prepare("INSERT INTO pull_request (repository_id, url, pr_no, title, body, head, base, state, created_by, created_year, created_month, created_day, created_hour, closed_year, closed_month, closed_day, closed_hour) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}
	defer insertStmt.Close()

	// insert again
	for _, pr := range prs {
		state := "open"
		if *pr.State == "closed" && pr.MergedAt != nil {
			state = "merged"
		} else if pr.State == nil {
			state = "rejected"
		}

		created := pr.CreatedAt.UTC()
		if pr.ClosedAt != nil {
			insertStmt.Exec(repoID, pr.HTMLURL, pr.Number, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, state, pr.User.Login, created.Year(), created.Month(), created.Day(), created.Hour(), pr.ClosedAt.Year(), pr.ClosedAt.Month(), pr.ClosedAt.Day(), pr.ClosedAt.Hour())
		} else {
			insertStmt.Exec(repoID, pr.HTMLURL, pr.Number, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, state, pr.User.Login, created.Year(), created.Month(), created.Day(), created.Hour(), nil, nil, nil, nil)
		}
	}
}
