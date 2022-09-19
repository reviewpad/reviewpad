// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Reviewers() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           reviewersCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func reviewersCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)
	reviewers := make([]aladino.Value, 0)
	existsInReviewersList := make(map[string]bool)

	reviews, err := t.GetReviews()
	if err != nil {
		return nil, err
	}

	for _, review := range reviews {
		reviewer := review.User.Login
		if _, ok := existsInReviewersList[reviewer]; !ok {
			reviewers = append(reviewers, aladino.BuildStringValue(reviewer))
			existsInReviewersList[reviewer] = true
		}
	}

	return aladino.BuildArrayValue(reviewers), nil
}
