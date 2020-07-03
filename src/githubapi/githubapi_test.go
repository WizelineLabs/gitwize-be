package githubapi

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_ParseGithubUrl_NotCorrectRepoUrl(t *testing.T) {
	wrongGithubUrl := "http://wizeline/gitwize-fe"
	_, _, err := ParseGithubUrl(wrongGithubUrl)
	assert.Equal(t, notAbleToParseUrl, err.Error())

}

func Test_GetListBranches_NotExistRepoName(t *testing.T) {
	owner := "wizeline"
	notExistRepoName := "not-exist-repo-name"
	expectedResult := "Not Found"

	_, err := GetListBranches(owner, notExistRepoName, "")
	assert.Regexp(t, regexp.MustCompile(expectedResult), err.Error())
}

func Test_GetListBranches_BadAuthorization(t *testing.T) {
	owner := "wizeline"
	repoName := "gitwize-be"
	expectedResult := "Bad credentials"

	_, err := GetListBranches(owner, repoName, "Bad Token")
	assert.Regexp(t, regexp.MustCompile(expectedResult), err.Error())
}

func Test_GetListBranches_OK(t *testing.T) {
	owner := "wizeline"
	repoName := "gitwize-be"
	expectedResult := "\\[.*\\]"

	branches, err := GetListBranches(owner, repoName, "")
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile(expectedResult), branches)
}
