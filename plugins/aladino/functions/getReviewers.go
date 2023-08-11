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

func GetReviewers() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           getReviewersCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func getReviewersCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	state := args[0].(*lang.StringValue).Val
	pr := e.GetTarget().(*target.PullRequestTarget)
	targetEntity := e.GetTarget().GetTargetEntity()

	filteredReviewers, err := pr.GetAllReviewers(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, state, targetEntity.Number)
	if err != nil {
		return nil, fmt.Errorf("error getting all reviewers. %v", err.Error())
	}

	reviewers := []lang.Value{}
	for _, reviewer := range filteredReviewers {
		reviewers = append(reviewers, lang.BuildStringValue(reviewer))
	}

	return lang.BuildArrayValue(reviewers), nil
}
