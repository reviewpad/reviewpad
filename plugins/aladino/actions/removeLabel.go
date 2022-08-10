// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"log"

	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func RemoveLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: removeLabelCode,
	}
}

func removeLabelCode(e aladino.Env, args []aladino.Value) error {
	labelID := args[0].(*aladino.StringValue).Val

	prNum := gh.GetPullRequestNumber(e.GetPullRequest())
	owner := gh.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := gh.GetPullRequestBaseRepoName(e.GetPullRequest())

	internalLabelID := aladino.BuildInternalLabelID(labelID)

	var labelName string

	if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
		labelName = val.(*aladino.StringValue).Val
	} else {
		labelName = labelID
		log.Printf("[warn]: the %v label was not found in the environment", labelID)
	}

	labelIsAppliedToPullRequest := false
	for _, ghLabel := range e.GetPullRequest().Labels {
		if ghLabel.GetName() == labelName {
			labelIsAppliedToPullRequest = true
			break
		}
	}

	if !labelIsAppliedToPullRequest {
		return nil
	}

	_, err := e.GetGithubClient().RemoveLabelForIssue(e.GetCtx(), owner, repo, prNum, labelName)

	return err
}
