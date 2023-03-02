// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasUnaddressedThreads() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           hasUnaddressedThreadsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasUnaddressedThreadsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	reviewThreads, err := pullRequest.GetReviewThreads()
	if err != nil {
		return nil, err
	}

	for _, reviewThread := range reviewThreads {
		if !reviewThread.IsResolved && !reviewThread.IsOutdated {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}
