// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ReviewerStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildStringType()),
		Code:           reviewerStatusCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func reviewerStatusCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	reviewerLogin := args[0].(*lang.StringValue)

	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	prNum := pullRequest.Number
	owner := pullRequest.GetBase().GetRepo().GetOwner()
	repo := pullRequest.GetBase().GetRepo().GetName()

	reviews, err := e.GetGithubClient().GetPullRequestReviews(e.GetCtx(), owner, repo, int(prNum))
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

	return lang.BuildStringValue(status), nil
}
