// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func RemoveLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: removeLabelCode,
	}
}

func removeLabelCode(e aladino.Env, args []aladino.Value) error {
	label := args[0].(*aladino.StringValue).Val

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	_, _, err := e.GetClient().Issues.GetLabel(e.GetCtx(), owner, repo, label)
	if err != nil {
		return err
	}

	var labelIsAppliedToPullRequest bool = false
	for _, ghLabel := range e.GetPullRequest().Labels {
		if ghLabel.GetName() == label {
			labelIsAppliedToPullRequest = true
			break
		}
	}

	if !labelIsAppliedToPullRequest {
		return nil
	}

	_, err = e.GetClient().Issues.RemoveLabelForIssue(e.GetCtx(), owner, repo, prNum, label)

	return err
}
