// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"os"

	git "github.com/libgit2/git2go/v31"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Rebase() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           rebaseCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func rebaseCode(e aladino.Env, args []aladino.Value) error {
	githubToken := os.Getenv("INPUT_TOKEN")
	t := e.GetTarget().(*target.PullRequestTarget)
	pr := t.PullRequest

	if !*pr.Rebaseable {
		return errors.New("the pull request is not rebaseable")
	}

	headHTMLUrl := pr.Head.GetRepo().GetHTMLURL()
	headRef := pr.Head.GetRef()
	baseRef := pr.Base.GetRef()

	repo, dir, err := gh.CloneRepository(e.GetLogger(), headHTMLUrl, githubToken, "", &git.CloneOptions{
		CheckoutBranch: baseRef,
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = gh.CheckoutBranch(e.GetLogger(), repo, headRef)
	if err != nil {
		return err
	}

	err = gh.RebaseOnto(e.GetLogger(), repo, baseRef, nil)
	if err != nil {
		return err
	}

	err = gh.Push(e.GetLogger(), repo, "origin", headRef, true)
	if err != nil {
		return err
	}

	return nil
}
