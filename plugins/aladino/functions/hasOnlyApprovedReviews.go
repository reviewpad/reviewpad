// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasOnlyApprovedReviews() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildBoolType()),
		Code:           hasOnlyApprovedReviewsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasOnlyApprovedReviewsCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pr := e.GetTarget().(*target.PullRequestTarget)
	targetEntity := e.GetTarget().GetTargetEntity()

	hasOnlyApprovedReviews, err := pr.GetAllReviewersApproved(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
	if err != nil {
		return nil, fmt.Errorf("error getting approved reviews. %v", err.Error())
	}

	return lang.BuildBoolValue(hasOnlyApprovedReviews), nil
}
