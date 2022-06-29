// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func AddLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: addLabelCode,
	}
}

func addLabelCode(e aladino.Env, args []aladino.Value) error {
	labelVal := args[0]
	if !labelVal.HasKindOf(aladino.STRING_VALUE) {
		return fmt.Errorf("addLabel: expecting string argument, got %v", labelVal.Kind())
	}

	label := labelVal.(*aladino.StringValue).Val

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	_, _, err := e.GetClient().Issues.GetLabel(e.GetCtx(), owner, repo, label)
	if err != nil {
		return err
	}

	_, _, err = e.GetClient().Issues.AddLabelsToIssue(e.GetCtx(), owner, repo, prNum, []string{label})

	return err
}
