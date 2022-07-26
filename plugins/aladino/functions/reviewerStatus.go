// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func ReviewerStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code: reviewerStatusCode,
	}
}

func reviewerStatusCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	reviewerLogin := args[0].(*aladino.StringValue)

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	reviews, err := utils.GetPullRequestReviews(e.GetCtx(), e.GetClient(), owner, repo, prNum, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	status := ""

	for _, review := range reviews {
		if review.User == nil {
			continue
		}

		if *review.User.Login != reviewerLogin.Val {
			continue
		}

		switch *review.State {
		case "COMMENTED":
			status = "commented"
		case "CHANGES_REQUESTED":
			status = "requested_changes"
		case "APPROVED":
			status = "approved"
		}
	}

	return aladino.BuildStringValue(status), nil
}
