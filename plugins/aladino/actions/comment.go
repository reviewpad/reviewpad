// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func Comment() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: commentCode,
	}
}

func commentCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	commentBody := args[0].(*aladino.StringValue).Val

	_, _, err := e.GetClient().Issues.CreateComment(e.GetCtx(), owner, repo, prNum, &github.IssueComment{
		Body: &commentBody,
	})

	return err
}
