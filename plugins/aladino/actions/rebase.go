// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"os"

	git "github.com/libgit2/git2go/v31"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Rebase() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           rebaseCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func rebaseCode(e aladino.Env, args []aladino.Value) error {
	githubToken := e.GetGithubClient().GetToken()
	log := e.GetLogger().WithField("builtin", "rebase")
	t := e.GetTarget().(*target.PullRequestTarget)
	pr := t.PullRequest

	if !pr.Rebaseable {
		return errors.New("the pull request is not rebaseable")
	}

	headHTMLUrl := pr.Head.GetRepo().GetUri()
	headRef := pr.Head.Name
	baseRef := pr.Base.Name

	repo, dir, err := gh.CloneRepository(log, headHTMLUrl, githubToken, "", &git.CloneOptions{
		CheckoutBranch: baseRef,
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = gh.CheckoutBranch(log, repo, headRef)
	if err != nil {
		return err
	}

	err = gh.RebaseOnto(log, repo, baseRef, nil)
	if err != nil {
		return err
	}

	err = gh.Push(log, repo, "origin", headRef, githubToken, true)
	if err != nil {
		return err
	}

	return nil
}
