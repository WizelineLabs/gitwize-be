package lambda

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"testing"
)

func Test_GetDTOFromCommitObject(t *testing.T) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: "https://github.com/sang-d/mock-repo.git",
	})
	if err != nil {
		panic(err.Error())
	}
	expectedHash := "e15d6dad1576edf08811cb1b85a80c23b6d91153"
	expectedEmail := "sang.dinh@wizeline.com"

	commit, _ := r.CommitObject(plumbing.NewHash(expectedHash))
	dto := getCommitDTO(commit)

	if dto.Hash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, dto.Hash)
	}
	if dto.Author != expectedEmail {
		t.Errorf("expected author %s, got %s", expectedEmail, dto.Author)
	}
	if dto.NumParents != 1 {
		t.Errorf("expected number parents %d, got %d", 1, dto.NumParents)
	}
	if dto.AdditionLOC != 2 {
		t.Errorf("expected addition loc %d, got %d", 2, dto.AdditionLOC)
	}
	if dto.DeletionLOC != 0 {
		t.Errorf("expected deletion loc %d, got %d", 0, dto.DeletionLOC)
	}
	if dto.NumFiles != 1 {
		t.Errorf("expected num files loc %d, got %d", 1, dto.NumFiles)
	}
}

func Test_GetFileStatsFromCommitObject(t *testing.T) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: "https://github.com/sang-d/mock-repo.git",
	})
	if err != nil {
		panic(err.Error())
	}
	expectedHash := "e15d6dad1576edf08811cb1b85a80c23b6d91153"

	commit, _ := r.CommitObject(plumbing.NewHash(expectedHash))
	dtos := getFileStatDTO(commit, 1)

	if len(dtos) != 1 {
		t.Errorf("expected num files changes %d, got %d", 1, len(dtos))
	}
	dto := dtos[0]
	if dto.Hash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, dto.Hash)
	}
	if dto.FileName != "hello.txt" {
		t.Errorf("expected file name %s, got %s", "hello.txt", dto.FileName)
	}
	if dto.AdditionLOC != 2 {
		t.Errorf("expected addition loc %d, got %d", 2, dto.AdditionLOC)
	}
	if dto.DeletionLOC != 0 {
		t.Errorf("expected deletion loc %d, got %d", 0, dto.DeletionLOC)
	}
}
