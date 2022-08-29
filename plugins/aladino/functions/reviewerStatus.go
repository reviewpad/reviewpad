// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func ReviewerStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           reviewerStatusCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func reviewerStatusCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	reviewerLogin := args[0].(*aladino.StringValue)

	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	reviews, err := e.GetGithubClient().GetPullRequestReviews(e.GetCtx(), owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	status := ""
	reviewerHasDecision := false

	for _, review := range reviews {
		if review.User == nil || review.State == nil {
			continue
		}

		if *review.User.Login != reviewerLogin.Val {
			continue
		}

		reviewState := *review.State
		if reviewState == "COMMENTED" {
			if reviewerHasDecision {
				continue
			} else {
				status = reviewState
			}
		} else {
			status = reviewState
			reviewerHasDecision = true
		}
	}

	return aladino.BuildStringValue(status), nil
}
