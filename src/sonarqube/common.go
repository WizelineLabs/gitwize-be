package sonarqube

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	gogit_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"gitwize-be/src/configuration"
	"gitwize-be/src/cypher"
	"gitwize-be/src/utils"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	sonarQubeAPIProjectCreate = "/api/projects/create?"
	sonarQubeAPITokenCreate   = "/api/user_tokens/generate?"
	sonarQubeAPIGetMetric     = "/api/project_badges/measure?"
	sonarAdminUser            = "admin"
)

const (
	curDirectory string = "~/"
)

type SonarQubeToken struct {
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Token     string    `json:"token"`
}

func performHttpRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	req.SetBasicAuth(sonarAdminUser, configuration.CurConfiguration.SonarQube.AdminSecret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	return resp, nil
}

// Clone repo to local file sys to avoid memory issue
func cloneRepo(repoName, repoURL, token string) {
	defer utils.TimeTrack(time.Now(), "getRepo: "+repoName)

	repoPath := curDirectory + repoName
	os.RemoveAll(repoPath)
	if _, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: &gogit_http.BasicAuth{
			Username: "nonempty",
			Password: cypher.DecryptString(token, os.Getenv("CYPHER_PASS_PHASE")),
		},
		URL:               repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}); err != nil {
		log.Printf("ERR repo: %s, %s\n", repoName, err)
	}

}
