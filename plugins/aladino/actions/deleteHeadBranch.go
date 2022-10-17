// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"

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

	if !*target.PullRequest.Merged && target.PullRequest.ClosedAt == nil {
		return errors.New("pull request should be merged or closed before deleting head branch")
	}

	return e.GetGithubClient().DeleteReference(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, "heads/"+*target.PullRequest.Head.Ref)
}
