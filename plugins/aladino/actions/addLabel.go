// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"log"

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
	labelID := args[0].(*aladino.StringValue).Val

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	internalLabelID := aladino.BuildInternalLabelID(labelID)

	var labelName string

	if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
		labelName = val.(*aladino.StringValue).Val
	} else {
		labelName = labelID
		log.Printf("[warn]: addLabel %v was not found in the environment", labelID)
	}

	_, _, err := e.GetClient().Issues.AddLabelsToIssue(e.GetCtx(), owner, repo, prNum, []string{labelName})

	return err
}
