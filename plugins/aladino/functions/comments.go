// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func Comments() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: commentsCode,
	}
}

func commentsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	ghComments, err := utils.GetPullRequestComments(e.GetCtx(), e.GetClient(), owner, repo, prNum, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	commentsBody := make([]aladino.Value, len(ghComments))
	for i, ghComment := range ghComments {
		commentsBody[i] = aladino.BuildStringValue(ghComment.GetBody())
	}

	return aladino.BuildArrayValue(commentsBody), nil
}
