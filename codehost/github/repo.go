// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"os"
	"strings"

	git "github.com/libgit2/git2go/v31"
	"github.com/sirupsen/logrus"
)

// CloneRepository clones a repository from a given URL to the provided path.
// url needs to be an HTTPS uri (e.g. https://github.com/libgit2/TestGitRepository)
// path can be empty. In this case, the repository will be cloned to a temporary location and the location is returned.
func CloneRepository(log *logrus.Entry, url string, token string, path string, options *git.CloneOptions) (*git.Repository, string, error) {
	dir := path
	if dir == "" {
		tempDir, err := os.MkdirTemp("", "repository")
		if err != nil {
			log.Error("failed to create temporary folder")
			return nil, "", err
		}
		dir = tempDir
	}

	// TODO: Validate url has the correct format
	splitted := strings.Split(url, "https://")
	if token != "" {
		url = fmt.Sprintf("https://%v@%v", token, splitted[1])
	}

	if options == nil {
		options = &git.CloneOptions{}
	}

	options.FetchOptions = &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: func(url, usernameFromURL string, allowedTypes git.CredentialType) (*git.Credential, error) {
				log.WithFields(logrus.Fields{
					"url":               url,
					"username_from_url": usernameFromURL,
					"allowed_types":     allowedTypes,
				}).Debug("called credentials callback in clone")

				return git.NewCredentialUserpassPlaintext(usernameFromURL, token)
			},
		},
	}

	log.Debugf("cloning %s to %s", url, dir)

	repo, err := git.Clone(url, dir, options)

	return repo, dir, err
}

// CheckoutBranch checks out a given branch in the given repository.
func CheckoutBranch(log *logrus.Entry, repo *git.Repository, branchName string) error {
	checkoutOpts := &git.CheckoutOptions{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	// Lookup for remote branch
	remoteBranch, err := repo.LookupBranch("origin/"+branchName, git.BranchRemote)
	if err != nil {
		log.Errorf("failed to find remote branch: %v", branchName)
		return err
	}
	defer remoteBranch.Free()

	// Lookup for remote branch commit
	remoteCommit, err := repo.LookupCommit(remoteBranch.Target())
	if err != nil {
		log.Errorf("failed to find remote branch commit: %v", branchName)
		return err
	}
	defer remoteCommit.Free()

	// Lookup for local branch
	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if localBranch == nil || err != nil {
		// Create local branch
		localBranch, err = repo.CreateBranch(branchName, remoteCommit, false)
		if err != nil {
			log.Errorf("failed to create local branch: %v", branchName)
			return err
		}

		// Setting upstream to origin branch
		err = localBranch.SetUpstream("origin/" + branchName)
		if err != nil {
			log.Errorf("failed to create upstream to origin/%v", branchName)
			return err
		}
	}
	if localBranch == nil {
		return fmt.Errorf("failed to locate/create local branch: %v", branchName)
	}
	defer localBranch.Free()

	// Lookup for local branch commit
	localCommit, err := repo.LookupCommit(localBranch.Target())
	if err != nil {
		log.Errorf("failed to lookup for commit in local branch: %v", branchName)
		return err
	}
	defer localCommit.Free()

	// Lookup for local branch tree
	tree, err := repo.LookupTree(localCommit.TreeId())
	if err != nil {
		log.Errorf("failed to lookup for local tree: %v", branchName)
		return err
	}
	defer tree.Free()

	// Checkout the tree
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		log.Errorf("failed to checkout tree: %v", branchName)
		return err
	}

	// Set current Head to the checkout branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Errorf("failed to set head: %v", branchName)
		return err
	}

	return nil
}

// RebaseOnto performs a rebase of the current repository Head onto the given branch.
// Inspired by https://github.com/libgit2/git2go/blob/main/rebase_test.go#L359
func RebaseOnto(log *logrus.Entry, repo *git.Repository, branchName string, rebaseOptions *git.RebaseOptions) error {
	ontoBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		log.Errorf("failed to lookup for onto branch: %v", branchName)
		return err
	}
	defer ontoBranch.Free()

	onto, err := repo.AnnotatedCommitFromRef(ontoBranch.Reference)
	if err != nil {
		log.Errorf("failed to extract annotated commit from ref of branch: %v", branchName)
		return err
	}
	defer onto.Free()

	// Start rebase operation
	rebase, err := repo.InitRebase(nil, nil, onto, rebaseOptions)
	if err != nil {
		log.Error("failed to init rebase")
		return err
	}

	// Verify that no operations are already in progress
	rebaseOperationIndex, err := rebase.CurrentOperationIndex()
	if rebaseOperationIndex != git.RebaseNoOperation && err != git.ErrRebaseNoOperation {
		return errors.New("rebase operation already in progress")
	}

	// Iterate over rebase operations based on operation count
	opCount := int(rebase.OperationCount())
	for op := 0; op < opCount; op++ {
		operation, errRebaseNext := rebase.Next()
		if errRebaseNext != nil {
			return errRebaseNext
		}

		// Verify that operation index is correct
		rebaseOperationIndex, err = rebase.CurrentOperationIndex()
		if err != nil {
			return err
		}

		if int(rebaseOperationIndex) != op {
			return errors.New("bad operation index on rebase")
		}

		if !operationsAreEqual(rebase.OperationAt(uint(op)), operation) {
			return errors.New("rebase operations should be equal")
		}

		// Get current rebase operation created commit
		commit, errCommitLookup := repo.LookupCommit(operation.Id)
		if errCommitLookup != nil {
			return errCommitLookup
		}
		defer commit.Free()

		// Apply commit
		err = rebase.Commit(operation.Id, nil, signatureFromCommit(commit), commit.Message())
		if err != nil {
			return err
		}
	}

	defer rebase.Free()

	err = rebase.Finish()
	if err != nil {
		return err
	}

	return nil
}

// Push performs a push of the provided remote/branch.
func Push(log *logrus.Entry, repo *git.Repository, remoteName string, branchName string, token string, force bool) error {
	remote, err := repo.Remotes.Lookup(remoteName)
	if err != nil {
		log.Errorf("failed to find remote: %v", remoteName)
		return err
	}

	refspec := fmt.Sprintf("refs/heads/%s", branchName)
	if force {
		refspec = fmt.Sprintf("+%s", refspec)
	}

	err = remote.Push([]string{refspec}, &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: func(url, usernameFromURL string, allowedTypes git.CredentialType) (*git.Credential, error) {
				log.WithFields(logrus.Fields{
					"url":               url,
					"username_from_url": usernameFromURL,
					"allowed_types":     allowedTypes,
				}).Debug("called credentials callback in push")

				return git.NewCredentialUserpassPlaintext(usernameFromURL, token)
			},
		},
	})
	if err != nil {
		log.Errorf("failed to push to: %v", branchName)
		return err
	}

	return nil
}

func signatureFromCommit(commit *git.Commit) *git.Signature {
	return &git.Signature{
		Name:  commit.Committer().Name,
		Email: commit.Committer().Email,
		When:  commit.Committer().When,
	}
}

func operationsAreEqual(l, r *git.RebaseOperation) bool {
	return l.Exec == r.Exec && l.Type == r.Type && l.Id.String() == r.Id.String()
}
