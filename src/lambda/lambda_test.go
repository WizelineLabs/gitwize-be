package lambda

import (
	"os"
	"strings"
	"testing"
)

func TestTrigger(t *testing.T) {
	repoPayload := RepoPayload{
		RepoID:   1,
		RepoName: "test-repo",
		URL:      "https://non-exist-url-test-repo",
		RepoPass: "",
		Branch:   "",
	}
	funcName := GetLoadFullRepoLambdaFunc()
	err := Trigger(repoPayload, funcName, "ap-southeast-1")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetLoadFullRepoLambdaFunc(t *testing.T) {
	env := os.Getenv("GW_DEPLOY_ENV")
	expectedName := "gitwize-lambda-dev-load_full_one_repo"
	if strings.ToLower(env) == "qa" {
		expectedName = "gitwize-lambda-qa-load_full_one_repo"
	}
	if expectedName != GetLoadFullRepoLambdaFunc() {
		t.Errorf("Error, expected %s, got %s", expectedName, GetLoadFullRepoLambdaFunc())
	}
}
