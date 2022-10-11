// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"

	git "github.com/libgit2/git2go/v31"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/stretchr/testify/assert"
)

const (
	TestRepo string = "testrepo"
)

func TestCloneRepository_WhenNoPathProvided(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t, false)
	defer cleanupTestRepo(t, repo)

	seedTestRepo(t, repo, "main")

	ref, err := repo.References.Lookup("refs/heads/main")
	checkFatal(t, err)

	repo2, _, err := gh.CloneRepository(repo.Path(), "", "", &git.CloneOptions{
		Bare:           true,
		CheckoutBranch: "main",
	})
	defer cleanupTestRepo(t, repo2)

	checkFatal(t, err)

	ref2, err := repo2.References.Lookup("refs/heads/main")
	checkFatal(t, err)

	if ref.Cmp(ref2) != 0 {
		assert.FailNow(t, "the repository should be cloned into a temporary directory")
	}
}

func TestCloneRepository_WhenPathProvided(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t, false)
	defer cleanupTestRepo(t, repo)

	seedTestRepo(t, repo, "main")

	ref, err := repo.References.Lookup("refs/heads/main")
	checkFatal(t, err)

	path, err := ioutil.TempDir("", TestRepo)
	checkFatal(t, err)

	repo2, _, err := gh.CloneRepository(repo.Path(), "", path, &git.CloneOptions{
		Bare:           true,
		CheckoutBranch: "main",
	})
	defer cleanupTestRepo(t, repo2)

	checkFatal(t, err)

	ref2, err := repo2.References.Lookup("refs/heads/main")
	checkFatal(t, err)

	if ref.Cmp(ref2) != 0 {
		assert.FailNow(t, "the repository should be cloned into the provided path")
	}
}

func TestCloneRepository_WithExternalHTTPUrl(t *testing.T) {
	path, err := ioutil.TempDir("", TestRepo)
	defer os.RemoveAll(path)
	checkFatal(t, err)

	url := "https://github.com/reviewpad/TestGitRepository"
	token := "TOKEN"
	_, _, err = gh.CloneRepository(url, token, path, &git.CloneOptions{})

	assert.Nil(t, err, "cannot clone remote repo via https")
}

func TestCheckoutBranch_BranchDoesNotExists(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t, false)
	defer cleanupTestRepo(t, repo)

	seedTestRepo(t, repo, "main")

	err := gh.CheckoutBranch(repo, "test")

	assert.ErrorContains(t, err, "cannot locate remote-tracking branch 'origin/test'")
}

func TestCheckoutBranch_BranchExists(t *testing.T) {
	t.Parallel()
	remoteRepo := createTestRepo(t, false)
	defer cleanupTestRepo(t, remoteRepo)

	head, _ := seedTestRepo(t, remoteRepo, "main")
	commit, err := remoteRepo.LookupCommit(head)
	checkFatal(t, err)
	defer commit.Free()

	branchName := "test"

	remoteRef, err := remoteRepo.CreateBranch(branchName, commit, true)
	checkFatal(t, err)
	defer remoteRef.Free()

	repo := createTestRepo(t, false)
	defer cleanupTestRepo(t, repo)

	config, err := repo.Config()
	checkFatal(t, err)
	defer config.Free()

	remoteUrl := fmt.Sprintf("file://%s", remoteRepo.Workdir())
	remote, err := repo.Remotes.Create("origin", remoteUrl)
	checkFatal(t, err)
	defer remote.Free()

	err = remote.Fetch([]string{branchName}, nil, "")
	checkFatal(t, err)

	err = gh.CheckoutBranch(repo, branchName)
	currentHead, _ := repo.Head()

	assert.Nil(t, err)
	assert.Equal(t, currentHead.Name(), "refs/heads/"+branchName)
}

func TestRebaseOnto(t *testing.T) {
	// TODO: #309
}

