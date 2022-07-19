// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func Close() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: closeCode,
	}
}

func closeCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestBaseOwnerName(pullRequest)
	repo := utils.GetPullRequestBaseRepoName(pullRequest)

	closedState := "closed"
	pullRequest.State = &closedState
	_, _, err := e.GetClient().PullRequests.Edit(e.GetCtx(), owner, repo, prNum, pullRequest)

	return err
}
