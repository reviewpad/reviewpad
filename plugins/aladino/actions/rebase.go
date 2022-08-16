// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"os"

	git "github.com/libgit2/git2go/v31"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Rebase() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: rebaseCode,
	}
}

func rebaseCode(e aladino.Env, args []aladino.Value) error {
	githubToken := os.Getenv("INPUT_TOKEN")
	baseOwner := gh.GetPullRequestBaseOwnerName(e.GetPullRequest())
	baseRepo := gh.GetPullRequestBaseRepoName(e.GetPullRequest())
	number := gh.GetPullRequestNumber(e.GetPullRequest())

	pr, _, err := e.GetGithubClient().GetPullRequest(e.GetCtx(), baseOwner, baseRepo, number)
	if err != nil {
		return err
	}

	if !*pr.Rebaseable {
		return errors.New("the pull request is not rebaseable")
	}

	headHTMLUrl := pr.Head.GetRepo().GetHTMLURL()
	headRef := pr.Head.GetRef()
	baseRef := pr.Base.GetRef()

	repo, dir, err := gh.CloneRepository(headHTMLUrl, githubToken, "", &git.CloneOptions{
		CheckoutBranch: baseRef,
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = gh.CheckoutBranch(repo, headRef)
	if err != nil {
		return err
	}

	err = gh.RebaseOnto(repo, baseRef, nil)
	if err != nil {
		return err
	}

	err = gh.Push(repo, "origin", headRef, true)
	if err != nil {
		return err
	}

	return nil
}
