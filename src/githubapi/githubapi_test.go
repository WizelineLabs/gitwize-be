package githubapi

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_GetListBranches_NotCorrectRepoUrl(t *testing.T) {
	wrongGithubUrl := "http://wizeline/gitwize-fe"
	_, err := GetListBranches(wrongGithubUrl, "")
	assert.Equal(t, notAbleToParseUrl, err.Error())

}

func Test_GetListBranches_NotExistRepoName(t *testing.T) {
	notExistRepoName := "https://github.com/wizeline/not-exist-repo-name"
	expectedResult := ".*404 Not Found.*"

	_, err := GetListBranches(notExistRepoName, "")
	assert.Regexp(t, regexp.MustCompile(expectedResult), err.Error())
}

func Test_GetListBranches_BadAuthorization(t *testing.T) {
	repoName := "https://github.com/wizeline/gitwize-be"
	expectedResult := ".*Bad credentials.*"

	_, err := GetListBranches(repoName, "Bad Token")
	assert.Regexp(t, regexp.MustCompile(expectedResult), err.Error())
}

func Test_GetListBranches_OK(t *testing.T) {
	repoName := "https://github.com/wizeline/gitwize-be"
	expectedResult := "\\[.*\\]"

	branches, err := GetListBranches(repoName, "")
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile(expectedResult), branches)
}
