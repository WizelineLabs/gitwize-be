package lambda

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"testing"
)

type mockRemoteClient struct{}

func (m *mockRemoteClient) List(options *git.ListOptions) ([]*plumbing.Reference, error) {
	return []*plumbing.Reference{
		plumbing.NewSymbolicReference("HEAD", "refs/heads/master"),
		plumbing.NewReferenceFromStrings("refs/heads/master", "000000000001"),
		plumbing.NewReferenceFromStrings("refs/heads/hotfix/abc", "000000000002"),
		plumbing.NewReferenceFromStrings("refs/heads/feature/xxx", "000000000003"),
	}, nil
}

func Test_GetRemoteBranches(t *testing.T) {
	client := &mockRemoteClient{}
	branches := GetRemoteBranches(client)
	if len(branches) != 3 {
		t.Errorf("expected %d, got %d", 3, len(branches))
	}
	if branches[0] != "master" {
		t.Errorf("expected %s, got %s", "master", branches[0])
	}
}

func Test_IntegrationGetRemoteBranchesNoAuth(t *testing.T) {
	url := "git@github.com:go-git/go-git.git"
	branches := GetRemoteBranchesNoAuth(url)
	if len(branches) == 0 {
		t.Errorf("expected to have branches, got none")
	}
}
