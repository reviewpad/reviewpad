// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func RequestedReviewers() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           requestedReviewersCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func requestedReviewersCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	usersReviewers := pullRequest.RequestedReviewers.Users
	teamReviewers := pullRequest.RequestedReviewers.Teams
	totalReviewers := len(usersReviewers) + len(teamReviewers)
	reviewersLogin := make([]aladino.Value, totalReviewers)

	for i, userReviewer := range usersReviewers {
		reviewersLogin[i] = aladino.BuildStringValue(userReviewer.GetLogin())
	}

	for i, teamReviewer := range teamReviewers {
		reviewersLogin[i+len(usersReviewers)] = aladino.BuildStringValue(teamReviewer.GetAlias())
	}

	return aladino.BuildArrayValue(reviewersLogin), nil
}
