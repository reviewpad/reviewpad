// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasUnaddressedThreads() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildBoolType()),
		Code:           hasUnaddressedThreadsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasUnaddressedThreadsCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	reviewThreads, err := pullRequest.GetReviewThreads()
	if err != nil {
		return nil, err
	}

	for _, reviewThread := range reviewThreads {
		if !reviewThread.IsResolved && !reviewThread.IsOutdated {
			return lang.BuildTrueValue(), nil
		}
	}

	return lang.BuildFalseValue(), nil
}
