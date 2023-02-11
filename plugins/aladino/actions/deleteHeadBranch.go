// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func DeleteHeadBranch() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           deleteHeadBranch,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func deleteHeadBranch(e aladino.Env, args []aladino.Value) error {
	target := e.GetTarget().(*target.PullRequestTarget)
	targetEntity := target.GetTargetEntity()
	ctx := e.GetCtx()
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	if !*target.PullRequest.Merged && target.PullRequest.ClosedAt == nil {
		return nil
	}

	if target.PullRequest.GetHead().GetRepo().GetFork() {
		e.GetLogger().Warnln("$deleteHeadBranch built-in action doesn't work across forks")
		return nil
	}

	ref := "heads/" + *target.PullRequest.Head.Ref

	exists, err := e.GetGithubClient().RefExists(ctx, owner, repo, "refs/"+ref)
	if err != nil {
		return fmt.Errorf("error getting reference: %w", err)
	}

	if !exists {
		return nil
	}

	return e.GetGithubClient().DeleteReference(ctx, owner, repo, ref)
}