func TestPush(t *testing.T) {
	// Local push doesn't (yet) support pushing to non-bare repos so we need to work with bare repos.
	repo := createTestRepo(t, true)
	defer cleanupTestRepo(t, repo)

	remoteUrl := fmt.Sprintf("file://%s", repo.Path())
	remote, err := repo.Remotes.Create("origin", remoteUrl)
	checkFatal(t, err)
	defer remote.Free()

	branchName := "main"

	loc, err := time.LoadLocation("Europe/Berlin")
	checkFatal(t, err)
	sig := &git.Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.Write()
	checkFatal(t, err)
	treeId, err := idx.WriteTree()
	checkFatal(t, err)

	message := "This is a commit\n"
	tree, err := repo.LookupTree(treeId)
	checkFatal(t, err)
	commitId, err := repo.CreateCommit("HEAD", sig, sig, message, tree)
	checkFatal(t, err)

	commit, err := repo.LookupCommit(commitId)
	checkFatal(t, err)
	_, err = repo.CreateBranch(branchName, commit, false)
	checkFatal(t, err)

	tests := map[string]struct {
		inputRemote     string
		inputBranchName string
		isForcePush     bool
		wantErr         string
	}{
		"when given remote cannot be found": {
			inputRemote:     "non-existing-remote",
			inputBranchName: "test",
			isForcePush:     false,
			wantErr:         "remote 'non-existing-remote' does not exist",
		},
		"when is a force push and push fails": {
			inputRemote:     remote.Name(),
			inputBranchName: "test",
			isForcePush:     true,
			wantErr:         "src refspec 'refs/heads/test' does not match any existing object",
		},
		"when is not a force push and push fails": {
			inputRemote:     remote.Name(),
			inputBranchName: "test",
			isForcePush:     true,
			wantErr:         "src refspec 'refs/heads/test' does not match any existing object",
		},
		"when is a force push and push is successful": {
			inputRemote:     remote.Name(),
			inputBranchName: branchName,
			isForcePush:     true,
			wantErr:         "",
		},
		"when is not a force push and push is successful": {
			inputRemote:     remote.Name(),
			inputBranchName: branchName,
			isForcePush:     false,
			wantErr:         "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := gh.Push(repo, test.inputRemote, test.inputBranchName, test.isForcePush)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "Push() error = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}

func createTestRepo(t *testing.T, isBare bool) *git.Repository {
	path, err := ioutil.TempDir("", TestRepo)
	checkFatal(t, err)

	repo, err := git.InitRepository(path, isBare)
	checkFatal(t, err)

	tmpfile := "README"
	err = ioutil.WriteFile(path+"/"+tmpfile, []byte("foo\n"), 0644)

	checkFatal(t, err)

	return repo
}

func cleanupTestRepo(t *testing.T, r *git.Repository) {
	var err error
	if r.IsBare() {
		err = os.RemoveAll(r.Path())
	} else {
		err = os.RemoveAll(r.Workdir())
	}
	checkFatal(t, err)

	r.Free()
}

func seedTestRepo(t *testing.T, repo *git.Repository, defaultBranch string) (*git.Oid, *git.Oid) {
	loc, err := time.LoadLocation("Europe/Berlin")
	checkFatal(t, err)
	sig := &git.Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.AddByPath("README")
	checkFatal(t, err)
	err = idx.Write()
	checkFatal(t, err)
	treeId, err := idx.WriteTree()
	checkFatal(t, err)

	message := "This is a commit\n"
	tree, err := repo.LookupTree(treeId)
	checkFatal(t, err)
	commitId, err := repo.CreateCommit("HEAD", sig, sig, message, tree)
	checkFatal(t, err)

	mainBranch, _ := repo.LookupBranch(defaultBranch, git.BranchLocal)
	if mainBranch == nil {
		commit, err := repo.LookupCommit(commitId)
		checkFatal(t, err)
		_, err = repo.CreateBranch(defaultBranch, commit, false)
		checkFatal(t, err)
	}

	return commitId, treeId
}

func checkFatal(t *testing.T, err error) {
	if err == nil {
		return
	}

	// The failure happens at wherever we were called, not here
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		assert.FailNow(t, "Unable to get caller")
	}
	assert.FailNow(t, fmt.Sprintf("Fail at %v:%v; %v", file, line, err))
}
