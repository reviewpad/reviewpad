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

func HasAnyReviewers() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildBoolType()),
		Code:           hasAnyReviewersCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasAnyReviewersCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pr := e.GetTarget().(*target.PullRequestTarget)
	targetEntity := e.GetTarget().GetTargetEntity()

	reviewers, err := pr.GetAllReviewers(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
	if err != nil {
		return nil, fmt.Errorf("error getting reviewers. %v", err.Error())
	}

	return lang.BuildBoolValue(len(reviewers) > 0), nil
}
