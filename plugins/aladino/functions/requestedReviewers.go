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

func RequestedReviewers() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           requestedReviewersCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func requestedReviewersCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	usersReviewers := pullRequest.RequestedReviewers.Users
	teamReviewers := pullRequest.RequestedReviewers.Teams
	totalReviewers := len(usersReviewers) + len(teamReviewers)
	reviewersLogin := make([]lang.Value, totalReviewers)

	for i, userReviewer := range usersReviewers {
		reviewersLogin[i] = lang.BuildStringValue(userReviewer.GetLogin())
	}

	for i, teamReviewer := range teamReviewers {
		reviewersLogin[i+len(usersReviewers)] = lang.BuildStringValue(teamReviewer.Slug)
	}

	return lang.BuildArrayValue(reviewersLogin), nil
}
